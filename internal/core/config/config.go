// Package config 提供统一的配置管理
//
// 所有 dong 家族 CLI 共享 ~/.dong/config.json
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// DONG_DIR 统一数据目录
var DONG_DIR = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".dong")
}()

// CONFIG_FILE 统一配置文件
var CONFIG_FILE = filepath.Join(DONG_DIR, "config.json")

// allConfig 全局配置缓存
var (
	allConfig   map[string]interface{}
	allConfigMu sync.RWMutex
)

// Config 配置管理基类
//
// 子类需要实现：
// - GetName() 返回 CLI 名称
// - GetDefaults() 返回默认配置字典
//
// 配置文件格式:
// {
//     "log": {"default_group": "work", ...},
//     "cang": {"default_account": 1, ...},
//     ...
// }
type Config struct {
	name     string
	defaults map[string]interface{}
}

// NewConfig 创建新的配置实例
func NewConfig(name string, defaults map[string]interface{}) *Config {
	if defaults == nil {
		defaults = make(map[string]interface{})
	}
	return &Config{
		name:     name,
		defaults: defaults,
	}
}

// GetName 返回 CLI 名称
func (c *Config) GetName() string {
	return c.name
}

// GetDefaults 返回默认配置
func (c *Config) GetDefaults() map[string]interface{} {
	return c.defaults
}

// GetConfigFile 获取统一配置文件路径
func (c *Config) GetConfigFile() string {
	return CONFIG_FILE
}

// loadAll 加载整个配置文件
func (c *Config) loadAll() (map[string]interface{}, error) {
	allConfigMu.RLock()
	if allConfig != nil {
		allConfigMu.RUnlock()
		return allConfig, nil
	}
	allConfigMu.RUnlock()

	allConfigMu.Lock()
	defer allConfigMu.Unlock()

	// 双重检查
	if allConfig != nil {
		return allConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// 文件不存在，返回空配置
			allConfig = make(map[string]interface{})
			return allConfig, nil
		}
		return nil, err
	}

	// 解析 JSON
	if err := json.Unmarshal(data, &allConfig); err != nil {
		// JSON 解析失败，返回空配置
		allConfig = make(map[string]interface{})
		return allConfig, nil
	}

	if allConfig == nil {
		allConfig = make(map[string]interface{})
	}

	return allConfig, nil
}

// saveAll 保存整个配置文件
func (c *Config) saveAll(config map[string]interface{}) error {
	// 确保目录存在
	if err := os.MkdirAll(DONG_DIR, 0755); err != nil {
		return err
	}

	// 格式化输出
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	if err := os.WriteFile(CONFIG_FILE, data, 0644); err != nil {
		return err
	}

	// 更新缓存
	allConfigMu.Lock()
	allConfig = config
	allConfigMu.Unlock()

	return nil
}

// Load 加载当前 CLI 的配置
func (c *Config) Load() (map[string]interface{}, error) {
	allCfg, err := c.loadAll()
	if err != nil {
		return nil, err
	}

	name := c.GetName()

	// 如果没有该 CLI 的配置，初始化为默认值
	if _, exists := allCfg[name]; !exists {
		allCfg[name] = c.defaults
		if err := c.saveAll(allCfg); err != nil {
			return nil, err
		}
	}

	// 合并默认配置
	result := make(map[string]interface{})
	for k, v := range c.defaults {
		result[k] = v
	}
	if cliConfig, ok := allCfg[name].(map[string]interface{}); ok {
		for k, v := range cliConfig {
			result[k] = v
		}
	}

	return result, nil
}

// Save 保存当前 CLI 的配置
func (c *Config) Save(config map[string]interface{}) error {
	allCfg, err := c.loadAll()
	if err != nil {
		return err
	}

	allCfg[c.GetName()] = config
	return c.saveAll(allCfg)
}

// Get 获取配置项
func (c *Config) Get(key string) (interface{}, error) {
	config, err := c.Load()
	if err != nil {
		return nil, err
	}

	if val, exists := config[key]; exists {
		return val, nil
	}

	return nil, nil
}

// GetString 获取字符串配置项
func (c *Config) GetString(key string, defaultValue string) string {
	val, err := c.Get(key)
	if err != nil || val == nil {
		return defaultValue
	}

	if str, ok := val.(string); ok {
		return str
	}

	return defaultValue
}

// GetInt 获取整数配置项
func (c *Config) GetInt(key string, defaultValue int) int {
	val, err := c.Get(key)
	if err != nil || val == nil {
		return defaultValue
	}

	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	}

	return defaultValue
}

// GetBool 获取布尔配置项
func (c *Config) GetBool(key string, defaultValue bool) bool {
	val, err := c.Get(key)
	if err != nil || val == nil {
		return defaultValue
	}

	if b, ok := val.(bool); ok {
		return b
	}

	return defaultValue
}

// Set 设置配置项
func (c *Config) Set(key string, value interface{}) error {
	config, err := c.Load()
	if err != nil {
		return err
	}

	config[key] = value
	return c.Save(config)
}

// SetMulti 批量设置配置项
func (c *Config) SetMulti(items map[string]interface{}) error {
	config, err := c.Load()
	if err != nil {
		return err
	}

	for k, v := range items {
		config[k] = v
	}

	return c.Save(config)
}

// Delete 删除配置项
func (c *Config) Delete(key string) error {
	config, err := c.Load()
	if err != nil {
		return err
	}

	delete(config, key)
	return c.Save(config)
}

// Reset 重置当前 CLI 的配置到默认值
func (c *Config) Reset() error {
	allCfg, err := c.loadAll()
	if err != nil {
		return err
	}

	delete(allCfg, c.GetName())

	// 清除缓存，重新加载
	allConfigMu.Lock()
	allConfig = nil
	allConfigMu.Unlock()

	return c.saveAll(allCfg)
}

// ClearCache 清除全局配置缓存
func ClearCache() {
	allConfigMu.Lock()
	allConfig = nil
	allConfigMu.Unlock()
}

// LoadConfig 加载指定名称的配置
func LoadConfig(name string) (map[string]interface{}, error) {
	ClearCache()

	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var allCfg map[string]interface{}
	if err := json.Unmarshal(data, &allCfg); err != nil {
		return nil, err
	}

	if allCfg == nil {
		allCfg = make(map[string]interface{})
	}

	if cliConfig, ok := allCfg[name].(map[string]interface{}); ok {
		return cliConfig, nil
	}

	return make(map[string]interface{}), nil
}

// SaveConfig 保存指定名称的配置
func SaveConfig(name string, config map[string]interface{}) error {
	ClearCache()

	// 加载现有配置
	data, err := os.ReadFile(CONFIG_FILE)
	var allCfg map[string]interface{}
	if err == nil {
		json.Unmarshal(data, &allCfg)
	}
	if allCfg == nil {
		allCfg = make(map[string]interface{})
	}

	// 更新配置
	allCfg[name] = config

	// 保存
	if err := os.MkdirAll(DONG_DIR, 0755); err != nil {
		return err
	}

	output, err := json.MarshalIndent(allCfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(CONFIG_FILE, output, 0644)
}

// GetDongDir 获取咚咚目录
func GetDongDir() string {
	return DONG_DIR
}

// GetConfigFilePath 获取配置文件路径
func GetConfigFilePath() string {
	return CONFIG_FILE
}
