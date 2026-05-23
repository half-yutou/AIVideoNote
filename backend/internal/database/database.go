package database

import (
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/logger"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) error {
	var logLevel gormlogger.LogLevel = gormlogger.Warn

	db, err := gorm.Open(sqlite.Open(cfg.FilePath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	if err := db.AutoMigrate(
		&model.Task{},
		&model.LLMProvider{},
		&model.NoteRecord{},
		&model.Group{},
		&model.PlatformCookie{},
	); err != nil {
		return err
	}

	DB = db
	logger.L.Infof("SQLite 数据库已连接: %s (WAL模式)", cfg.FilePath)
	return nil
}
