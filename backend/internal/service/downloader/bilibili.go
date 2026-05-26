package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aivideonote/backend/internal/pkg/logger"
)

type BilibiliDownloader struct {
	cookie *CookieOption
}

func NewBilibiliDownloader() *BilibiliDownloader {
	return &BilibiliDownloader{}
}

func (d *BilibiliDownloader) SetCookie(opt *CookieOption) {
	d.cookie = opt
}

func (d *BilibiliDownloader) Name() string { return "bilibili" }

func (d *BilibiliDownloader) Download(ctx context.Context, videoURL, quality, outputDir string, needVideo, skipDownload bool) (*AudioMeta, error) {
	if skipDownload {
		return d.metaOnly(ctx, videoURL, outputDir)
	}

	os.MkdirAll(outputDir, 0755)

	format := "bestaudio[ext=m4a]/bestaudio/best"

	args := []string{
		"-f", format,
		"-o", "%(id)s.%(ext)s",
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", qualityToABR(quality),
		"--no-playlist",
		"--print-json",
		"--no-simulate",
		videoURL,
	}

	if needVideo {
		args[0] = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
		args = append(args[:5], args[7:]...)
	}

	meta, err := d.runYtDlp(ctx, args, outputDir)
	if err != nil {
		return nil, err
	}
	meta.Platform = "bilibili"
	return meta, nil
}

func (d *BilibiliDownloader) DownloadVideo(ctx context.Context, videoURL, outputDir string) (string, error) {
	os.MkdirAll(outputDir, 0755)

	args := []string{
		"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
		"-o", "%(id)s.%(ext)s",
		"--no-playlist",
		"--print-json",
		videoURL,
	}

	info, err := d.runYtDlp(ctx, args, outputDir)
	if err != nil {
		return "", err
	}

	videoPath := filepath.Join(outputDir, info.VideoID+".mp4")
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		videoPath = filepath.Join(outputDir, info.VideoID+".mkv")
		if _, err := os.Stat(videoPath); os.IsNotExist(err) {
			videoPath = filepath.Join(outputDir, info.VideoID+".webm")
		}
	}

	return videoPath, nil
}

func (d *BilibiliDownloader) DownloadSubtitles(ctx context.Context, videoURL string) (*TranscriptResult, error) {
	return nil, nil
}

func (d *BilibiliDownloader) metaOnly(ctx context.Context, videoURL, outputDir string) (*AudioMeta, error) {
	args := []string{
		"--print-json",
		"--skip-download",
		"--no-playlist",
		videoURL,
	}
	return d.runYtDlp(ctx, args, outputDir)
}

func (d *BilibiliDownloader) runYtDlp(ctx context.Context, args []string, outputDir string) (*AudioMeta, error) {
	cmdArgs := append([]string{
		"--referer", "https://www.bilibili.com",
	}, args...)
	cmdArgs = append(cmdArgs, cookieArgsFrom(d.cookie)...)

	logger.L.Infof("yt-dlp %s", cmdArgs)
	cmd := exec.CommandContext(ctx, ytDlpPath(), cmdArgs...)
	cmd.Dir = outputDir

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("yt-dlp 执行失败: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("yt-dlp 执行失败: %w", err)
	}

	return parseYtDlpOutput(output, outputDir)
}

func parseYtDlpOutput(output []byte, outputDir string) (*AudioMeta, error) {
	lines := strings.Split(string(output), "\n")
	var lastJSON string
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if len(trimmed) > 0 && trimmed[0] == '{' {
			lastJSON = trimmed
			break
		}
	}

	if lastJSON == "" {
		return nil, fmt.Errorf("yt-dlp 未输出 JSON 信息")
	}

	var info struct {
		ID        string  `json:"id"`
		Title     string  `json:"title"`
		Duration  float64 `json:"duration"`
		Thumbnail string  `json:"thumbnail"`
	}

	if err := json.Unmarshal([]byte(lastJSON), &info); err != nil {
		return nil, fmt.Errorf("解析 yt-dlp 输出失败: %w", err)
	}

	return &AudioMeta{
		VideoID:  info.ID,
		Title:    info.Title,
		Duration: int(info.Duration),
		CoverURL: info.Thumbnail,
		FilePath: filepath.Join(outputDir, info.ID+".mp3"),
		RawInfo:  lastJSON,
	}, nil
}

func qualityToABR(quality string) string {
	switch quality {
	case "fast":
		return "64"
	case "medium":
		return "128"
	case "slow":
		return "192"
	default:
		return "128"
	}
}
