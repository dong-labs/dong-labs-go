package models

import "time"

type Member struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Wechat      string    `json:"wechat,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	AccountID   string    `json:"account_id,omitempty"`
	MemberType  string    `json:"member_type"`
	Project     string    `json:"project"`
	JoinDate    string    `json:"join_date"`
	ExpireDate  string    `json:"expire_date"`
	Price       float64   `json:"price,omitempty"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Source      string    `json:"source,omitempty"`
	Region      string    `json:"region,omitempty"`
	Job         string    `json:"job,omitempty"`
	TechLevel   string    `json:"tech_level,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Renewal struct {
	ID             int       `json:"id"`
	MemberID       int       `json:"member_id"`
	OldExpireDate  string    `json:"old_expire_date"`
	NewExpireDate  string    `json:"new_expire_date"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	RenewedAt      time.Time `json:"renewed_at"`
	Notes          string    `json:"notes,omitempty"`
}

type MemberList struct {
	Total int       `json:"total"`
	Items []Member  `json:"items"`
}

type StatsResponse struct {
	Total        int           `json:"total"`
	Active       int           `json:"active"`
	Expired      int           `json:"expired"`
	ThisMonth    int           `json:"this_month"`
	ThisYear     int           `json:"this_year"`
	ByType       []TypeStat    `json:"by_type,omitempty"`
	ByProject    []ProjectStat `json:"by_project,omitempty"`
	ByRegion     []RegionStat  `json:"by_region,omitempty"`
}

type TypeStat struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type ProjectStat struct {
	Project string `json:"project"`
	Count   int    `json:"count"`
}

type RegionStat struct {
	Region string `json:"region"`
	Count  int    `json:"count"`
}
