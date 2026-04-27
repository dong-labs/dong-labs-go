// Package db 提供数据库 Schema 管理基类
package db

import (
	"fmt"
)

// VERSION_KEY Schema 版本在 meta 表中的 key
const VERSION_KEY = "schema_version"

// SchemaManager Schema 管理接口
//
// 子类需要实现：
// - InitSchema() 初始化数据库 schema
type SchemaManager struct {
	DB             *Database
	currentVersion  string
}

// NewSchemaManager 创建新的 Schema 管理器
func NewSchemaManager(db *Database, currentVersion string) *SchemaManager {
	return &SchemaManager{
		DB:             db,
		currentVersion: currentVersion,
	}
}

// InitSchema 初始化数据库 schema（子类必须实现）
func (s *SchemaManager) InitSchema() error {
	return fmt.Errorf("InitSchema 必须由子类实现")
}

// GetStoredVersion 获取存储的 schema 版本
func (s *SchemaManager) GetStoredVersion() (string, error) {
	s.DB.ensureMetaTable()
	return s.DB.GetMeta(VERSION_KEY)
}

// SetVersion 设置 schema 版本
func (s *SchemaManager) SetVersion(version string) error {
	return s.DB.SetMeta(VERSION_KEY, version)
}

// GetCurrentVersion 获取当前版本
func (s *SchemaManager) GetCurrentVersion() string {
	return s.currentVersion
}

// IsInitialized 检查数据库是否已初始化
func (s *SchemaManager) IsInitialized() (bool, error) {
	version, err := s.GetStoredVersion()
	if err != nil {
		return false, err
	}
	return version != "", nil
}

// RequiresMigration 检查是否需要迁移
func (s *SchemaManager) RequiresMigration() (bool, error) {
	stored, err := s.GetStoredVersion()
	if err != nil {
		return false, err
	}
	if stored == "" {
		return false, nil
	}
	return stored != s.currentVersion, nil
}

// GetVersionDelta 获取版本差异
// 返回 (存储的版本, 当前版本)
func (s *SchemaManager) GetVersionDelta() (string, string, error) {
	stored, err := s.GetStoredVersion()
	if err != nil {
		return "", "", err
	}
	return stored, s.currentVersion, nil
}

// Initialize 初始化数据库
// 创建 meta 表、执行 InitSchema、设置版本号
func (s *SchemaManager) Initialize() error {
	// 确保 meta 表存在
	s.DB.ensureMetaTable()

	// 检查是否已初始化
	initialized, err := s.IsInitialized()
	if err != nil {
		return fmt.Errorf("检查初始化状态失败: %w", err)
	}

	// 如果已经初始化且版本一致，跳过
	if initialized {
		storedVersion, _ := s.GetStoredVersion()
		if storedVersion == s.currentVersion {
			return nil
		}
	}

	// 执行 schema 初始化
	if err := s.InitSchema(); err != nil {
		return fmt.Errorf("初始化 schema 失败: %w", err)
	}

	// 设置版本号
	if err := s.SetVersion(s.currentVersion); err != nil {
		return fmt.Errorf("设置版本号失败: %w", err)
	}

	return nil
}

// MustInitialize 初始化数据库，失败时 panic
func (s *SchemaManager) MustInitialize() {
	if err := s.Initialize(); err != nil {
		panic(err)
	}
}

// CreateTable 创建表（如果不存在）
func (s *SchemaManager) CreateTable(sql string) error {
	_, err := s.DB.Exec(sql)
	return err
}

// CreateIndex 创建索引（如果不存在）
func (s *SchemaManager) CreateIndex(sql string) error {
	_, err := s.DB.Exec(sql)
	return err
}

// DropTable 删除表
func (s *SchemaManager) DropTable(tableName string) error {
	return s.DB.DropTable(tableName)
}

// TableExists 检查表是否存在
func (s *SchemaManager) TableExists(tableName string) (bool, error) {
	return s.DB.TableExists(tableName)
}
