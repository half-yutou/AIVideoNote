package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/pkg/response"
)

// AllowedMediaExtensions 前后端统一的允许上传的媒体文件扩展名
var AllowedMediaExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".webm": true,
	".flv":  true,
	".mp3":  true,
	".wav":  true,
	".flac": true,
	".aac":  true,
	".ogg":  true,
	".m4a":  true,
}

// MaxUploadSize 最大上传文件大小 2GB
const MaxUploadSize = 2 << 30

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

	// 校验文件大小
	if file.Size > MaxUploadSize {
		response.BadRequest(c, "文件过大，最大支持 2GB")
		return
	}

	// 校验文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedMediaExtensions[ext] {
		allowed := make([]string, 0, len(AllowedMediaExtensions))
		for k := range AllowedMediaExtensions {
			allowed = append(allowed, k)
		}
		response.BadRequest(c, fmt.Sprintf("不支持的文件格式 %s，允许: %s", ext, strings.Join(allowed, ", ")))
		return
	}

	uploadDir := h.cfg.Storage.UploadDir
	os.MkdirAll(uploadDir, 0755)

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
		os.Remove(savePath)
		response.InternalError(c, "保存文件失败")
		return
	}

	// 返回绝对路径，供 pipeline 直接使用
	absPath, _ := filepath.Abs(savePath)

	response.Success(c, gin.H{
		"path":     absPath,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// AllowedExtensionsList 返回允许的扩展名列表（供 API 查询）
func (h *UploadHandler) AllowedFormats(c *gin.Context) {
	exts := make([]string, 0, len(AllowedMediaExtensions))
	for k := range AllowedMediaExtensions {
		exts = append(exts, k)
	}
	response.Success(c, gin.H{"formats": exts})
}
