package domain

import (
	"context"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusNoShow    BookingStatus = "no_show"
)

type BookingSource string

const (
	BookingSourceZalo   BookingSource = "zalo"
	BookingSourceWeb    BookingSource = "web"
	BookingSourceAgent  BookingSource = "agent"
	BookingSourceManual BookingSource = "manual"
)

// Booking - Đơn booking của khách hàng
type Booking struct {
	ID              string        `json:"id" db:"id"`
	UserID          string        `json:"user_id" db:"user_id"`
	BranchID        string        `json:"branch_id" db:"branch_id"`
	StylistID       string        `json:"stylist_id" db:"stylist_id"`
	ServiceID       string        `json:"service_id" db:"service_id"`
	ScheduledAt     time.Time     `json:"scheduled_at" db:"scheduled_at"`
	DurationMinutes int           `json:"duration_minutes" db:"duration_minutes"`
	Status          BookingStatus `json:"status" db:"status"`
	CancelReason    *string       `json:"cancel_reason,omitempty" db:"cancel_reason"`
	Source          BookingSource `json:"source,omitempty" db:"source"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

type BookingFilter struct {
	UserID    *string
	BranchID  *string
	StylistID *string
	Status    []string
	FromDate  *time.Time
	ToDate    *time.Time
	Limit     int
	Offset    int
}

type BookingRepository interface {
	// Tạo booking mới
	CreateBooking(ctx context.Context, booking Booking) (*Booking, error)

	// Lấy chi tiết 1 booking
	GetBookingByID(ctx context.Context, bookingID string) (*Booking, error)

	// Lấy tất cả booking của customer
	GetBookingsByUser(ctx context.Context, userID string, status ...string) ([]Booking, error)

	// Lấy tất cả booking của stylist (dùng để check schedule)
	GetBookingsByStylist(ctx context.Context, stylistID string, fromTime, toTime time.Time) ([]Booking, error)

	// Lấy tất cả booking của chi nhánh trong khoảng thời gian
	GetBookingsByBranch(ctx context.Context, branchID string, fromTime, toTime time.Time) ([]Booking, error)

	// Tìm slot trống cho dịch vụ/stylist/ngày
	// Trả về list time.Time của các slot khả dụng
	GetAvailableSlots(ctx context.Context, branchID, stylistID, serviceID string, targetDate time.Time) ([]time.Time, error)

	// Cập nhật trạng thái booking
	UpdateBookingStatus(ctx context.Context, bookingID string, status BookingStatus, cancelReason *string) error

	// Sửa booking (đổi stylist, dịch vụ, thời gian)
	UpdateBooking(ctx context.Context, booking Booking) (Booking, error)

	// Hủy booking
	CancelBooking(ctx context.Context, bookingID string, cancelReason string) error

	// Lấy booking sắp tới trong N giờ (dùng cho reminder)
	GetUpcomingBookings(ctx context.Context, withinHours int) ([]Booking, error)

	// Lấy booking quá giờ hẹn mà không có check-in (no-show detection)
	GetNoShowBookings(ctx context.Context, minutesOverdue int) ([]Booking, error)

	// Kiểm tra stylist có bận tại khung giờ cụ thể không
	IsStylistBusy(ctx context.Context, stylistID string, startTime, endTime time.Time) (bool, error)

	// Lấy booking theo filter phức tạp (dùng cho admin)
	GetBookingsFiltered(ctx context.Context, filter *BookingFilter) ([]Booking, int64, error)
}
