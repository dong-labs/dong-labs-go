package db

import (
	"sync"
	"github.com/dong-labs/think/internal/core/db"
)

const NAME = "read"

type ReadDatabase struct { *db.Database }

func NewReadDatabase() *ReadDatabase {
	return &ReadDatabase{Database: db.NewDatabase(NAME)}
}

var dbInstance *ReadDatabase
var once sync.Once

func GetDB() *ReadDatabase {
	once.Do(func() { dbInstance = NewReadDatabase() })
	return dbInstance
}
