package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/middleware"
	"github.com/aivideonote/backend/internal/pkg/logger"
	"github.com/aivideonote/backend/internal/repository"
	"github.com/aivideonote/backend/internal/router"
	"github.com/aivideonote/backend/internal/service/chat"
	cookiesvc "github.com/aivideonote/backend/internal/service/cookie"
	"github.com/aivideonote/backend/internal/service/export"
	"github.com/aivideonote/backend/internal/service/group"
	"github.com/aivideonote/backend/internal/service/llm"
	provsvc "github.com/aivideonote/backend/internal/service/provider"
	"github.com/aivideonote/backend/internal/service/task"
	"github.com/aivideonote/backend/internal/service/transcriber"
	"github.com/gin-gonic/gin"
)

type repos struct {
	provider *repository.ProviderRepository
	task     *repository.TaskRepository
	note     *repository.NoteRepository
	group    *repository.GroupRepository
	cookie   *repository.CookieRepository
}

type services struct {
	provider *provsvc.Service
	task     *task.Service
	group    *group.Service
	chat     *chat.Service
	export   *export.Service
	cookie   *cookiesvc.Service
}

func main() {
	configPath, logLevel := parseFlags()
	initLogger(logLevel)
	defer logger.Sync()

	cfg := loadConfig(configPath)
	initDatabase(cfg)

	r := initGinEngine()
	llmClient, transcriberCl := initClients(cfg)
	transcriber.CheckTools(cfg.Tools.FfmpegPath, cfg.Tools.FfprobePath, cfg.Tools.YtDlpPath)
	repos := initRepos()
	svcs := initServices(cfg, repos, llmClient)
	executor := initTaskExecutor(cfg, repos, llmClient, transcriberCl)

	router.Setup(r, &router.Deps{
		Cfg:             cfg,
		ProviderService: svcs.provider,
		TaskService:     svcs.task,
		GroupService:    svcs.group,
		ChatService:     svcs.chat,
		ExportService:   svcs.export,
		CookieService:   svcs.cookie,
		Executor:        executor,
	})

	serve(r, fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}

func parseFlags() (configPath, logLevel string) {
	cp := flag.String("config", "config.yaml", "配置文件路径")
	ll := flag.String("log-level", "info", "日志级别: debug/info/warn/error")
	flag.Parse()
	return *cp, *ll
}

func initLogger(level string) {
	if err := logger.Init(level); err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig(path string) *config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		logger.L.Fatalf("加载配置失败: %v", err)
	}
	logger.L.Infof("配置加载成功")
	logger.L.Infof("  server: %s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.L.Infof("  数据目录: %s", cfg.Storage.DataDir)
	logger.L.Infof("  上传目录: %s", cfg.Storage.UploadDir)
	logger.L.Infof("  ffmpeg: %s", cfg.Tools.FfmpegPath)
	return cfg
}

func initDatabase(cfg *config.Config) {
	if err := database.Init(cfg.Database); err != nil {
		logger.L.Warnf("数据库连接失败（将以无数据库模式运行）: %v", err)
	}
}

func initGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	return r
}

func initClients(cfg *config.Config) (*llm.Client, *transcriber.Client) {
	return llm.NewClient(), transcriber.NewClient(cfg.PythonService)
}

func initRepos() *repos {
	r := &repos{
		provider: repository.NewProviderRepository(),
		task:     repository.NewTaskRepository(),
		note:     repository.NewNoteRepository(),
		group:    repository.NewGroupRepository(),
		cookie:   repository.NewCookieRepository(),
	}
	if _, err := r.group.EnsureDefault(); err != nil {
		logger.L.Warnf("创建默认分组失败: %v", err)
	}
	return r
}

func initServices(cfg *config.Config, r *repos, llmClient *llm.Client) *services {
	return &services{
		provider: provsvc.NewService(r.provider, llmClient),
		task:     task.NewService(r.task, r.note, r.group),
		group:    group.NewService(r.group),
		chat:     chat.NewService(r.note, r.provider, llmClient),
		export:   export.NewService(cfg, r.note, r.task),
		cookie:   cookiesvc.NewService(r.cookie),
	}
}

func initTaskExecutor(cfg *config.Config, r *repos, llmClient *llm.Client, transcriberCl *transcriber.Client) *task.Executor {
	pipeline := task.NewPipeline(
		cfg,
		r.task,
		r.note,
		r.provider,
		r.cookie,
		llmClient,
		transcriberCl,
		nil,
	)
	return task.NewExecutor(cfg, pipeline)
}

func serve(r *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.L.Infof("AIVideoNote 后端启动: http://%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatalf("服务启动失败: %v", err)
		}
	}()

	<-ctx.Done()
	logger.L.Info("收到关闭信号，正在优雅关闭...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.L.Errorf("服务关闭失败: %v", err)
	}

	logger.L.Info("服务已关闭")
}
