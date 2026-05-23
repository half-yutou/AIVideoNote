package handler

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/pkg/response"
)

type UploadHandler struct {
	cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择文件")
		return
	}

	uploadDir := h.cfg.Storage.UploadDir
	os.MkdirAll(uploadDir, 0755)

	ext := filepath.Ext(file.Filename)
	newName := uuid.New().String() + ext
	savePath := filepath.Join(uploadDir, newName)

	src, err := file.Open()
	if err != nil {
		response.InternalError(c, "打开文件失败")
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		response.InternalError(c, "创建文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		response.InternalError(c, "保存文件失败")
		return
	}

	response.Success(c, gin.H{
		"url":      "/uploads/" + newName,
		"filename": file.Filename,
	})
}
