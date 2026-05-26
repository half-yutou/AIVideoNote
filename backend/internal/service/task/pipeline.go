package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/logger"
	"github.com/aivideonote/backend/internal/repository"
	"github.com/aivideonote/backend/internal/service/downloader"
	"github.com/aivideonote/backend/internal/service/llm"
	"github.com/aivideonote/backend/internal/service/transcriber"
)

type Pipeline struct {
	cfg             *config.Config
	taskRepo        *repository.TaskRepository
	noteRepo        *repository.NoteRepository
	providerRepo    *repository.ProviderRepository
	cookieRepo      *repository.CookieRepository
	llmClient       *llm.Client
	transcriberCl   *transcriber.Client
	statusBroadcast func(taskID string, status model.TaskStatus, msg string)
}

func NewPipeline(
	cfg *config.Config,
	taskRepo *repository.TaskRepository,
	noteRepo *repository.NoteRepository,
	providerRepo *repository.ProviderRepository,
	cookieRepo *repository.CookieRepository,
	llmClient *llm.Client,
	transcriberCl *transcriber.Client,
	broadcast func(taskID string, status model.TaskStatus, msg string),
) *Pipeline {
	return &Pipeline{
		cfg:             cfg,
		taskRepo:        taskRepo,
		noteRepo:        noteRepo,
		providerRepo:    providerRepo,
		cookieRepo:      cookieRepo,
		llmClient:       llmClient,
		transcriberCl:   transcriberCl,
		statusBroadcast: broadcast,
	}
}

type stepError struct {
	failedStatus model.TaskStatus
	msg          string
	err          error
}

func (e *stepError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func (p *Pipeline) failStep(taskID string, s *stepError) {
	p.taskRepo.UpdateStatus(taskID, s.failedStatus, s.Error())
	if p.statusBroadcast != nil {
		p.statusBroadcast(taskID, s.failedStatus, s.Error())
	}
	logger.L.Errorf("任务 %s: %s", taskID[:8], s.Error())
}

func (p *Pipeline) Run(ctx context.Context, taskID string) {
	task, err := p.taskRepo.FindByID(taskID)
	if err != nil {
		logger.L.Errorf("任务不存在: %s", taskID)
		return
	}

	// 无论成功或失败，local 平台上传的临时文件都需要清理
	if task.Platform == "local" {
		defer p.cleanupLocalUpload(task.VideoURL)
	}

	result, err := p.execute(ctx, task)
	if err != nil {
		if se, ok := err.(*stepError); ok {
			p.failStep(taskID, se)
		} else {
			p.taskRepo.UpdateStatus(taskID, model.TaskStatusFailed, err.Error())
			logger.L.Errorf("任务执行失败 (task=%s): %v", taskID, err)
		}
		return
	}

	p.saveNote(taskID, result)
	p.taskRepo.UpdateStatus(taskID, model.TaskStatusSuccess, "")
	if p.statusBroadcast != nil {
		p.statusBroadcast(taskID, model.TaskStatusSuccess, "")
	}
	logger.L.Infof("任务 %s: SUCCESS", taskID[:8])
}

func (p *Pipeline) execute(ctx context.Context, task *model.Task) (*NoteResult, error) {
	var lastErr error
	for retry := 0; retry <= p.cfg.Task.MaxRetry; retry++ {
		if retry > 0 {
			logger.L.Infof("任务重试 %d/%d (task=%s)", retry, p.cfg.Task.MaxRetry, task.ID)
		}

		result, err := p.runOnce(ctx, task)
		if err == nil {
			return result, nil
		}
		if _, ok := err.(*stepError); ok {
			return nil, err
		}
		lastErr = err
	}

	return nil, &stepError{model.TaskStatusFailed, "任务失败", fmt.Errorf("已重试 %d 次: %w", p.cfg.Task.MaxRetry, lastErr)}
}

func (p *Pipeline) runOnce(ctx context.Context, task *model.Task) (*NoteResult, error) {
	dataDir := p.cfg.Storage.DataDir

	p.updateStatus(task.ID, model.TaskStatusParsing, "解析视频链接...")
	dl, err := downloader.New(ctx, task.Platform)
	if err != nil {
		return nil, &stepError{model.TaskStatusParsingFailed, "解析链接失败", err}
	}

	if task.Platform != "local" {
		ytPath := p.cfg.Tools.YtDlpPath
		checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if _, err := exec.CommandContext(checkCtx, ytPath, "--version").Output(); err != nil {
			cancel()
			return nil, &stepError{model.TaskStatusParsingFailed, "yt-dlp 不可用", fmt.Errorf("path=%s: %w — 请安装: pip install yt-dlp", ytPath, err)}
		}
		cancel()
	}

	cookieFile := p.prepareCookies(task.Platform)
	defer os.Remove(cookieFile)
	if cookieFile != "" {
		downloader.SetCookieOption(dl, &downloader.CookieOption{FilePath: cookieFile})
	}

	p.updateStatus(task.ID, model.TaskStatusDownloading, "获取视频字幕...")
	transcript, err := dl.DownloadSubtitles(ctx, task.VideoURL)
	if err != nil {
		logger.L.Warnf("获取字幕失败: %v", err)
	}

	skipDownload := transcript != nil

	var audioMeta *downloader.AudioMeta
	var workDir string

	if !skipDownload {
		audioMeta, err = dl.Download(ctx, task.VideoURL, task.Quality, dataDir, false, false)
		if err != nil {
			return nil, &stepError{model.TaskStatusDownloadingFailed, "下载失败", err}
		}
	} else {
		audioMeta, err = dl.Download(ctx, task.VideoURL, task.Quality, dataDir, false, true)
		if err != nil {
			return nil, &stepError{model.TaskStatusDownloadingFailed, "提取元信息失败", err}
		}
	}

	workDir = p.setupWorkDir(dataDir, audioMeta.VideoID)
	audioMeta = p.moveIntoWorkDir(audioMeta, workDir)

	task.VideoID = audioMeta.VideoID
	_ = p.taskRepo.UpdateVideoID(task.ID, audioMeta.VideoID)

	p.updateStatus(task.ID, model.TaskStatusDownloading, "下载完成")

	if transcript == nil {
		p.updateStatus(task.ID, model.TaskStatusTranscribing, "语音转录中...")

		audioFile := filepath.Join(workDir, audioMeta.VideoID+".mp3")
		if _, err := os.Stat(audioFile); os.IsNotExist(err) {
			return nil, &stepError{model.TaskStatusTranscribingFailed, "转录失败", fmt.Errorf("音频文件不存在: %s", audioFile)}
		}

		absAudioFile, _ := filepath.Abs(audioFile)
		transcriptResult, err := p.transcriberCl.Transcribe(ctx, absAudioFile)
		if err != nil {
			return nil, &stepError{model.TaskStatusTranscribingFailed, "转录失败", err}
		}

		transcript = &downloader.TranscriptResult{
			Language: transcriptResult.Language,
			FullText: transcriptResult.FullText,
		}
		for _, seg := range transcriptResult.Segments {
			transcript.Segments = append(transcript.Segments, downloader.TranscriptSegment{
				Start: seg.Start,
				End:   seg.End,
				Text:  seg.Text,
			})
		}

		transcriptPath := filepath.Join(workDir, audioMeta.VideoID+"_transcript.json")
		transcriptJSON, _ := json.MarshalIndent(transcript, "", "  ")
		if err := os.WriteFile(transcriptPath, transcriptJSON, 0644); err != nil {
			logger.L.Warnf("保存转录中间文件失败: %v", err)
		} else {
			logger.L.Infof("转录中间文件: %s", transcriptPath)
		}
	}

	p.updateStatus(task.ID, model.TaskStatusGenerating, "AI 生成笔记中...")

	provider, err := p.providerRepo.FindByID(task.ProviderID)
	if err != nil {
		return nil, &stepError{model.TaskStatusGeneratingFailed, "AI 生成失败", fmt.Errorf("模型供应商不存在: %w", err)}
	}

	rawInfoMap := make(map[string]interface{})
	json.Unmarshal([]byte(audioMeta.RawInfo), &rawInfoMap)
	description := ""
	if desc, ok := rawInfoMap["description"].(string); ok {
		description = desc
	}

	systemPrompt := llm.BuildVideoSummary(
		transcript.FullText,
		audioMeta.Title,
		description,
		task.Style,
	)

	messages := []llm.ChatMessage{
		{Role: "user", Content: systemPrompt},
	}

	markdown, err := p.llmClient.Chat(ctx, provider.BaseURL, provider.ApiKey, task.ModelName, messages)
	if err != nil {
		return nil, &stepError{model.TaskStatusGeneratingFailed, "AI 生成失败", fmt.Errorf("LLM 调用失败: %w", err)}
	}

	p.updateStatus(task.ID, model.TaskStatusPostProcessing, "后处理中...")

	markdown, noteTitle := postProcessMarkdown(markdown, task, audioMeta)

	return &NoteResult{
		Markdown:   markdown,
		Transcript: transcript,
		AudioMeta:  audioMeta,
		WorkDir:    workDir,
		NoteTitle:  noteTitle,
	}, nil
}

func (p *Pipeline) setupWorkDir(dataDir, videoID string) string {
	workDir := filepath.Join(dataDir, videoID)
	os.MkdirAll(workDir, 0755)
	return workDir
}

func (p *Pipeline) moveIntoWorkDir(meta *downloader.AudioMeta, workDir string) *downloader.AudioMeta {
	src := meta.FilePath
	dst := filepath.Join(workDir, meta.VideoID+".mp3")
	if src == dst {
		return meta
	}
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return meta
	}
	if err := os.Rename(src, dst); err != nil {
		logger.L.Warnf("移动音频文件失败: %v (src=%s, dst=%s)", err, src, dst)
		return meta
	}
	meta.FilePath = dst
	logger.L.Infof("音频已移入工作目录: %s", dst)
	return meta
}

func (p *Pipeline) cleanupLocalUpload(filePath string) {
	if filePath == "" {
		return
	}
	// 只清理 uploads 目录下的文件，避免误删用户本地文件
	uploadDir, _ := filepath.Abs(p.cfg.Storage.UploadDir)
	absFile, _ := filepath.Abs(filePath)
	if uploadDir != "" && strings.HasPrefix(absFile, uploadDir) {
		if err := os.Remove(absFile); err != nil {
			if !os.IsNotExist(err) {
				logger.L.Warnf("清理本地上传临时文件失败: %v (path=%s)", err, absFile)
			}
		} else {
			logger.L.Infof("已清理本地上传临时文件: %s", absFile)
		}
	}
}

func (p *Pipeline) updateStatus(taskID string, status model.TaskStatus, msg string) {
	p.taskRepo.UpdateStatus(taskID, status, "")
	if p.statusBroadcast != nil {
		p.statusBroadcast(taskID, status, msg)
	}
	logger.L.Infof("任务 %s: %s - %s", taskID[:8], status, msg)
}

func (p *Pipeline) saveNote(taskID string, result *NoteResult) {
	workDir := result.WorkDir

	mdPath := filepath.Join(workDir, taskID+".md")
	if err := os.WriteFile(mdPath, []byte(result.Markdown), 0644); err != nil {
		logger.L.Errorf("保存笔记文件失败: %v", err)
	} else {
		logger.L.Infof("笔记已保存: %s", mdPath)
	}

	transcriptJSON, _ := json.Marshal(result.Transcript)
	audioMetaJSON, _ := json.Marshal(result.AudioMeta)

	noteRecord := &model.NoteRecord{
		TaskID:     taskID,
		Markdown:   result.Markdown,
		Transcript: string(transcriptJSON),
		AudioMeta:  string(audioMetaJSON),
	}
	if err := p.noteRepo.Create(noteRecord); err != nil {
		logger.L.Errorf("保存笔记记录失败: %v", err)
	}

	if result.NoteTitle != "" {
		if err := p.taskRepo.UpdateName(taskID, result.NoteTitle); err != nil {
			logger.L.Warnf("保存笔记标题失败: %v", err)
		}
	}
}

type NoteResult struct {
	Markdown   string
	Transcript *downloader.TranscriptResult
	AudioMeta  *downloader.AudioMeta
	WorkDir    string
	NoteTitle  string
}

func (p *Pipeline) prepareCookies(platform string) string {
	if p.cookieRepo == nil {
		return ""
	}
	cookie, err := p.cookieRepo.FindByPlatform(platform)
	if err != nil || cookie.Content == "" {
		return ""
	}

	f, err := os.CreateTemp("", "aivideonote-cookies-*.txt")
	if err != nil {
		logger.L.Warnf("创建 cookie 临时文件失败: %v", err)
		return ""
	}

	content := cookie.Content
	if isHeaderCookieFormat(content) {
		content = convertHeaderToNetscape(content, platform)
	}

	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		logger.L.Warnf("写入 cookie 文件失败: %v", err)
		return ""
	}
	f.Close()
	logger.L.Infof("已加载 %s 平台 Cookie (%s)", platform, cookieFormatName(cookie.Content))
	return f.Name()
}

var platformDomains = map[string]string{
	"bilibili": ".bilibili.com",
}

func isHeaderCookieFormat(content string) bool {
	if strings.Contains(content, "\t") {
		return false
	}
	return strings.Contains(content, "; ") || strings.Contains(content, ";")
}

func cookieFormatName(content string) string {
	if isHeaderCookieFormat(content) {
		return "请求头格式→已转换"
	}
	return "Netscape 格式"
}

func convertHeaderToNetscape(content, platform string) string {
	if !strings.HasPrefix(content, "# Netscape") {
		var b strings.Builder
		b.WriteString("# Netscape HTTP Cookie File\n")
		b.WriteString("# Converted from browser header format by AIVideoNote\n\n")
		domain := platformDomains[platform]
		if domain == "" {
			domain = "." + platform + ".com"
		}
		pairs := strings.Split(content, ";")
		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}
			eqIdx := strings.Index(pair, "=")
			if eqIdx < 0 {
				continue
			}
			name := strings.TrimSpace(pair[:eqIdx])
			value := strings.TrimSpace(pair[eqIdx+1:])
			fmt.Fprintf(&b, "%s\tTRUE\t/\tFALSE\t0\t%s\t%s\n", domain, name, value)
		}
		return b.String()
	}
	return content
}

func postProcessMarkdown(markdown string, task *model.Task, meta *downloader.AudioMeta) (string, string) {
	trimmed := strings.TrimSpace(markdown)
	var noteTitle string

	lines := strings.SplitN(trimmed, "\n", 2)
	firstLine := strings.TrimSpace(lines[0])
	if strings.HasPrefix(firstLine, "# ") {
		noteTitle = strings.TrimPrefix(firstLine, "# ")
		noteTitle = strings.TrimSpace(noteTitle)
		sourceLine := "\n\n> 来源: " + task.VideoURL + "\n\n"
		if len(lines) > 1 {
			markdown = firstLine + sourceLine + lines[1]
		} else {
			markdown = firstLine + sourceLine
		}
	} else {
		noteTitle = meta.Title
		sourceLine := "> 来源: " + task.VideoURL + "\n\n"
		if noteTitle != "" {
			markdown = "# " + noteTitle + "\n\n" + sourceLine + trimmed
		} else {
			markdown = sourceLine + trimmed
		}
	}

	return markdown, noteTitle
}
