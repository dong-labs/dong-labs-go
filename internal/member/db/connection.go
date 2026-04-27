package db

import (
	"sync"

	"github.com/dong-labs/think/internal/core/db"
)

const (
	NAME = "member"
)

type MemberDatabase struct {
	*db.Database
}

func NewMemberDatabase() *MemberDatabase {
	return &MemberDatabase{
		Database: db.NewDatabase(NAME),
	}
}

var (
	dbInstance *MemberDatabase
	once       sync.Once
)

func GetDB() *MemberDatabase {
	once.Do(func() {
		dbInstance = NewMemberDatabase()
	})
	return dbInstance
}
