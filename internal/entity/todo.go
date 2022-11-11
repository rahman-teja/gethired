package entity

import "time"

type ToDo struct {
	ID              int64      `json:"id"`
	ActivityGroupID int64      `json:"activity_group_id"`
	Title           string     `json:"title"`
	IsActive        bool       `json:"is_active"`
	Priority        string     `json:"priority"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
