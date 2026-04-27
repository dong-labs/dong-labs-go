// Package models log 数据模型
package models

import "time"

// Log 日志模型
type Log struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	LogGroup  string    `json:"log_group"`
	Date      string    `json:"date"`
	Tags      string    `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LogList 日志列表响应
type LogList struct {
	Total int   `json:"total"`
	Items []Log `json:"items"`
}

// StatsResponse 统计响应
type StatsResponse struct {
	Total     int         `json:"total"`
	ThisWeek  int         `json:"this_week"`
	ThisMonth int         `json:"this_month"`
	ThisYear  int         `json:"this_year"`
	TopTags   []TagStat   `json:"top_tags,omitempty"`
	Groups    []GroupStat `json:"groups,omitempty"`
}

// GroupStat 组统计
type GroupStat struct {
	Group  string `json:"group"`
	Count  int    `json:"count"`
}

// TagStat 标签统计
type TagStat struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// GroupsResponse 组列表响应
type GroupsResponse struct {
	Total   int         `json:"total"`
	Groups  []GroupStat `json:"groups"`
	Default string      `json:"default_group"`
}
