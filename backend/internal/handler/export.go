package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/export"
)

type ExportHandler struct {
	svc *export.Service
}

func NewExportHandler(svc *export.Service) *ExportHandler {
	return &ExportHandler{svc: svc}
}

func (h *ExportHandler) Export(c *gin.Context) {
	taskID := c.Param("id")
	format := c.DefaultQuery("format", "md")

	if format == "md" || format == "markdown" {
		content, err := h.svc.ExportMarkdown(taskID)
		if err != nil {
			response.Error(c, http.StatusNotFound, 404, "笔记不存在")
			return
		}
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=note.md")
		c.String(http.StatusOK, content)
		return
	}

	filePath, err := h.svc.ExportFile(taskID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "笔记文件不存在")
		return
	}

	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+taskID+".md")
	c.File(filePath)
}
