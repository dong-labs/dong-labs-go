package db

import (
	"sync"

	"github.com/dong-labs/think/internal/core/db"
)

const (
	NAME = "timeline"
)

type TimelineDatabase struct {
	*db.Database
}

func NewTimelineDatabase() *TimelineDatabase {
	return &TimelineDatabase{
		Database: db.NewDatabase(NAME),
	}
}

var (
	dbInstance *TimelineDatabase
	once       sync.Once
)

func GetDB() *TimelineDatabase {
	once.Do(func() {
		dbInstance = NewTimelineDatabase()
	})
	return dbInstance
}
