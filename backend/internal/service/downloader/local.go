package downloader

import (
	"context"
	"os"
	"path/filepath"
)

type LocalDownloader struct{}

func NewLocalDownloader() *LocalDownloader {
	return &LocalDownloader{}
}

func (d *LocalDownloader) Name() string { return "local" }

func (d *LocalDownloader) Download(ctx context.Context, videoURL, quality, outputDir string, needVideo, skipDownload bool) (*AudioMeta, error) {
	filePath := videoURL

	ext := filepath.Ext(filePath)
	fileName := filepath.Base(filePath)
	nameWithoutExt := fileName[:len(fileName)-len(ext)]

	os.MkdirAll(outputDir, 0755)

	return &AudioMeta{
		FilePath: filePath,
		Title:    nameWithoutExt,
		Duration: 0,
		Platform: "local",
		VideoID:  nameWithoutExt,
		RawInfo:  "{}",
		VideoPath: func() string {
			if needVideo {
				return filePath
			}
			return ""
		}(),
	}, nil
}

func (d *LocalDownloader) DownloadVideo(ctx context.Context, videoURL, outputDir string) (string, error) {
	return videoURL, nil
}

func (d *LocalDownloader) DownloadSubtitles(ctx context.Context, videoURL string) (*TranscriptResult, error) {
	return nil, nil
}
