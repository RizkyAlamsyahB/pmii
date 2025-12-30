package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// AuthService interface untuk business logic authentication
type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	Logout(ctx context.Context, userID int, token string) error
	ChangePassword(userID int, req requests.ChangePasswordRequest) error
}

type authService struct {
	userRepo        repository.UserRepository
	activityLogRepo repository.ActivityLogRepository
}

// NewAuthService constructor untuk AuthService
func NewAuthService(userRepo repository.UserRepository, activityLogRepo repository.ActivityLogRepository) AuthService {
	return &authService{
		userRepo:        userRepo,
		activityLogRepo: activityLogRepo,
	}
}

// Login melakukan proses login user
func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	// 1. Cari user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// 2. Verify password dengan bcrypt
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	// 3. Cek status user aktif
	if !user.IsActive {
		return nil, "", errors.New("user account is inactive")
	}

	// 4. Generate JWT token (convert role int to string)
	token, err := utils.GenerateJWT(user.ID, strconv.Itoa(user.Role))
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	// 5. Log activity (synchronous)
	s.logActivity(ctx, user.ID, domain.ActionLogin, domain.ModuleAuth, "User berhasil login", nil, map[string]any{
		"user_id":   user.ID,
		"email":     user.Email,
		"full_name": user.FullName,
	})

	return user, token, nil
}

// Logout melakukan proses logout user dengan blacklist token
func (s *authService) Logout(ctx context.Context, userID int, token string) error {
	// Validate token untuk get expiry time
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return errors.New("invalid token")
	}

	// Add token to blacklist dengan expiry time
	utils.AddToBlacklist(token, claims.ExpiresAt.Time)

	// Log activity (synchronous)
	s.logActivity(ctx, userID, domain.ActionLogout, domain.ModuleAuth, "User berhasil logout", nil, nil)

	return nil
}

func (s *authService) ChangePassword(userID int, req requests.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if !user.IsActive {
		return errors.New("user tidak aktif")
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("password lama salah")
	}

	if req.OldPassword == req.NewPassword {
		return errors.New("password baru tidak boleh sama dengan password lama")
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	return s.userRepo.Update(user)
}

// logActivity helper untuk mencatat activity log
func (s *authService) logActivity(ctx context.Context, userID int, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any) {
	ipAddress := utils.GetIPAddress(ctx)
	userAgent := utils.GetUserAgent(ctx)

	var ipPtr, uaPtr *string
	if ipAddress != "" {
		ipPtr = &ipAddress
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	log := &domain.ActivityLog{
		UserID:      userID,
		ActionType:  actionType,
		Module:      module,
		Description: &description,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
	}

	// Ignore error - logging should not affect main operation
	_ = s.activityLogRepo.Create(log)
}
