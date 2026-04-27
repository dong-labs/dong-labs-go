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
		name TEXT NOT NULL,
		category TEXT,
		expire_date TEXT NOT NULL,
		reminder_days INTEGER DEFAULT 7,
		tags TEXT,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_items_expire_date ON items(expire_date)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_items_category ON items(category)`)

	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS renew_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id INTEGER NOT NULL,
		old_expire_date TEXT NOT NULL,
		new_expire_date TEXT NOT NULL,
		renewed_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}

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
