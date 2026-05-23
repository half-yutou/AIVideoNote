package media

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/aivideonote/backend/internal/config"
)

type ProbeInfo struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func Probe(ctx context.Context, filePath string) (*ProbeInfo, error) {
	cmd := exec.CommandContext(ctx,
		config.Cfg.Tools.FfprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe 执行失败: %w", err)
	}

	var info ProbeInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("解析 ffprobe 输出失败: %w", err)
	}

	return &info, nil
}

func ExtractAudio(ctx context.Context, videoPath, audioOutputPath string) error {
	cmd := exec.CommandContext(ctx,
		config.Cfg.Tools.FfmpegPath,
		"-i", videoPath,
		"-vn",
		"-acodec", "libmp3lame",
		"-q:a", "2",
		"-y",
		audioOutputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg 提取音频失败: %w\n%s", err, string(output))
	}

	return nil
}

func ExtractFrame(ctx context.Context, videoPath, outputPath string, timestamp int) error {
	cmd := exec.CommandContext(ctx,
		config.Cfg.Tools.FfmpegPath,
		"-ss", fmt.Sprintf("%d", timestamp),
		"-i", videoPath,
		"-frames:v", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg 截帧失败: %w\n%s", err, string(output))
	}

	return nil
}
