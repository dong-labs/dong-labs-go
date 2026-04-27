// Package db log 数据库连接管理
package db

import (
	"github.com/dong-labs/think/internal/core/db"
	"sync"
)

const (
	NAME = "log"
)

// LogDatabase 日志数据库
type LogDatabase struct {
	*db.Database
}

// NewLogDatabase 创建新的数据库实例
func NewLogDatabase() *LogDatabase {
	return &LogDatabase{
		Database: db.NewDatabase(NAME),
	}
}

var (
	dbInstance *LogDatabase
	once      sync.Once
)

// GetDB 获取数据库实例（单例）
func GetDB() *LogDatabase {
	once.Do(func() {
		dbInstance = NewLogDatabase()
	})
	return dbInstance
}
