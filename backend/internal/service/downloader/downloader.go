package downloader

import (
	"context"
	"fmt"

	"github.com/aivideonote/backend/internal/config"
)

type AudioMeta struct {
	FilePath  string `json:"file_path"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	CoverURL  string `json:"cover_url"`
	Platform  string `json:"platform"`
	VideoID   string `json:"video_id"`
	RawInfo   string `json:"raw_info"`
	VideoPath string `json:"video_path,omitempty"`
}

type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type TranscriptResult struct {
	Language string              `json:"language"`
	FullText string              `json:"full_text"`
	Segments []TranscriptSegment `json:"segments"`
}

type Downloader interface {
	Download(ctx context.Context, videoURL, quality, outputDir string, needVideo, skipDownload bool) (*AudioMeta, error)
	DownloadVideo(ctx context.Context, videoURL, outputDir string) (string, error)
	DownloadSubtitles(ctx context.Context, videoURL string) (*TranscriptResult, error)
	Name() string
}

// New 根据平台创建对应的下载器实例
// 当前仅支持 bilibili 和 local，但接口设计支持后续扩展其他平台
func New(ctx context.Context, platform string) (Downloader, error) {
	switch platform {
	case "bilibili":
		return NewBilibiliDownloader(), nil
	case "local":
		return NewLocalDownloader(), nil
	default:
		return nil, fmt.Errorf("不支持的平台: %s", platform)
	}
}

func ytDlpPath() string {
	return config.Cfg.Tools.YtDlpPath
}

// CookieOption 用于在每次调用时传递 cookie 文件路径，避免全局变量并发竞争
type CookieOption struct {
	FilePath string
}

func cookieArgsFrom(opt *CookieOption) []string {
	if opt == nil || opt.FilePath == "" {
		return nil
	}
	return []string{"--cookies", opt.FilePath}
}

// SetCookieOption 为 Downloader 设置 cookie 选项（线程安全，每个实例独立）
func SetCookieOption(dl Downloader, opt *CookieOption) {
	if s, ok := dl.(interface{ SetCookie(*CookieOption) }); ok {
		s.SetCookie(opt)
	}
}

// ErrDownloadFailed 下载失败错误
type downloadError struct {
	msg string
}

func (e *downloadError) Error() string { return e.msg }

func ErrDownloadFailed(msg string) error {
	return &downloadError{msg: msg}
}
