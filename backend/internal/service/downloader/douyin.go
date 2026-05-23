package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GenericYtDlpDownloader struct {
	platform string
}

func NewGenericYtDlpDownloader(platform string) *GenericYtDlpDownloader {
	return &GenericYtDlpDownloader{platform: platform}
}

func (d *GenericYtDlpDownloader) Name() string { return d.platform }

func (d *GenericYtDlpDownloader) Download(ctx context.Context, videoURL, quality, outputDir string, needVideo, skipDownload bool) (*AudioMeta, error) {
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
	meta.Platform = d.platform
	return meta, nil
}

func (d *GenericYtDlpDownloader) DownloadVideo(ctx context.Context, videoURL, outputDir string) (string, error) {
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

func (d *GenericYtDlpDownloader) DownloadSubtitles(ctx context.Context, videoURL string) (*TranscriptResult, error) {
	return nil, nil
}

func (d *GenericYtDlpDownloader) metaOnly(ctx context.Context, videoURL, outputDir string) (*AudioMeta, error) {
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
	meta.Platform = d.platform
	return meta, nil
}

func parseSRT(_ string) (*TranscriptResult, error) {
	return nil, fmt.Errorf("SRT 解析未实现")
}

type downloadError struct {
	msg string
}

func (e *downloadError) Error() string { return e.msg }

func ErrDownloadFailed(msg string) error {
	return &downloadError{msg: msg}
}
