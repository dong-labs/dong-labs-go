package db

import (
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "1.0.0"
)

type SchemaManager struct {
	*db.SchemaManager
}

func NewSchemaManager() *SchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &SchemaManager{SchemaManager: schemaMgr}
}

func (s *SchemaManager) InitSchema() error {
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		title TEXT NOT NULL,
		note TEXT,
		source TEXT,
		type_val TEXT,
		tags TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_items_type ON items(type_val)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_items_tags ON items(tags)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_items_source ON items(source)`)

	return nil
}

func (s *SchemaManager) Initialize() error {
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
	return NewSchemaManager().Initialize()
}
