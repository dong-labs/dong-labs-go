// Package models 定义数据模型
package models

import (
	"time"
)

// Thought 想法模型
type Thought struct {
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	Tags        string    `json:"tags,omitempty"`
	Priority    string    `json:"priority,omitempty"`    // low/normal/high
	Status      string    `json:"status,omitempty"`      // active/completed/archived
	Context     string    `json:"context,omitempty"`     // 上下文
	SourceAgent string    `json:"source_agent,omitempty"` // 来源智能体
	Note        string    `json:"note,omitempty"`        // 备注
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ThoughtList 想法列表响应
type ThoughtList struct {
	Total int       `json:"total"`
	Items []Thought `json:"items"`
}

// StatsResponse 统计响应
type StatsResponse struct {
	Total        int                `json:"total"`
	ThisWeek     int                `json:"this_week"`
	ThisMonth    int                `json:"this_month"`
	ThisYear     int                `json:"this_year"`
	StatusStats  map[string]int     `json:"status_stats,omitempty"`
	PriorityStats map[string]int    `json:"priority_stats,omitempty"`
	TopTags      []TagStat          `json:"top_tags,omitempty"`
}

// TagStat 标签统计
type TagStat struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// ReviewResponse 回顾响应
type ReviewResponse struct {
	Random       []Thought `json:"random,omitempty"`
	Recent       []Thought `json:"recent,omitempty"`
	ThisWeek     []Thought `json:"this_week,omitempty"`
	ThisMonth    []Thought `json:"this_month,omitempty"`
	UnreadCount  int       `json:"unread_count"`
}

// TagsResponse 标签列表响应
type TagsResponse struct {
	Total int               `json:"total"`
	Tags  []TagStat         `json:"tags"`
}
