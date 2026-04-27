package db

import (
	"github.com/dong-labs/think/internal/core/db"
)

const (
	SCHEMA_VERSION = "1.3.0"
)

type MemberSchemaManager struct {
	*db.SchemaManager
}

func NewMemberSchemaManager() *MemberSchemaManager {
	schemaMgr := db.NewSchemaManager(GetDB().Database, SCHEMA_VERSION)
	return &MemberSchemaManager{SchemaManager: schemaMgr}
}

func (s *MemberSchemaManager) InitSchema() error {
	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		wechat TEXT,
		phone TEXT,
		email TEXT,
		account_id TEXT,
		member_type TEXT DEFAULT 'yearly',
		project TEXT DEFAULT 'donglijuan',
		join_date TEXT NOT NULL,
		expire_date TEXT,
		price REAL,
		currency TEXT DEFAULT 'CNY',
		status TEXT DEFAULT 'active',
		source TEXT,
		region TEXT,
		job TEXT,
		tech_level TEXT,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	if err := s.CreateTable(`CREATE TABLE IF NOT EXISTS renewals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		member_id INTEGER NOT NULL,
		old_expire_date TEXT,
		new_expire_date TEXT NOT NULL,
		amount REAL,
		currency TEXT DEFAULT 'CNY',
		renewed_at TEXT DEFAULT CURRENT_TIMESTAMP,
		notes TEXT,
		FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
	)`); err != nil {
		return err
	}

	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_members_status ON members(status)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_members_expire_date ON members(expire_date)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_members_type ON members(member_type)`)
	s.CreateIndex(`CREATE INDEX IF NOT EXISTS idx_renewals_member_id ON renewals(member_id)`)

	return nil
}

func (s *MemberSchemaManager) Initialize() error {
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
	return NewMemberSchemaManager().Initialize()
}
