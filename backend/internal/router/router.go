package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/handler"
	"github.com/aivideonote/backend/internal/service/chat"
	cookiesvc "github.com/aivideonote/backend/internal/service/cookie"
	"github.com/aivideonote/backend/internal/service/export"
	"github.com/aivideonote/backend/internal/service/group"
	"github.com/aivideonote/backend/internal/service/provider"
	"github.com/aivideonote/backend/internal/service/task"
)

type Deps struct {
	Cfg             *config.Config
	ProviderService *provider.Service
	TaskService     *task.Service
	GroupService    *group.Service
	ChatService     *chat.Service
	ExportService   *export.Service
	CookieService   *cookiesvc.Service
	Executor        *task.Executor
}

func Setup(r *gin.Engine, deps *Deps) {
	r.Static("/uploads", deps.Cfg.Storage.UploadDir)

	providerHandler := handler.NewProviderHandler(deps.ProviderService)
	taskHandler := handler.NewTaskHandler(deps.TaskService, deps.Executor)
	groupHandler := handler.NewGroupHandler(deps.GroupService)
	chatHandler := handler.NewChatHandler(deps.ChatService)
	cookieHandler := handler.NewCookieHandler(deps.CookieService)
	uploadHandler := handler.NewUploadHandler(deps.Cfg)
	exportHandler := handler.NewExportHandler(deps.ExportService)

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		provider := api.Group("/provider")
		{
			provider.POST("", providerHandler.Create)
			provider.GET("", providerHandler.GetAll)
			provider.GET("/:id", providerHandler.GetByID)
			provider.PUT("/:id", providerHandler.Update)
			provider.DELETE("/:id", providerHandler.Delete)
			provider.GET("/:id/test", providerHandler.TestConnection)
		}

		api.GET("/model/list", providerHandler.ListModels)

		api.POST("/note/generate", taskHandler.Generate)
		api.GET("/task/:id/status", taskHandler.GetStatus)
		api.GET("/task/list", taskHandler.GetList)
		api.DELETE("/task/:id", taskHandler.Delete)
		api.PUT("/task/:id/name", taskHandler.RenameTask)

		api.GET("/group/list", groupHandler.List)
		api.POST("/group/create", groupHandler.Create)
		api.PUT("/group/:id/name", groupHandler.Rename)
		api.DELETE("/group/:id", groupHandler.Delete)

		api.GET("/cookies", cookieHandler.GetAll)
		api.POST("/cookies", cookieHandler.Save)
		api.DELETE("/cookies/:platform", cookieHandler.Delete)

		api.POST("/chat/ask", chatHandler.Ask)
		api.GET("/chat/status", chatHandler.IndexStatus)

		api.POST("/upload", uploadHandler.Upload)

		api.GET("/export/:id", exportHandler.Export)
	}
}
