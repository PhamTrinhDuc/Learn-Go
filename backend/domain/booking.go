package domain

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

type Source string

const (
	SourceZalo   Source = "zalo"
	SourceWeb    Source = "web"
	SourceAgent  Source = "agent"
	SourceManual Source = "manual"
)

type Booking struct {
	ID              string `json:"id" db:"id"`
	UserID          string `json:"user_id" db:"user_id"`
	BranchID        string `json:"branch_id" db:"branch_id"`
	StylistID       string `json:"stylist_id" db:"stylist_id"`
	ServiceID       string `json:"service_id" db:"service_id"`
	ScheduledAt     string `json:"scheduled_at" db:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes" db:"duration_minutes"`
	Status          Status `json:"status" db:"status"`
	Source          Source `json:"source" db:"source"`
	CancelReason    string `json:"cancel_reason,omitempty" db:"cancel_reason"`
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
}

type BookingRepo interface {
	GetByID(id string) (*Booking, error)
	Create(booking *Booking) error
	Update(booking *Booking) error
	Delete(id string) error
	ListByUserID(userId string) ([]*Booking, error)
	ListByBranchID(branchId string) ([]*Booking, error)
	ListByStylistID(stylistId string) ([]*Booking, error)
}
