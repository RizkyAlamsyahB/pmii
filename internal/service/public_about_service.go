package service

import (
	"context"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// PublicAboutService interface untuk public about page
type PublicAboutService interface {
	GetAboutPage(ctx context.Context, limit int, search string) (*responses.PublicAboutPageResponse, error)
	GetMembersByDepartment(ctx context.Context, department string, page, limit int, search string) ([]responses.PublicMemberResponse, int, int, int64, error)
}

type publicAboutService struct {
	aboutRepo         repository.AboutRepository
	memberRepo        repository.MemberRepository
	contactRepo       repository.ContactRepository
	cloudinaryService CloudinaryService
}

// NewPublicAboutService constructor untuk PublicAboutService
func NewPublicAboutService(
	aboutRepo repository.AboutRepository,
	memberRepo repository.MemberRepository,
	contactRepo repository.ContactRepository,
	cloudinaryService CloudinaryService,
) PublicAboutService {
	return &publicAboutService{
		aboutRepo:         aboutRepo,
		memberRepo:        memberRepo,
		contactRepo:       contactRepo,
		cloudinaryService: cloudinaryService,
	}
}

// GetAboutPage mengambil data about page lengkap dengan members per department dan contact
func (s *publicAboutService) GetAboutPage(ctx context.Context, limit int, search string) (*responses.PublicAboutPageResponse, error) {
	// Set default limit
	if limit < 1 {
		limit = 8 // Default 8 members per department sesuai design
	}

	// Get about data
	about, _ := s.aboutRepo.Get() // Ignore error, return empty if not found

	// Get contact data
	contact, _ := s.contactRepo.Get() // Ignore error, return empty if not found

	// Build departments response
	departments := make([]responses.DepartmentMembersResponse, 0, 4)

	// Iterate through all valid departments
	for _, dept := range domain.ValidDepartments() {
		members, total, err := s.memberRepo.FindActiveByDepartment(string(dept), 1, limit, search)
		if err != nil {
			continue // Skip department on error
		}

		// Jika ada search query, skip department yang tidak memiliki hasil
		if search != "" && total == 0 {
			continue
		}

		// Calculate last page
		lastPage := int(total) / limit
		if int(total)%limit != 0 {
			lastPage++
		}
		if lastPage < 1 {
			lastPage = 1
		}

		deptResponse := responses.DepartmentMembersResponse{
			Department:      string(dept),
			DepartmentLabel: dept.GetLabel(),
			Members:         s.toPublicMembersResponse(members, 1, limit, total, lastPage),
		}
		departments = append(departments, deptResponse)
	}

	// Build response
	response := &responses.PublicAboutPageResponse{
		About:       s.toPublicAboutResponse(about),
		Departments: departments,
		Contact:     s.toPublicContactResponse(contact),
	}

	return response, nil
}

// GetMembersByDepartment mengambil list members per department (dengan pagination & search)
func (s *publicAboutService) GetMembersByDepartment(ctx context.Context, department string, page, limit int, search string) ([]responses.PublicMemberResponse, int, int, int64, error) {
	// Set default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 8
	}

	members, total, err := s.memberRepo.FindActiveByDepartment(department, page, limit, search)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Calculate last page
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}
	if lastPage < 1 {
		lastPage = 1
	}

	// Convert to response
	result := make([]responses.PublicMemberResponse, len(members))
	for i, m := range members {
		result[i] = s.toPublicMemberResponse(&m)
	}

	return result, page, lastPage, total, nil
}

// toPublicAboutResponse converts domain.About to PublicAboutResponse
func (s *publicAboutService) toPublicAboutResponse(a *domain.About) responses.PublicAboutResponse {
	if a == nil {
		return responses.PublicAboutResponse{}
	}

	var imageURL string
	if a.ImageURI != nil {
		imageURL = s.cloudinaryService.GetImageURL("about", *a.ImageURI)
	}

	return responses.PublicAboutResponse{
		History:  a.History,
		Vision:   a.Vision,
		Mission:  a.Mission,
		ImageUrl: imageURL,
		VideoURL: a.VideoURL,
	}
}

// toPublicMemberResponse converts domain.Member to PublicMemberResponse
func (s *publicAboutService) toPublicMemberResponse(m *domain.Member) responses.PublicMemberResponse {
	var photoURL string
	if m.PhotoURI != nil {
		photoURL = s.cloudinaryService.GetImageURL("members", *m.PhotoURI)
	}

	return responses.PublicMemberResponse{
		ID:          m.ID,
		FullName:    m.FullName,
		Position:    m.Position,
		Photo:       photoURL,
		SocialLinks: m.SocialLinks,
	}
}

// toPublicMembersResponse builds the members response with pagination
func (s *publicAboutService) toPublicMembersResponse(members []domain.Member, page, limit int, total int64, lastPage int) responses.PublicMembersResponse {
	data := make([]responses.PublicMemberResponse, len(members))
	for i, m := range members {
		data[i] = s.toPublicMemberResponse(&m)
	}

	return responses.PublicMembersResponse{
		Data: data,
		Pagination: responses.PaginationMeta{
			Page:     page,
			Limit:    limit,
			Total:    total,
			LastPage: lastPage,
		},
	}
}

// toPublicContactResponse converts domain.Contact to PublicContactResponse
func (s *publicAboutService) toPublicContactResponse(c *domain.Contact) responses.PublicContactResponse {
	if c == nil {
		return responses.PublicContactResponse{}
	}

	return responses.PublicContactResponse{
		Address:       c.Address,
		Email:         c.Email,
		Phone:         c.Phone,
		GoogleMapsURL: c.GoogleMapsURL,
	}
}
