package database

import (
	"gorm.io/gorm"
)

// LeaveRepository handles database operations for LeaveRecords
type LeaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

func (r *LeaveRepository) Create(record *LeaveRecord) error {
	return r.db.Create(record).Error
}

func (r *LeaveRepository) GetAll() ([]LeaveRecord, error) {
	var records []LeaveRecord
	err := r.db.Find(&records).Error
	return records, err
}

func (r *LeaveRepository) GetByUserID(userID string) ([]LeaveRecord, error) {
	var records []LeaveRecord
	err := r.db.Where("user_id = ?", userID).Find(&records).Error
	return records, err
}

func (r *LeaveRepository) Delete(id uint) error {
	return r.db.Delete(&LeaveRecord{}, id).Error
}

// ScheduleRepository handles database operations for MeetingSchedules
type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetAll() ([]MeetingSchedule, error) {
	var schedules []MeetingSchedule
	err := r.db.Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) Create(s *MeetingSchedule) error {
	return r.db.Create(s).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&MeetingSchedule{}, id).Error
}

// ConfigRepository handles bot configurations
type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Set(key, value string) error {
	var config Config
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return r.db.Create(&Config{Key: key, Value: value}).Error
	}
	config.Value = value
	return r.db.Save(&config).Error
}

func (r *ConfigRepository) Get(key string) (string, error) {
	var config Config
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}
