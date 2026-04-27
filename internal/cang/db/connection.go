package db

import (
	"sync"

	"github.com/dong-labs/think/internal/core/db"
)

const (
	NAME = "cang"
)

type CangDatabase struct {
	*db.Database
}

func NewCangDatabase() *CangDatabase {
	return &CangDatabase{
		Database: db.NewDatabase(NAME),
	}
}

var (
	dbInstance *CangDatabase
	once       sync.Once
)

func GetDB() *CangDatabase {
	once.Do(func() {
		dbInstance = NewCangDatabase()
	})
	return dbInstance
}
