package handler

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/chat"
	"github.com/aivideonote/backend/internal/service/llm"
)

type ChatHandler struct {
	svc *chat.Service
}

func NewChatHandler(svc *chat.Service) *ChatHandler {
	return &ChatHandler{svc: svc}
}

type askRequest struct {
	TaskID     string    `json:"task_id" binding:"required"`
	Question   string    `json:"question" binding:"required"`
	History    []chatMsg `json:"history"`
	ProviderID string    `json:"provider_id" binding:"required"`
	ModelName  string    `json:"model_name" binding:"required"`
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (h *ChatHandler) Ask(c *gin.Context) {
	var req askRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	markdownSnippets, transcriptSnippets, err := h.svc.BuildContext(req.TaskID, req.Question)
	if err != nil {
		response.BadRequest(c, "笔记不存在，请先生成笔记")
		return
	}

	ctxPrompt := llm.BuildChatContext(markdownSnippets, transcriptSnippets, req.Question)

	messages := []llm.ChatMessage{
		{Role: "system", Content: "你是一个视频笔记问答助手。请基于提供的笔记和转录内容回答问题。用中文回答，保持简洁准确。"},
	}
	for _, msg := range req.History {
		messages = append(messages, llm.ChatMessage{Role: msg.Role, Content: msg.Content})
	}
	messages = append(messages, llm.ChatMessage{Role: "user", Content: ctxPrompt})

	acceptSSE := strings.Contains(c.GetHeader("Accept"), "text/event-stream")

	if acceptSSE {
		ch, errCh, err := h.svc.AskStream(c.Request.Context(), req.ProviderID, req.ModelName, messages)
		if err != nil {
			response.BadRequest(c, "模型供应商不存在")
			return
		}

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		c.Stream(func(w io.Writer) bool {
			select {
			case content, ok := <-ch:
				if !ok {
					return false
				}
				c.SSEvent("message", content)
				return true
			case err, ok := <-errCh:
				if ok && err != nil {
					c.SSEvent("error", err.Error())
				}
				return false
			}
		})
		return
	}

	answer, err := h.svc.Ask(c.Request.Context(), req.ProviderID, req.ModelName, messages)
	if err != nil {
		response.InternalError(c, "问答失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"answer": answer})
}

func (h *ChatHandler) IndexStatus(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		response.BadRequest(c, "请指定 task_id")
		return
	}

	indexed := h.svc.IndexStatus(taskID)
	response.Success(c, gin.H{"task_id": taskID, "indexed": indexed})
}
