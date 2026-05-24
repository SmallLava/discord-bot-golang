package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the SQLite database and performs migrations
func InitDB(dbPath string) (*gorm.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// AutoMigrate models
	err = db.AutoMigrate(&LeaveRecord{}, &MeetingSchedule{}, &Config{})

	if err != nil {
		return nil, err
	}

	// Initialize default schedules if empty
	initDefaultSchedules(db)

	DB = db
	return db, nil
}

func initDefaultSchedules(db *gorm.DB) {
	var count int64
	db.Model(&MeetingSchedule{}).Count(&count)
	if count == 0 {
		schedules := []MeetingSchedule{
			{DayOfWeek: time.Monday, StartTime: "21:00", EndTime: "21:30"},
			{DayOfWeek: time.Friday, StartTime: "19:00", EndTime: "21:30"},
		}
		for _, s := range schedules {
			db.Create(&s)
		}
		log.Println("Default meeting schedules initialized.")
	}
}
