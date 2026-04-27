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
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT,
		status TEXT DEFAULT 'pending',
		priority TEXT DEFAULT 'medium',
		due_date TEXT,
		tags TEXT,
		note TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
		completed_at TEXT
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date)`)

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
