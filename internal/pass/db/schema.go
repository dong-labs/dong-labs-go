package db

import (
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "1.0.0"
)

type PassSchemaManager struct {
	*db.SchemaManager
}

func NewPassSchemaManager() *PassSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &PassSchemaManager{SchemaManager: schemaMgr}
}

func (s *PassSchemaManager) InitSchema() error {
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		username TEXT,
		password TEXT NOT NULL,
		url TEXT,
		category TEXT,
		tags TEXT,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_passwords_category ON passwords(category)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_passwords_tags ON passwords(tags)`)

	return nil
}

func (s *PassSchemaManager) Initialize() error {
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
	return NewPassSchemaManager().Initialize()
}
