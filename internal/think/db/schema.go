// Package db 提供 think-cli 的数据库 Schema 管理
package db

import (
	"fmt"
	"github.com/dong-labs/think/internal/core/db"
)

const (
	// SCHEMA_VERSION 当前 schema 版本
	SCHEMA_VERSION = "1.0.0"
)

// ThinkSchemaManager 思咚咚 Schema 管理器
type ThinkSchemaManager struct {
	*db.SchemaManager
}

// NewThinkSchemaManager 创建新的 Schema 管理器
func NewThinkSchemaManager() *ThinkSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)

	return &ThinkSchemaManager{
		SchemaManager: schemaMgr,
	}
}

// InitSchema 初始化数据库 schema
func (s *ThinkSchemaManager) InitSchema() error {
	// 创建 thoughts 表
	if err := s.CreateTable(`
		CREATE TABLE IF NOT EXISTS thoughts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			tags TEXT,
			priority TEXT DEFAULT 'normal',
			status TEXT DEFAULT 'active',
			context TEXT,
			source_agent TEXT,
			note TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	// 创建索引
	if err := s.CreateIndex(`
		CREATE INDEX IF NOT EXISTS idx_thoughts_tags
		ON thoughts(tags)
	`); err != nil {
		return err
	}

	if err := s.CreateIndex(`
		CREATE INDEX IF NOT EXISTS idx_thoughts_priority
		ON thoughts(priority)
	`); err != nil {
		return err
	}

	if err := s.CreateIndex(`
		CREATE INDEX IF NOT EXISTS idx_thoughts_status
		ON thoughts(status)
	`); err != nil {
		return err
	}

	if err := s.CreateIndex(`
		CREATE INDEX IF NOT EXISTS idx_thoughts_created_at
		ON thoughts(created_at)
	`); err != nil {
		return err
	}

	return nil
}

// Initialize 初始化数据库（覆盖基类方法）
func (s *ThinkSchemaManager) Initialize() error {
	// 确保 meta 表存在（通过调用 GetMeta 触发）
	_, _ = s.SchemaManager.DB.GetMeta("_ensure_meta")

	// 检查是否已初始化
	initialized, err := s.IsInitialized()
	if err != nil {
		return fmt.Errorf("检查初始化状态失败: %w", err)
	}

	// 如果已经初始化且版本一致，跳过
	if initialized {
		storedVersion, _ := s.GetStoredVersion()
		if storedVersion == SCHEMA_VERSION {
			return nil
		}
	}

	// 执行 schema 初始化（调用当前类型的 InitSchema）
	if err := s.InitSchema(); err != nil {
		return fmt.Errorf("初始化 schema 失败: %w", err)
	}

	// 设置版本号
	if err := s.SetVersion(SCHEMA_VERSION); err != nil {
		return fmt.Errorf("设置版本号失败: %w", err)
	}

	return nil
}

// GetSchemaVersion 获取存储的 schema 版本
func GetSchemaVersion() (string, error) {
	mgr := NewThinkSchemaManager()
	return mgr.GetStoredVersion()
}

// SetSchemaVersion 设置 schema 版本
func SetSchemaVersion(version string) error {
	mgr := NewThinkSchemaManager()
	return mgr.SetVersion(version)
}

// IsInitialized 检查数据库是否已初始化
func IsInitialized() (bool, error) {
	mgr := NewThinkSchemaManager()
	return mgr.IsInitialized()
}

// InitDatabase 初始化数据库
func InitDatabase() error {
	mgr := NewThinkSchemaManager()
	return mgr.Initialize()
}
