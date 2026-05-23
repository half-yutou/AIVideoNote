package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aivideonote/backend/internal/pkg/crypto"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig `mapstructure:"server"`
	Tools         ToolsConfig
	PythonService PythonServiceConfig
	Database      DatabaseConfig `mapstructure:"database"`
	Storage       StorageConfig
	LLM           LLMConfig
	Task          TaskConfig `mapstructure:"task"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type ToolsConfig struct {
	FfmpegPath  string `mapstructure:"ffmpeg_path"`
	FfprobePath string `mapstructure:"ffprobe_path"`
	YtDlpPath   string `mapstructure:"yt_dlp_path"`
}

type PythonServiceConfig struct {
	URL               string `mapstructure:"url"`
	TranscribeTimeout int    `mapstructure:"transcribe_timeout"`
}

type DatabaseConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type StorageConfig struct {
	DataDir   string `mapstructure:"data_dir"`
	UploadDir string `mapstructure:"upload_dir"`
}

type LLMConfig struct {
	DefaultBaseURL string `mapstructure:"default_base_url"`
}

type TaskConfig struct {
	MaxConcurrent int `mapstructure:"max_concurrent"`
	MaxRetry      int `mapstructure:"max_retry"`
	RetryDelay    int `mapstructure:"retry_delay"`
}

var Cfg *Config

func Load(configPath string) (*Config, error) {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("..")
		if exePath, err := os.Executable(); err == nil {
			v.AddConfigPath(filepath.Dir(exePath))
		}
	}

	v.SetEnvPrefix("VN")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	cfg.ensureDefaults()

	if err := cfg.initCrypto(); err != nil {
		return nil, fmt.Errorf("初始化加密模块失败: %w", err)
	}

	if err := cfg.ensureDirs(); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	Cfg = &cfg
	return &cfg, nil
}

func (c *Config) ensureDefaults() {
	c.applyEnvOverrides()
	c.applyEnvSections()

	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Tools.FfmpegPath == "" {
		c.Tools.FfmpegPath = "ffmpeg"
	}
	if c.Tools.FfprobePath == "" {
		c.Tools.FfprobePath = "ffprobe"
	}
	if c.Tools.YtDlpPath == "" {
		c.Tools.YtDlpPath = "yt-dlp"
	}
	if c.PythonService.URL == "" {
		transcriberPort := "9090"
		if p := os.Getenv("TRANSCRIBER_PORT"); p != "" {
			transcriberPort = p
		}
		c.PythonService.URL = "http://127.0.0.1:" + transcriberPort
	}
	if c.PythonService.TranscribeTimeout == 0 {
		c.PythonService.TranscribeTimeout = 600
	}
	if c.Database.FilePath == "" {
		c.Database.FilePath = "./data/aivideonote.db"
	}
	if c.Storage.DataDir == "" {
		c.Storage.DataDir = "./data"
	}
	if c.Storage.UploadDir == "" {
		c.Storage.UploadDir = "./data/uploads"
	}
	if c.LLM.DefaultBaseURL == "" {
		c.LLM.DefaultBaseURL = "https://api.openai.com/v1"
	}
	if c.Task.MaxConcurrent == 0 {
		c.Task.MaxConcurrent = 3
	}
	if c.Task.MaxRetry == 0 {
		c.Task.MaxRetry = 3
	}
	if c.Task.RetryDelay == 0 {
		c.Task.RetryDelay = 30
	}
}

func (c *Config) initCrypto() error {
	return crypto.EnsureKey(os.Getenv("ENCRYPTION_KEY"))
}

func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("BACKEND_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.Server.Port = p
		}
	}
}

func (c *Config) applyEnvSections() {
	if v := os.Getenv("FFMPEG_PATH"); v != "" {
		c.Tools.FfmpegPath = v
	}
	if v := os.Getenv("FFPROBE_PATH"); v != "" {
		c.Tools.FfprobePath = v
	}
	if v := os.Getenv("YT_DLP_PATH"); v != "" {
		c.Tools.YtDlpPath = v
	}
	if v := os.Getenv("TRANSCRIBE_TIMEOUT"); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			c.PythonService.TranscribeTimeout = t
		}
	}
	if v := os.Getenv("DATA_DIR"); v != "" {
		c.Storage.DataDir = v
	}
	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		c.Storage.UploadDir = v
	}
	if v := os.Getenv("LLM_DEFAULT_BASE_URL"); v != "" {
		c.LLM.DefaultBaseURL = v
	}
}

func (c *Config) ensureDirs() error {
	dirs := []string{
		c.Storage.DataDir,
		c.Storage.UploadDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
