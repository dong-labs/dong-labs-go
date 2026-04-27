// Package db log 数据库 Schema 管理
package db

import (
	"fmt"
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "1.0.0"
)

// LogSchemaManager 日志 Schema 管理器
type LogSchemaManager struct {
	*db.SchemaManager
}

// NewLogSchemaManager 创建新的 Schema 管理器
func NewLogSchemaManager() *LogSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &LogSchemaManager{
		SchemaManager: schemaMgr,
	}
}

// InitSchema 初始化数据库 schema
func (s *LogSchemaManager) InitSchema() error {
	if err := s.CreateTable(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			log_group TEXT DEFAULT 'default',
			date TEXT,
			tags TEXT DEFAULT '',
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_logs_group ON logs(log_group)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_logs_date ON logs(date)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_logs_tags ON logs(tags)`)

	return nil
}

// Initialize 初始化数据库（覆盖基类方法）
func (s *LogSchemaManager) Initialize() error {
	_, _ = s.SchemaManager.DB.GetMeta("_ensure_meta")

	initialized, err := s.IsInitialized()
	if err != nil {
		return fmt.Errorf("检查初始化状态失败: %w", err)
	}

	if initialized {
		storedVersion, _ := s.GetStoredVersion()
		if storedVersion == SCHEMA_VERSION {
			return nil
		}
	}

	if err := s.InitSchema(); err != nil {
		return fmt.Errorf("初始化 schema 失败: %w", err)
	}

	return s.SetVersion(SCHEMA_VERSION)
}

// InitDatabase 初始化数据库
func InitDatabase() error {
	mgr := NewLogSchemaManager()
	return mgr.Initialize()
}
