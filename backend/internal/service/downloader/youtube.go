package downloader

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aivideonote/backend/internal/pkg/logger"
)

type YoutubeDownloader struct{}

func NewYoutubeDownloader() *YoutubeDownloader {
	return &YoutubeDownloader{}
}

func (d *YoutubeDownloader) Name() string { return "youtube" }

func (d *YoutubeDownloader) Download(ctx context.Context, videoURL, quality, outputDir string, needVideo, skipDownload bool) (*AudioMeta, error) {
	if skipDownload {
		return d.metaOnly(ctx, videoURL, outputDir)
	}

	os.MkdirAll(outputDir, 0755)

	args := []string{
		"-f", "bestaudio[ext=m4a]/bestaudio/best",
		"-o", "%(id)s.%(ext)s",
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", qualityToABR(quality),
		"--no-playlist",
		"--print-json",
		videoURL,
	}

	cmd := exec.CommandContext(ctx, ytDlpPath(), append(args, cookieArgs()...)...)
	cmd.Dir = outputDir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, ErrDownloadFailed(string(exitErr.Stderr))
		}
		return nil, ErrDownloadFailed(err.Error())
	}

	meta, err := parseYtDlpOutput(output, outputDir)
	if err != nil {
		return nil, err
	}
	meta.Platform = "youtube"
	return meta, nil
}

func (d *YoutubeDownloader) DownloadVideo(ctx context.Context, videoURL, outputDir string) (string, error) {
	os.MkdirAll(outputDir, 0755)

	args := []string{
		"-f", "best[ext=mp4]/best",
		"-o", "%(id)s.%(ext)s",
		"--no-playlist",
		"--print-json",
		videoURL,
	}

	cmd := exec.CommandContext(ctx, ytDlpPath(), append(args, cookieArgs()...)...)
	cmd.Dir = outputDir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", ErrDownloadFailed(string(exitErr.Stderr))
		}
		return "", ErrDownloadFailed(err.Error())
	}

	info, err := parseYtDlpOutput(output, outputDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(outputDir, info.VideoID+".mp4"), nil
}

func (d *YoutubeDownloader) DownloadSubtitles(ctx context.Context, videoURL string) (*TranscriptResult, error) {
	tmpDir, err := os.MkdirTemp("", "yt_subs_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	args := []string{
		"--write-auto-subs",
		"--sub-langs", "zh-Hans,zh,en",
		"--convert-subs", "srt",
		"--skip-download",
		"-o", "%(id)s.%(ext)s",
		"--no-playlist",
		"--print-json",
		videoURL,
	}

	cmd := exec.CommandContext(ctx, ytDlpPath(), append(args, cookieArgs()...)...)
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		logger.L.Warnf("获取 YouTube 字幕失败: %v", err)
		return nil, nil
	}

	info, err := parseYtDlpOutput(output, tmpDir)
	if err != nil || info == nil {
		return nil, nil
	}

	srtFiles, _ := filepath.Glob(filepath.Join(tmpDir, "*.srt"))
	if len(srtFiles) == 0 {
		return nil, nil
	}

	return parseSRT(srtFiles[0])
}

func (d *YoutubeDownloader) metaOnly(ctx context.Context, videoURL, outputDir string) (*AudioMeta, error) {
	args := []string{
		"--print-json",
		"--skip-download",
		"--no-playlist",
		videoURL,
	}

	cmd := exec.CommandContext(ctx, ytDlpPath(), append(args, cookieArgs()...)...)
	cmd.Dir = outputDir
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrDownloadFailed(err.Error())
	}

	meta, err := parseYtDlpOutput(output, outputDir)
	if err != nil {
		return nil, err
	}
	meta.Platform = "youtube"
	return meta, nil
}
