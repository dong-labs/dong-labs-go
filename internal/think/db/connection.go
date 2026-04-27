// Package db 提供 think-cli 的数据库连接管理
package db

import (
	"github.com/dong-labs/think/internal/core/db"
	"sync"
)

const (
	// NAME CLI 名称
	NAME = "think"
)

// ThinkDatabase 思咚咚数据库
type ThinkDatabase struct {
	*db.Database
}

// NewThinkDatabase 创建新的数据库实例
func NewThinkDatabase() *ThinkDatabase {
	return &ThinkDatabase{
		Database: db.NewDatabase(NAME),
	}
}

// GetName 实现 Database 接口
func (d *ThinkDatabase) GetName() string {
	return NAME
}

var (
	dbInstance *ThinkDatabase
	once      sync.Once
)

// GetDB 获取数据库实例（单例）
func GetDB() *ThinkDatabase {
	once.Do(func() {
		dbInstance = NewThinkDatabase()
	})
	return dbInstance
}
