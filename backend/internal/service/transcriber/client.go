package transcriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/pkg/logger"
)

var toolInstallHints = map[string][]string{
	"ffmpeg": {
		"  Windows: winget install ffmpeg",
		"           或下载: https://www.gyan.dev/ffmpeg/builds/",
		"  macOS:   brew install ffmpeg",
		"  Linux:   sudo apt install ffmpeg",
	},
	"ffprobe": {
		"  ffprobe 随 ffmpeg 一起安装，安装 ffmpeg 即可",
	},
	"yt-dlp": {
		"  pip install yt-dlp",
		"  或 winget install yt-dlp",
		"  或下载 .exe: https://github.com/yt-dlp/yt-dlp/releases",
	},
}

func CheckTools(ffmpegPath, ffprobePath, ytDlpPath string) {
	tools := map[string]string{
		"ffmpeg":  ffmpegPath,
		"ffprobe": ffprobePath,
		"yt-dlp":  ytDlpPath,
	}

	for name, path := range tools {
		if _, err := exec.LookPath(path); err != nil {
			logger.L.Warnf("⚠ 未找到 %s，视频处理和转录功能将不可用", name)
			for _, hint := range toolInstallHints[name] {
				logger.L.Warnf("%s", hint)
			}
		} else {
			logger.L.Infof("✓ %s: %s", name, path)
		}
	}
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type Result struct {
	Language string    `json:"language"`
	FullText string    `json:"full_text"`
	Segments []Segment `json:"segments"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg config.PythonServiceConfig) *Client {
	return &Client{
		baseURL: cfg.URL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TranscribeTimeout) * time.Second,
		},
	}
}

type transcribeRequest struct {
	AudioPath string `json:"audio_path"`
}

func (c *Client) Transcribe(ctx context.Context, audioFilePath string) (*Result, error) {
	url := c.baseURL + "/api/v1/transcribe"

	reqBody := transcribeRequest{AudioPath: audioFilePath}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建转录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	logger.L.Infof("调用 Python 转录服务: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("转录服务请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("转录服务返回错误 (status=%d): %s", resp.StatusCode, string(body))
	}

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析转录结果失败: %w", err)
	}

	return &result, nil
}
