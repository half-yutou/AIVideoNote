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

func New(ctx context.Context, platform string) (Downloader, error) {
	switch platform {
	case "bilibili":
		return NewBilibiliDownloader(), nil
	case "youtube":
		return NewYoutubeDownloader(), nil
	case "local":
		return NewLocalDownloader(), nil
	case "douyin", "tiktok":
		return NewGenericYtDlpDownloader("douyin"), nil
	case "kuaishou":
		return NewGenericYtDlpDownloader("kuaishou"), nil
	default:
		return nil, fmt.Errorf("不支持的平台: %s", platform)
	}
}

func videoOutputPath(outputDir, videoID string) string {
	return fmt.Sprintf("%s/%s_video.mp4", outputDir, videoID)
}

func audioOutputPath(outputDir, videoID string) string {
	return fmt.Sprintf("%s/%s_audio.mp3", outputDir, videoID)
}

func ytDlpPath() string {
	return config.Cfg.Tools.YtDlpPath
}

var cookieFile string

func SetCookieFile(path string) {
	cookieFile = path
}

func cookieArgs() []string {
	if cookieFile == "" {
		return nil
	}
	return []string{"--cookies", cookieFile}
}
