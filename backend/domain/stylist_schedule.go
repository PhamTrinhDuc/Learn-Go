package domain

type StylistSchedule struct {
	ID        string `json:"id" db:"id"`
	StylistID string `json:"stylist_id" db:"stylist_id"`
	DayOfWeek int    `json:"day_of_week" db:"day_of_week"` // 0 = Sunday, 6 = Saturday
	StartTime string `json:"start_time" db:"start_time"`   // "09:00"
	EndTime   string `json:"end_time" db:"end_time"`       // "18:00"
	IsActive  bool   `json:"is_active" db:"is_active"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type StylistScheduleRepo interface {
	GetByID(id string) (*StylistSchedule, error)
	Create(schedule *StylistSchedule) error
	Update(schedule *StylistSchedule) error
	Delete(id string) error
	ListByStylistID(stylistID string) ([]*StylistSchedule, error)
}
