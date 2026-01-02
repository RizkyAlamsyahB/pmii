package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

// MemberService interface untuk business logic member
type MemberService interface {
	Create(ctx context.Context, req requests.CreateMemberRequest, photoFile *multipart.FileHeader) (*responses.MemberResponse, error)
	GetAll(ctx context.Context, page, limit int, search string) ([]responses.MemberResponse, int, int, int64, error)
	GetByID(ctx context.Context, id int) (*responses.MemberResponse, error)
	Update(ctx context.Context, id int, req requests.UpdateMemberRequest, photoFile *multipart.FileHeader) (*responses.MemberResponse, error)
	Delete(ctx context.Context, id int) error
}

type memberService struct {
	memberRepo        repository.MemberRepository
	cloudinaryService CloudinaryService
	activityLogRepo   repository.ActivityLogRepository
}

// NewMemberService constructor untuk MemberService
func NewMemberService(memberRepo repository.MemberRepository, cloudinaryService CloudinaryService, activityLogRepo repository.ActivityLogRepository) MemberService {
	return &memberService{
		memberRepo:        memberRepo,
		cloudinaryService: cloudinaryService,
		activityLogRepo:   activityLogRepo,
	}
}

// Create membuat member baru dengan upload foto ke Cloudinary
func (s *memberService) Create(ctx context.Context, req requests.CreateMemberRequest, photoFile *multipart.FileHeader) (*responses.MemberResponse, error) {
	// Upload photo ke Cloudinary (jika ada)
	var photoFilename *string
	if photoFile != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "members", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		photoFilename = &filename
	}

	// Prepare domain entity
	member := &domain.Member{
		FullName:    req.FullName,
		Position:    req.Position,
		Department:  domain.MemberDepartment(req.Department),
		PhotoURI:    photoFilename,
		SocialLinks: req.SocialLinks,
		IsActive:    true,
	}

	// Save ke database
	if err := s.memberRepo.Create(member); err != nil {
		// Rollback: hapus foto dari Cloudinary jika save gagal
		if photoFilename != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "members", *photoFilename)
		}
		return nil, errors.New("gagal menyimpan member")
	}

	// Convert to response DTO
	resp := s.toResponseDTO(member)

	// Log activity - Create Member
	s.logActivity(ctx, domain.ActionCreate, domain.ModuleMembers, "Membuat member baru: "+member.FullName, nil, map[string]any{
		"id":           member.ID,
		"full_name":    member.FullName,
		"position":     member.Position,
		"department":   string(member.Department),
		"photo_uri":    member.PhotoURI,
		"social_links": member.SocialLinks,
		"is_active":    member.IsActive,
		"created_at":   member.CreatedAt,
	}, &member.ID)

	return resp, nil
}

// GetAll mengambil semua member dengan pagination dan search
func (s *memberService) GetAll(ctx context.Context, page, limit int, search string) ([]responses.MemberResponse, int, int, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	members, total, err := s.memberRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, 0, 0, 0, errors.New("gagal mengambil data member")
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	// Auto-clamp: jika page melebihi lastPage dan ada data, clamp ke lastPage
	if page > lastPage && lastPage > 0 {
		page = lastPage
		// Re-fetch dengan page yang sudah di-clamp
		members, _, err = s.memberRepo.FindAll(page, limit, search)
		if err != nil {
			return nil, 0, 0, 0, errors.New("gagal mengambil data member")
		}
	}

	// Convert to response DTOs
	result := make([]responses.MemberResponse, len(members))
	for i, m := range members {
		result[i] = *s.toResponseDTO(&m)
	}

	return result, page, lastPage, total, nil
}

// GetByID mengambil member berdasarkan ID
func (s *memberService) GetByID(ctx context.Context, id int) (*responses.MemberResponse, error) {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("member tidak ditemukan")
	}

	return s.toResponseDTO(member), nil
}

// Update mengupdate member dengan optional upload foto baru
func (s *memberService) Update(ctx context.Context, id int, req requests.UpdateMemberRequest, photoFile *multipart.FileHeader) (*responses.MemberResponse, error) {
	// Ambil member existing
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("member tidak ditemukan")
	}

	// Simpan foto lama untuk rollback
	oldPhotoURI := member.PhotoURI

	// Store old values for audit log
	oldValues := map[string]any{
		"id":           member.ID,
		"full_name":    member.FullName,
		"position":     member.Position,
		"department":   string(member.Department),
		"photo_uri":    member.PhotoURI,
		"social_links": member.SocialLinks,
		"is_active":    member.IsActive,
	}

	// Upload foto baru ke Cloudinary (jika ada)
	var newPhotoFilename *string
	if photoFile != nil {
		filename, err := s.cloudinaryService.UploadImage(ctx, "members", photoFile)
		if err != nil {
			return nil, errors.New("gagal mengupload foto")
		}
		newPhotoFilename = &filename
		member.PhotoURI = &filename
	}

	// Update fields yang dikirim
	if req.FullName != "" {
		member.FullName = req.FullName
	}
	if req.Position != "" {
		member.Position = req.Position
	}
	if req.Department != "" {
		member.Department = domain.MemberDepartment(req.Department)
	}
	if len(req.SocialLinks) > 0 {
		member.SocialLinks = req.SocialLinks
	}
	if req.IsActive != nil {
		member.IsActive = *req.IsActive
	}

	// Save ke database
	if err := s.memberRepo.Update(member); err != nil {
		// Rollback: hapus foto baru jika update gagal
		if newPhotoFilename != nil {
			_ = s.cloudinaryService.DeleteImage(ctx, "members", *newPhotoFilename)
		}
		return nil, errors.New("gagal mengupdate member")
	}

	// Hapus foto lama SETELAH database update berhasil
	if newPhotoFilename != nil && oldPhotoURI != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "members", *oldPhotoURI)
	}

	// Log activity - Update Member
	s.logActivity(ctx, domain.ActionUpdate, domain.ModuleMembers, "Mengupdate member: "+member.FullName, oldValues, map[string]any{
		"id":           member.ID,
		"full_name":    member.FullName,
		"position":     member.Position,
		"department":   string(member.Department),
		"photo_uri":    member.PhotoURI,
		"social_links": member.SocialLinks,
		"is_active":    member.IsActive,
	}, &member.ID)

	return s.toResponseDTO(member), nil
}

// Delete menghapus member dan foto dari Cloudinary
func (s *memberService) Delete(ctx context.Context, id int) error {
	// Ambil member untuk mendapatkan info foto
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return errors.New("member tidak ditemukan")
	}

	// Log activity sebelum delete
	s.logActivity(ctx, domain.ActionDelete, domain.ModuleMembers, "Menghapus member: "+member.FullName, map[string]any{
		"id":           member.ID,
		"full_name":    member.FullName,
		"position":     member.Position,
		"department":   string(member.Department),
		"photo_uri":    member.PhotoURI,
		"social_links": member.SocialLinks,
		"is_active":    member.IsActive,
		"created_at":   member.CreatedAt,
	}, nil, &member.ID)

	// Hapus dari database
	if err := s.memberRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus member")
	}

	// Hapus foto dari Cloudinary (jika ada)
	if member.PhotoURI != nil {
		_ = s.cloudinaryService.DeleteImage(ctx, "members", *member.PhotoURI)
	}

	return nil
}

// toResponseDTO converts domain.Member to responses.MemberResponse
func (s *memberService) toResponseDTO(m *domain.Member) *responses.MemberResponse {
	var imageURL string
	if m.PhotoURI != nil {
		imageURL = s.cloudinaryService.GetImageURL("members", *m.PhotoURI)
	}

	return &responses.MemberResponse{
		ID:          m.ID,
		FullName:    m.FullName,
		Position:    m.Position,
		Department:  string(m.Department),
		Photo:       imageURL,
		SocialLinks: m.SocialLinks,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
	}
}

// logActivity helper untuk mencatat activity log
func (s *memberService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return // Skip if no user in context
	}

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
		TargetID:    targetID,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
	}

	// Ignore error - logging should not affect main operation
	_ = s.activityLogRepo.Create(log)
}
