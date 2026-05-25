package service

import (
	"time"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
)

type LeaveService struct {
	leaveRepo    *database.LeaveRepository
	scheduleRepo *database.ScheduleRepository
	configRepo   *database.ConfigRepository
}

func NewLeaveService(lRepo *database.LeaveRepository, sRepo *database.ScheduleRepository, cRepo *database.ConfigRepository) *LeaveService {
	return &LeaveService{
		leaveRepo:    lRepo,
		scheduleRepo: sRepo,
		configRepo:   cRepo,
	}
}

// GetMembersOnLeave returns user IDs of members on leave for the given date
func (s *LeaveService) GetMembersOnLeave(date time.Time) ([]string, error) {
	allLeaves, err := s.leaveRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var userIDs []string
	// Use UTC for consistent date comparison
	checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	for _, leave := range allLeaves {
		// Normalize leave range to start and end of their respective days in UTC
		start := time.Date(leave.StartDate.Year(), leave.StartDate.Month(), leave.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(leave.EndDate.Year(), leave.EndDate.Month(), leave.EndDate.Day(), 23, 59, 59, 0, time.UTC)

		if (checkDate.After(start) || checkDate.Equal(start)) && (checkDate.Before(end) || checkDate.Equal(end)) {
			userIDs = append(userIDs, leave.UserID)
		}
	}

	return userIDs, nil
}

func (s *LeaveService) SetAnnouncementChannel(channelID string) error {
	return s.configRepo.Set("announcement_channel", channelID)
}

func (s *LeaveService) GetAnnouncementChannel() (string, error) {
	return s.configRepo.Get("announcement_channel")
}

func (s *LeaveService) GetAllSchedules() ([]database.MeetingSchedule, error) {
	return s.scheduleRepo.GetAll()
}

// GetMeetingsInRange returns all actual meeting times between start and end date
func (s *LeaveService) GetMeetingsInRange(start, end time.Time) ([]time.Time, error) {
	schedules, err := s.scheduleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var meetings []time.Time
	// Standardize to start of day for comparison
	current := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	limit := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

	for !current.After(limit) {
		for _, sched := range schedules {
			if current.Weekday() == sched.DayOfWeek {
				// Parse HH:MM
				t, _ := time.Parse("15:04", sched.StartTime)
				meetingTime := time.Date(current.Year(), current.Month(), current.Day(), t.Hour(), t.Minute(), 0, 0, current.Location())
				
				// Ensure it's within the specific time range if start/end have hours
				if meetingTime.After(start) && meetingTime.Before(end) {
					meetings = append(meetings, meetingTime)
				} else if meetingTime.Equal(start) || meetingTime.Equal(end) {
					meetings = append(meetings, meetingTime)
				}
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	return meetings, nil
}

// AddLeaveRecord creates a new leave record
func (s *LeaveService) AddLeaveRecord(userID, userName string, leaveType database.LeaveType, start, end time.Time, reason string) error {
	record := &database.LeaveRecord{
		UserID:    userID,
		UserName:  userName,
		Type:      leaveType,
		StartDate: start,
		EndDate:   end,
		Reason:    reason,
	}
	return s.leaveRepo.Create(record)
}

// GetUserLeaves returns all leave records for a specific user
func (s *LeaveService) GetUserLeaves(userID string) ([]database.LeaveRecord, error) {
	return s.leaveRepo.GetByUserID(userID)
}

// GetNextNMeetings returns the next N meeting dates starting from now
func (s *LeaveService) GetNextNMeetings(n int) ([]time.Time, error) {
	schedules, err := s.scheduleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var meetings []time.Time
	now := time.Now()
	current := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	// Search forward up to 60 days to find N meetings
	for i := 0; i < 60 && len(meetings) < n; i++ {
		for _, sched := range schedules {
			if current.Weekday() == sched.DayOfWeek {
				t, _ := time.Parse("15:04", sched.StartTime)
				meetingTime := time.Date(current.Year(), current.Month(), current.Day(), t.Hour(), t.Minute(), 0, 0, current.Location())
				
				if meetingTime.After(now) {
					meetings = append(meetings, meetingTime)
					if len(meetings) == n {
						break
					}
				}
			}
		}
		current = current.AddDate(0, 0, 1)
		// Reset time to start of day for subsequent days to catch all meetings
		current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
	}

	return meetings, nil
}

// ParseDate is a helper to parse YYYY-MM-DD
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
