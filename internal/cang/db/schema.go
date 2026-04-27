package db

import (
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "3.0.0"
)

type CangSchemaManager struct {
	*db.SchemaManager
}

func NewCangSchemaManager() *CangSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &CangSchemaManager{SchemaManager: schemaMgr}
}

func (s *CangSchemaManager) InitSchema() error {
	// accounts 表
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		type TEXT NOT NULL,
		currency TEXT DEFAULT 'CNY',
		balance_cents INTEGER DEFAULT 0,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	// transactions 表
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		amount_cents INTEGER NOT NULL,
		account_id INTEGER,
		category TEXT,
		note TEXT,
		tags TEXT DEFAULT '',
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	// categories 表
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`); err != nil {
		return err
	}

	// budgets 表
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		amount_cents INTEGER NOT NULL,
		period TEXT NOT NULL,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_transactions_account ON transactions(account_id)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_budgets_period ON budgets(period)`)

	return nil
}

func (s *CangSchemaManager) Initialize() error {
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
	return NewCangSchemaManager().Initialize()
}
