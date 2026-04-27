package db

import (
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "1.0.0"
)

type TimelineSchemaManager struct {
	*db.SchemaManager
}

func NewTimelineSchemaManager() *TimelineSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &TimelineSchemaManager{SchemaManager: schemaMgr}
}

func (s *TimelineSchemaManager) InitSchema() error {
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		date TEXT NOT NULL,
		description TEXT,
		category TEXT,
		tags TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_events_date ON events(date)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_events_category ON events(category)`)

	return nil
}

func (s *TimelineSchemaManager) Initialize() error {
	_, _ = s.SchemaManager.DB.GetMeta("_ensure_meta")

	initialized, err := s.IsInitialized()
	if err != nil {
		return err
	}

	if initialized {
		return nil
	}

	if err := s.InitSchema(); err != nil {
		return err
	}

	return s.SetVersion(SCHEMA_VERSION)
}

func InitDatabase() error {
	return NewTimelineSchemaManager().Initialize()
}
