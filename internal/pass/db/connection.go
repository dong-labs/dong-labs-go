package db

import (
	"sync"

	"github.com/dong-labs/think/internal/core/db"
)

const (
	NAME = "pass"
)

type PassDatabase struct {
	*db.Database
}

func NewPassDatabase() *PassDatabase {
	return &PassDatabase{
		Database: db.NewDatabase(NAME),
	}
}

var (
	dbInstance *PassDatabase
	once       sync.Once
)

func GetDB() *PassDatabase {
	once.Do(func() {
		dbInstance = NewPassDatabase()
	})
	return dbInstance
}
