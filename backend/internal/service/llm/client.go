package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aivideonote/backend/internal/pkg/logger"
)

type Client struct {
	httpClient *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream"`
}

type chatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type modelListResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

func (c *Client) ListModels(ctx context.Context, baseURL, apiKey string) ([]string, error) {
	url := strings.TrimRight(baseURL, "/") + "/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求模型列表失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取模型列表失败 (status=%d): %s", resp.StatusCode, string(body))
	}

	var result modelListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析模型列表失败: %w", err)
	}

	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}
	return models, nil
}

func (c *Client) Chat(ctx context.Context, baseURL, apiKey, model string, messages []ChatMessage) (string, error) {
	url := strings.TrimRight(baseURL, "/") + "/chat/completions"

	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
		Stream:      false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	logger.L.Infof("LLM Chat: model=%s base_url=%s", model, baseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM返回错误 (status=%d): %s", resp.StatusCode, string(body))
	}

	var result chatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析LLM响应失败: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("LLM错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("LLM返回空响应")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *Client) ChatStream(ctx context.Context, baseURL, apiKey, model string, messages []ChatMessage) (<-chan string, <-chan error) {
	ch := make(chan string, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		url := strings.TrimRight(baseURL, "/") + "/chat/completions"

		reqBody := chatRequest{
			Model:       model,
			Messages:    messages,
			Temperature: 0.7,
			MaxTokens:   4096,
			Stream:      true,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			errCh <- err
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errCh <- fmt.Errorf("LLM返回错误 (status=%d): %s", resp.StatusCode, string(body))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()

			// SSE 格式: "data: {...}" 或 "data: [DONE]"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var streamChunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamChunk); err != nil {
				continue // 跳过无法解析的行
			}

			if len(streamChunk.Choices) > 0 && streamChunk.Choices[0].Delta.Content != "" {
				ch <- streamChunk.Choices[0].Delta.Content
			}
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			errCh <- err
		}
	}()

	return ch, errCh
}
