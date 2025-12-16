package service

import (
	"context"
	"errors"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// ContactService interface untuk business logic contact
type ContactService interface {
	Get(ctx context.Context) (*responses.ContactResponse, error)
	Update(ctx context.Context, req requests.UpdateContactRequest) (*responses.ContactResponse, error)
}

type contactService struct {
	contactRepo repository.ContactRepository
}

// NewContactService constructor untuk ContactService
func NewContactService(contactRepo repository.ContactRepository) ContactService {
	return &contactService{
		contactRepo: contactRepo,
	}
}

// Get mengambil contact info
func (s *contactService) Get(ctx context.Context) (*responses.ContactResponse, error) {
	contact, err := s.contactRepo.Get()
	if err != nil {
		// Return empty response if not found
		return &responses.ContactResponse{}, nil
	}

	return s.toResponseDTO(contact), nil
}

// Update mengupdate contact info
func (s *contactService) Update(ctx context.Context, req requests.UpdateContactRequest) (*responses.ContactResponse, error) {
	// Get existing contact
	contact, _ := s.contactRepo.Get()
	if contact == nil {
		contact = &domain.Contact{}
	}

	// Update fields
	if req.Address != nil {
		contact.Address = req.Address
	}
	if req.Email != nil {
		contact.Email = req.Email
	}
	if req.Phone != nil {
		contact.Phone = req.Phone
	}
	if req.GoogleMapsURL != nil {
		contact.GoogleMapsURL = req.GoogleMapsURL
	}

	// Save to database
	if err := s.contactRepo.Update(contact); err != nil {
		return nil, errors.New("gagal menyimpan informasi kontak")
	}

	return s.toResponseDTO(contact), nil
}

// toResponseDTO converts domain.Contact to responses.ContactResponse
func (s *contactService) toResponseDTO(contact *domain.Contact) *responses.ContactResponse {
	return &responses.ContactResponse{
		Address:       contact.Address,
		Email:         contact.Email,
		Phone:         contact.Phone,
		GoogleMapsURL: contact.GoogleMapsURL,
		UpdatedAt:     contact.UpdatedAt,
	}
}
