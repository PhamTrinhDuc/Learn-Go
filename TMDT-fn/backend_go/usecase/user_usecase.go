package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"tmdt-backend/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
	redis    *redis.Client
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo domain.UserRepository, redis *redis.Client) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.UserRegisterRequest) (*domain.User, error) {
	// 1. Check if email or phone already exists
	existingEmail, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email này đã được sử dụng")
	}

	existingPhone, err := u.userRepo.GetByPhone(ctx, req.NumPhone)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, errors.New("số điện thoại đã được đăng ký")
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create user domain model
	newUser := &domain.User{
		Password: string(hashedPassword),
		FullName: req.FullName,
		Email:    req.Email,
		NumPhone: req.NumPhone,
		Role:     "Customer",
		IsLock:   false,
		Gender:   "none",
	}

	// 4. Save to database
	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// Hide password before returning
	newUser.Password = ""
	return newUser, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.UserLoginRequest) (*domain.User, error) {
	// 1. Retrieve user by email or phone
	user, err := u.userRepo.GetByEmailOrPhone(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("tài khoản không tồn tại")
	}

	// 2. Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("mật khẩu không chính xác")
	}

	// 3. Check if user is locked
	if user.IsLock {
		return nil, errors.New("tài khoản này đã bị khóa bởi Admin")
	}

	// Hide password before returning
	user.Password = ""
	return user, nil
}

func (u *userUsecase) GoogleLogin(ctx context.Context, idToken string) (*domain.User, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call google oauth api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("xác thực Google thất bại")
	}

	var googleUserInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Sub   string `json:"sub"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode google oauth response: %w", err)
	}

	if googleUserInfo.Email == "" {
		return nil, errors.New("không thể lấy thông tin email từ Google")
	}

	// Find user by email
	user, err := u.userRepo.GetByEmail(ctx, googleUserInfo.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Create a new user with Google sub as password hash
		hashedSub, err := bcrypt.GenerateFromPassword([]byte(googleUserInfo.Sub), 10)
		if err != nil {
			return nil, fmt.Errorf("failed to hash google sub id: %w", err)
		}

		user = &domain.User{
			Password: string(hashedSub),
			FullName: googleUserInfo.Name,
			Email:    googleUserInfo.Email,
			Role:     "Customer",
			IsLock:   false,
			Gender:   "none",
		}

		err = u.userRepo.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	if user.IsLock {
		return nil, errors.New("tài khoản Google này đã bị khóa")
	}

	user.Password = ""
	return user, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("không tìm thấy người dùng")
	}
	user.Password = ""
	return user, nil
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := u.userRepo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	// Hide passwords
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

func (u *userUsecase) GetRoles(ctx context.Context) ([]string, int, error) {
	return u.userRepo.GetDistinctRoles(ctx)
}

func (u *userUsecase) UpdateProfile(ctx context.Context, id int, req *domain.UserUpdateRequest) (*domain.User, error) {
	// 1. Check if email or phone is already taken by another user
	existingEmail, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil && existingEmail.ID != id {
		return nil, errors.New("email đã được sử dụng bởi tài khoản khác")
	}

	existingPhone, err := u.userRepo.GetByPhone(ctx, req.NumPhone)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil && existingPhone.ID != id {
		return nil, errors.New("số điện thoại đã được sử dụng bởi tài khoản khác")
	}

	// 2. Fetch existing user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("không tìm thấy người dùng")
	}

	// 3. Update fields
	user.FullName = req.FullName
	user.Email = req.Email
	user.NumPhone = req.NumPhone
	user.Gender = req.Gender

	if req.DOB != nil && *req.DOB != "" {
		// Try parsing different common layouts
		if t, err := time.Parse("2006-01-02", *req.DOB); err == nil {
			user.DOB = &t
		} else if t, err := time.Parse("2006-01-02T15:04:05.000Z", *req.DOB); err == nil {
			user.DOB = &t
		} else if t, err := time.Parse(time.RFC3339, *req.DOB); err == nil {
			user.DOB = &t
		} else {
			user.DOB = nil
		}
	} else {
		user.DOB = nil
	}

	// 4. Save to db
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *userUsecase) UpdateRole(ctx context.Context, id int, role string) (*domain.User, error) {
	err := u.userRepo.UpdateRole(ctx, id, role)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (u *userUsecase) UpdateStatus(ctx context.Context, id int, isLock bool) (*domain.User, error) {
	err := u.userRepo.UpdateStatus(ctx, id, isLock)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (u *userUsecase) CheckPasswordValidity(ctx context.Context, email string, newPassword string) (bool, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("người dùng không tồn tại")
	}

	// Compare current password with new password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		// Passwords are same, so it's invalid
		return false, nil
	}

	return true, nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) error {
	// 1. Verify resetToken in Redis
	redisKey := fmt.Sprintf("resetToken:%s", req.Email)
	storedToken, err := u.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return errors.New("yêu cầu không hợp lệ hoặc phiên làm việc đã hết hạn. Vui lòng xác thực lại OTP")
	} else if err != nil {
		return fmt.Errorf("failed to fetch token from redis: %w", err)
	}

	if storedToken != req.ResetToken {
		return errors.New("yêu cầu không hợp lệ hoặc phiên làm việc đã hết hạn. Vui lòng xác thực lại OTP")
	}

	// 2. Hash new password
	saltCost := 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), saltCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Update password in database
	err = u.userRepo.UpdatePassword(ctx, req.Email, string(hashedPassword))
	if err != nil {
		return err
	}

	// 4. Delete token in Redis
	_ = u.redis.Del(ctx, redisKey).Err()

	return nil
}

func (u *userUsecase) CheckEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("email này chưa được đăng ký trong hệ thống")
	}

	if user.IsLock {
		return nil, errors.New("tài khoản này hiện đang bị khóa. Vui lòng liên hệ quản trị viên")
	}

	user.Password = ""
	return user, nil
}

func (u *userUsecase) AdminResetPassword(ctx context.Context, id int, newPassword string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("không tìm thấy người dùng")
	}

	err = u.userRepo.UpdatePassword(ctx, user.Email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
