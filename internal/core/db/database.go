// Package db 提供统一的数据库管理基类
//
// 所有 dong 家族 CLI 共享 ~/.dong/ 目录
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DONG_DIR 统一数据目录
var DONG_DIR = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".dong")
}()

// Database 数据库管理基类
//
// 子类需要实现：
// - GetName() 返回 CLI 名称
//
// 数据库路径: ~/.dong/<name>/<name>.db
type Database struct {
	name       string
	dbPath     string
	moduleDir  string
	connection *sql.DB
	once       sync.Once
	initErr    error
}

// NewDatabase 创建新的数据库实例
func NewDatabase(name string) *Database {
	moduleDir := filepath.Join(DONG_DIR, name)
	return &Database{
		name:      name,
		dbPath:    filepath.Join(moduleDir, name+".db"),
		moduleDir: moduleDir,
	}
}

// GetName 返回 CLI 名称
func (d *Database) GetName() string {
	return d.name
}

// GetDongDir 获取咚咚家族统一目录
func (d *Database) GetDongDir() string {
	return DONG_DIR
}

// GetModuleDir 获取模块目录 ~/.dong/<name>/
func (d *Database) GetModuleDir() string {
	return d.moduleDir
}

// GetDBPath 获取数据库文件路径 ~/.dong/<name>/<name>.db
func (d *Database) GetDBPath() string {
	return d.dbPath
}

// ensureModuleDir 确保模块目录存在
func (d *Database) ensureModuleDir() error {
	return os.MkdirAll(d.moduleDir, 0755)
}

// GetConnection 获取数据库连接（单例）
func (d *Database) GetConnection() (*sql.DB, error) {
	d.once.Do(func() {
		// 确保目录存在
		if err := d.ensureModuleDir(); err != nil {
			d.initErr = fmt.Errorf("创建目录失败: %w", err)
			return
		}

		// 连接数据库
		conn, err := sql.Open("sqlite3", d.dbPath)
		if err != nil {
			d.initErr = fmt.Errorf("打开数据库失败: %w", err)
			return
		}

		// 设置连接参数
		conn.SetMaxOpenConns(1) // SQLite 不支持多写入并发
		conn.SetMaxIdleConns(1)

		d.connection = conn
	})

	if d.initErr != nil {
		return nil, d.initErr
	}

	return d.connection, nil
}

// CloseConnection 关闭数据库连接
func (d *Database) CloseConnection() error {
	if d.connection != nil {
		err := d.connection.Close()
		d.connection = nil
		d.once = sync.Once{} // 重置 once，允许重新连接
		return err
	}
	return nil
}

// GetCursor 获取数据库游标的上下文管理器
// 返回一个事务，使用完后需要 Commit 或 Rollback
func (d *Database) GetCursor() (*sql.Tx, error) {
	conn, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return conn.Begin()
}

// GetMeta 获取元数据
func (d *Database) GetMeta(key string) (string, error) {
	d.ensureMetaTable()

	conn, err := d.GetConnection()
	if err != nil {
		return "", err
	}

	var value string
	err = conn.QueryRow("SELECT value FROM __meta WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SetMeta 设置元数据
func (d *Database) SetMeta(key, value string) error {
	d.ensureMetaTable()

	conn, err := d.GetConnection()
	if err != nil {
		return err
	}

	_, err = conn.Exec(`
		INSERT OR REPLACE INTO __meta (key, value)
		VALUES (?, ?)
	`, key, value)
	return err
}

// ensureMetaTable 确保元数据表存在
func (d *Database) ensureMetaTable() error {
	conn, err := d.GetConnection()
	if err != nil {
		return err
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS __meta (
			key TEXT PRIMARY KEY,
			value TEXT
		)
	`)
	return err
}

// Exec 执行 SQL 语句（自动提交）
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	conn, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return conn.Exec(query, args...)
}

// Query 执行查询语句
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	conn, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return conn.Query(query, args...)
}

// QueryRow 执行查询单行语句
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	conn, err := d.GetConnection()
	if err != nil {
		// 返回一个已经出错的 Row
		return &sql.Row{}
	}
	return conn.QueryRow(query, args...)
}

// WithTransaction 在事务中执行函数
// 如果函数返回错误，事务会回滚；否则提交
func (d *Database) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.GetCursor()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// TableExists 检查表是否存在
func (d *Database) TableExists(tableName string) (bool, error) {
	d.ensureMetaTable()

	conn, err := d.GetConnection()
	if err != nil {
		return false, err
	}

	var count int
	err = conn.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master
		WHERE type='table' AND name=?
	`, tableName).Scan(&count)

	return count > 0, err
}

// DropTable 删除表
func (d *Database) DropTable(tableName string) error {
	conn, err := d.GetConnection()
	if err != nil {
		return err
	}

	_, err = conn.Exec("DROP TABLE IF EXISTS " + tableName)
	return err
}

// GetTableNames 获取所有表名
func (d *Database) GetTableNames() ([]string, error) {
	conn, err := d.GetConnection()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(`
		SELECT name FROM sqlite_master
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, rows.Err()
}

// CountRows 统计表中的行数
func (d *Database) CountRows(tableName string) (int, error) {
	conn, err := d.GetConnection()
	if err != nil {
		return 0, err
	}

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	return count, err
}

// Vacuum 优化数据库（VACUUM 命令）
func (d *Database) Vacuum() error {
	conn, err := d.GetConnection()
	if err != nil {
		return err
	}

	_, err = conn.Exec("VACUUM")
	return err
}

// BackupTo 备份数据库到指定路径
func (d *Database) BackupTo(destPath string) error {
	// 读取源文件
	data, err := os.ReadFile(d.dbPath)
	if err != nil {
		return fmt.Errorf("读取数据库失败: %w", err)
	}

	// 确保目标目录存在
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 写入备份文件
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("写入备份文件失败: %w", err)
	}

	return nil
}

// GetSize 获取数据库文件大小（字节）
func (d *Database) GetSize() (int64, error) {
	info, err := os.Stat(d.dbPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Exists 检查数据库文件是否存在
func (d *Database) Exists() bool {
	_, err := os.Stat(d.dbPath)
	return err == nil
}

// Delete 删除数据库文件
func (d *Database) Delete() error {
	d.CloseConnection()
	return os.Remove(d.dbPath)
}
