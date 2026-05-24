package database

import (
	"time"

	"gorm.io/gorm"
)

// LeaveType defines the category of leave
type LeaveType string

const (
	LeaveTypeSingle   LeaveType = "single"   // 單次請假
	LeaveTypeInterval LeaveType = "interval" // 時段請假
)

// LeaveRecord represents a leave request from a user
type LeaveRecord struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index;not null" json:"user_id"` // Discord User ID
	UserName  string         `json:"user_name"`
	Type      LeaveType      `gorm:"not null" json:"type"`
	StartDate time.Time      `gorm:"not null" json:"start_date"`
	EndDate   time.Time      `gorm:"not null" json:"end_date"`
	Reason    string         `json:"reason"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MeetingSchedule represents the fixed meeting times
type MeetingSchedule struct {
	ID        uint      `gorm:"primaryKey"`
	DayOfWeek time.Weekday `gorm:"not null"` // 0 (Sunday) to 6 (Saturday)
	StartTime string    `gorm:"not null"` // e.g., "21:00"
	EndTime   string    `gorm:"not null"` // e.g., "21:30"
}

// Config stores bot-wide configurations like announcement channel
type Config struct {
	Key   string `gorm:"primaryKey"`
	Value string
}
