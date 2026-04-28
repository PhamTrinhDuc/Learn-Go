package repository

import (
	"backend/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookingRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) *bookingRepository {
	return &bookingRepository{pool: pool}
}

// ============================================================
// CREATE
// ============================================================

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := `
		INSERT INTO booking (
			id, user_id, branch_id, stylist_id, service_id,
			scheduled_at, duration_minutes, status, source
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, branch_id, stylist_id, service_id,
				  scheduled_at, duration_minutes, status, cancel_reason, source
	`

	row := r.pool.QueryRow(ctx, query,
		booking.ID, booking.UserID, booking.BranchID, booking.StylistID, booking.ServiceID,
		booking.ScheduledAt, booking.DurationMinutes, booking.Status, booking.Source,
	)

	err := row.Scan(
		&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
		&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return booking, nil
}

// ============================================================
// READ - Single
// ============================================================

func (r *bookingRepository) GetBookingByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, bookingID)

	booking := &domain.Booking{}
	err := row.Scan(
		&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
		&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
		&booking.CreatedAt, &booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return booking, nil
}

// ============================================================
// READ - Multiple
// ============================================================

// GetBookingsByUser - Lấy lịch sử booking của customer
func (r *bookingRepository) GetBookingsByUser(ctx context.Context, userID string, status ...string) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE user_id = $1
	`

	args := []interface{}{userID}

	// Nếu filter theo status
	if len(status) > 0 {
		query += ` AND status = ANY($2)`
		args = append(args, status)
	}

	query += ` ORDER BY scheduled_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings by user: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

// GetBookingsByStylist - Lấy lịch booking của stylist trong khoảng thời gian
func (r *bookingRepository) GetBookingsByStylist(ctx context.Context, stylistID string, fromTime, toTime time.Time) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE stylist_id = $1
		  AND scheduled_at >= $2
		  AND scheduled_at <= $3
		  AND status != 'cancelled'
		ORDER BY scheduled_at ASC
	`

	rows, err := r.pool.Query(ctx, query, stylistID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings by stylist: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

// GetBookingsByBranch - Lấy tất cả booking của chi nhánh trong khoảng thời gian
func (r *bookingRepository) GetBookingsByBranch(ctx context.Context, branchID string, fromTime, toTime time.Time) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE branch_id = $1
		  AND scheduled_at >= $2
		  AND scheduled_at <= $3
		  AND status != 'cancelled'
		ORDER BY scheduled_at ASC
	`

	rows, err := r.pool.Query(ctx, query, branchID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings by branch: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

// ============================================================
// SLOTS & AVAILABILITY
// ============================================================

// GetAvailableSlots - Tìm slot trống cho stylist trong ngày
// Business logic:
// - Stylist phải có schedule trên ngày đó
// - Slot = duration dịch vụ + 5 phút buffer
// - Tránh booking trùng lặp
func (r *bookingRepository) GetAvailableSlots(ctx context.Context, branchID, stylistID, serviceID string, targetDate time.Time) ([]time.Time, error) {
	// 1. Lấy giờ làm việc của stylist hôm đó
	scheduleQuery := `
		SELECT start_time, end_time
		FROM stylist_schedule
		WHERE stylist_id = $1
		  AND day_of_week = $2
		  AND is_active = true
	`

	dayOfWeek := int(targetDate.Weekday())
	row := r.pool.QueryRow(ctx, scheduleQuery, stylistID, dayOfWeek)

	var startTime, endTime string
	err := row.Scan(&startTime, &endTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []time.Time{}, nil // Stylist không làm ca đó
		}
		return nil, fmt.Errorf("failed to get stylist schedule: %w", err)
	}

	// 2. Lấy duration dịch vụ
	serviceQuery := `SELECT duration_minutes FROM service WHERE id = $1`
	serviceRow := r.pool.QueryRow(ctx, serviceQuery, serviceID)

	var durationMinutes int
	err = serviceRow.Scan(&durationMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to get service duration: %w", err)
	}

	// 3. Lấy tất cả booking đã confirm/pending hôm đó
	bookingQuery := `
		SELECT scheduled_at, duration_minutes
		FROM booking
		WHERE stylist_id = $1
		  AND DATE(scheduled_at) = $2
		  AND status IN ('pending', 'confirmed')
		ORDER BY scheduled_at ASC
	`

	bookingRows, err := r.pool.Query(ctx, bookingQuery, stylistID, targetDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer bookingRows.Close()

	type BookedSlot struct {
		Start time.Time
		End   time.Time
	}

	var bookedSlots []BookedSlot
	for bookingRows.Next() {
		var scheduledAt time.Time
		var bookingDuration int
		err := bookingRows.Scan(&scheduledAt, &bookingDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookedSlots = append(bookedSlots, BookedSlot{
			Start: scheduledAt,
			End:   scheduledAt.Add(time.Duration(bookingDuration+5) * time.Minute),
		})
	}

	// ============================================================
	// PHẦN 4: GENERATE AVAILABLE SLOTS
	// ============================================================
	// Mục đích: Tạo danh sách slot khả dụng từ 8am-5pm với interval 30 phút
	// Logic:
	//   - Bắt đầu từ giờ mở cửa, kết thúc tại giờ đóng cửa
	//   - Mỗi slot = khung thời gian để bắt đầu booking
	//   - Slot phải đủ dài cho dịch vụ + 5p buffer
	//   - Slot không được overlap với booking hiện tại
	// ============================================================

	// Step 4a: Parse giờ làm việc từ chuỗi HH:MM
	// Ví dụ: startTime = "08:00", endTime = "20:00"
	// Parse thành time.Time để có thể tính toán
	startTimeOfDay, _ := time.Parse("15:04", startTime)
	endTimeOfDay, _ := time.Parse("15:04", endTime)

	// Step 4b: Ghép giờ làm việc với ngày cụ thể
	// Ví dụ: 2024-04-28 08:00:00, 2024-04-28 20:00:00
	// Để so sánh với booking thực tế (có cả năm-tháng-ngày)
	startDateTime := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		startTimeOfDay.Hour(), startTimeOfDay.Minute(), 0, 0,
		targetDate.Location(),
	)
	endDateTime := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		endTimeOfDay.Hour(), endTimeOfDay.Minute(), 0, 0,
		targetDate.Location(),
	)

	var availableSlots []time.Time

	// Step 4c: Iterate qua từng slot (30 phút một lần)
	// Ví dụ: 08:00, 08:30, 09:00, 09:30, ... 20:00
	for t := startDateTime;
	// Điều kiện: slot kết thúc trước hoặc bằng giờ đóng cửa
	t.Add(time.Duration(durationMinutes)*time.Minute).Before(endDateTime) ||
		t.Add(time.Duration(durationMinutes)*time.Minute).Equal(endDateTime); t = t.Add(30 * time.Minute) {

		// Tính thời gian kết thúc của slot này (tính duration dịch vụ, không có buffer)
		// Ví dụ: nếu dịch vụ cắt tóc 30p, slot 08:00 có slotEnd = 08:30
		slotEnd := t.Add(time.Duration(durationMinutes) * time.Minute)

		// Step 4d: Check xem slot này có bị chiếm không (conflict detection)
		isAvailable := true
		for _, booked := range bookedSlots {
			// Overlap logic (2D interval intersection):
			// Slot [t, slotEnd] và Booking [booked.Start, booked.End] trùng nếu:
			// - Slot bắt đầu trước khi booking kết thúc: t < booked.End
			// - Và slot kết thúc sau khi booking bắt đầu: slotEnd > booked.Start
			//
			// Ví dụ:
			// Slot:    [08:00 -------- 08:30]
			// Booking:         [08:15 -------- 08:45]  ❌ OVERLAP (không được đặt)
			//
			// Slot:    [08:00 -------- 08:30]
			// Booking:                       [08:30 -------- 09:00]  ✅ OK (buffer tách rời)
			//
			// Slot:    [08:30 -------- 09:00]
			// Booking: [08:00 -------- 08:35]  ❌ OVERLAP (booking vượt qua chút)
			if t.Before(booked.End) && slotEnd.After(booked.Start) {
				isAvailable = false
				break // Tìm thấy conflict, không cần check tiếp
			}
		}

		// Step 4e: Nếu slot trống, thêm vào danh sách
		if isAvailable {
			availableSlots = append(availableSlots, t)
		}
	}

	// Step 4f: Return danh sách slot khả dụng
	// Ví dụ output:
	// [
	//   2024-04-28 08:00:00,
	//   2024-04-28 08:30:00,
	//   2024-04-28 09:00:00,
	//   ...
	// ]
	return availableSlots, nil
}

// IsStylistBusy - Kiểm tra stylist có bận tại khung giờ cụ thể
func (r *bookingRepository) IsStylistBusy(ctx context.Context, stylistID string, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM booking
			WHERE stylist_id = $1
			  AND status IN ('pending', 'confirmed')
			  -- Check overlap: start < other_end AND end > other_start
			  AND $2 < (scheduled_at + INTERVAL '1 minute' * duration_minutes)
			  AND $3 > scheduled_at
		)
	`

	var isBusy bool
	err := r.pool.QueryRow(ctx, query, stylistID, startTime, endTime).Scan(&isBusy)
	if err != nil {
		return false, fmt.Errorf("failed to check stylist availability: %w", err)
	}

	return isBusy, nil
}

// ============================================================
// UPDATE
// ============================================================

// UpdateBookingStatus - Cập nhật trạng thái booking
func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, bookingID string, status domain.BookingStatus, cancelReason *string) error {
	query := `
		UPDATE booking
		SET status = $1, cancel_reason = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.pool.Exec(ctx, query, status, cancelReason, bookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// UpdateBooking - Sửa booking (đổi stylist, dịch vụ, thời gian)
func (r *bookingRepository) UpdateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	query := `
		UPDATE booking
		SET stylist_id = $1, service_id = $2, scheduled_at = $3,
		    duration_minutes = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, user_id, branch_id, stylist_id, service_id,
				  scheduled_at, duration_minutes, status, cancel_reason, source,
				  created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		booking.StylistID, booking.ServiceID, booking.ScheduledAt,
		booking.DurationMinutes, booking.ID,
	)

	err := row.Scan(
		&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
		&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
		&booking.CreatedAt, &booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	return booking, nil
}

// CancelBooking - Hủy booking
func (r *bookingRepository) CancelBooking(ctx context.Context, bookingID string, cancelReason string) error {
	return r.UpdateBookingStatus(ctx, bookingID, domain.BookingStatusCancelled, &cancelReason)
}

// ============================================================
// SPECIAL QUERIES
// ============================================================

// GetUpcomingBookings - Lấy booking sắp tới trong N giờ (dùng cho reminder)
func (r *bookingRepository) GetUpcomingBookings(ctx context.Context, withinHours int) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE status IN ('pending', 'confirmed')
		  AND scheduled_at > NOW()
		  AND scheduled_at <= NOW() + INTERVAL '1 hour' * $1
		ORDER BY scheduled_at ASC
	`

	rows, err := r.pool.Query(ctx, query, withinHours)
	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

// GetNoShowBookings - Lấy booking quá giờ hẹn mà trạng thái vẫn pending/confirmed
func (r *bookingRepository) GetNoShowBookings(ctx context.Context, minutesOverdue int) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, branch_id, stylist_id, service_id,
		       scheduled_at, duration_minutes, status, cancel_reason, source,
		       created_at, updated_at
		FROM booking
		WHERE status IN ('pending', 'confirmed')
		  AND scheduled_at + INTERVAL '1 minute' * duration_minutes < NOW()
		  AND scheduled_at + INTERVAL '1 minute' * duration_minutes > NOW() - INTERVAL '1 minute' * $1
		ORDER BY scheduled_at ASC
	`

	rows, err := r.pool.Query(ctx, query, minutesOverdue)
	if err != nil {
		return nil, fmt.Errorf("failed to query no-show bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

// ============================================================
// FILTER
// ============================================================

// GetBookingsFiltered - Lấy booking với filter phức tạp (dùng cho admin dashboard)
func (r *bookingRepository) GetBookingsFiltered(ctx context.Context, filter *domain.BookingFilter) ([]domain.Booking, int64, error) {
	query := `SELECT id, user_id, branch_id, stylist_id, service_id,
	       scheduled_at, duration_minutes, status, cancel_reason, source,
	       created_at, updated_at
	FROM booking
	WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM booking WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(` AND user_id = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND user_id = $%d`, argCount)
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.BranchID != nil {
		query += fmt.Sprintf(` AND branch_id = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND branch_id = $%d`, argCount)
		args = append(args, *filter.BranchID)
		argCount++
	}

	if filter.StylistID != nil {
		query += fmt.Sprintf(` AND stylist_id = $%d`, argCount)
		countQuery += fmt.Sprintf(` AND stylist_id = $%d`, argCount)
		args = append(args, *filter.StylistID)
		argCount++
	}

	if len(filter.Status) > 0 {
		query += fmt.Sprintf(` AND status = ANY($%d)`, argCount)
		countQuery += fmt.Sprintf(` AND status = ANY($%d)`, argCount)
		args = append(args, filter.Status)
		argCount++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(` AND scheduled_at >= $%d`, argCount)
		countQuery += fmt.Sprintf(` AND scheduled_at >= $%d`, argCount)
		args = append(args, *filter.FromDate)
		argCount++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(` AND scheduled_at <= $%d`, argCount)
		countQuery += fmt.Sprintf(` AND scheduled_at <= $%d`, argCount)
		args = append(args, *filter.ToDate)
		argCount++
	}

	query += ` ORDER BY scheduled_at DESC`
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argCount, argCount+1)

	// Get total count
	countArgs := args[:len(args)]
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Get filtered data
	args = append(args, filter.Limit, filter.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query filtered bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		booking := domain.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.BranchID, &booking.StylistID, &booking.ServiceID,
			&booking.ScheduledAt, &booking.DurationMinutes, &booking.Status, &booking.CancelReason, &booking.Source,
			&booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, total, rows.Err()
}
