package service

import (
	"fmt"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

type InboxService interface {
	GetList(userID int) ([]responses.InboxListItemResponse, error)
	GetChatHistory(senderID, receiverID int) ([]responses.ChatHistoryResponse, error)
}

type inboxService struct {
	repo     repository.InboxRepository
	userRepo repository.UserRepository // Untuk mengambil detail profil lawan chat
}

func NewInboxService(repo repository.InboxRepository, userRepo repository.UserRepository) InboxService {
	return &inboxService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GetList mengembalikan daftar orang yang pernah chat dengan user (untuk Modal Inbox)
func (s *inboxService) GetList(userID int) ([]responses.InboxListItemResponse, error) {
	// 1. Ambil pesan terakhir dari setiap percakapan unik
	messages, err := s.repo.GetLatestMessagesPerUser(userID)
	if err != nil {
		return nil, err
	}

	var result []responses.InboxListItemResponse
	for _, msg := range messages {
		// Tentukan siapa "lawan bicara" (bukan diri sendiri)
		opponentID := msg.SenderID
		if msg.SenderID == userID {
			opponentID = msg.ReceiverID
		}

		// Ambil info user lawan bicara
		opponent, _ := s.userRepo.FindByID(opponentID)

		// Handling pointer to string untuk PhotoURI
		photoStr := ""
		if opponent.PhotoURI != nil {
			photoStr = *opponent.PhotoURI
		}

		result = append(result, responses.InboxListItemResponse{
			UserID:      opponentID,
			FullName:    opponent.FullName,
			PhotoURI:    photoStr,
			LastMessage: msg.Message,
			Time:        msg.CreatedAt.Format("3:04 PM"), // Format sesuai desain "9:56 PM"
			UnreadCount: s.repo.CountUnread(opponentID, userID),
		})
	}

	return result, nil
}

// GetChatHistory mengembalikan riwayat bubble chat antara dua user (untuk Modal Chat)
func (s *inboxService) GetChatHistory(senderID, receiverID int) ([]responses.ChatHistoryResponse, error) {
	messages, err := s.repo.GetMessagesBetweenUsers(senderID, receiverID)
	if err != nil {
		return nil, err
	}

	var result []responses.ChatHistoryResponse
	for _, msg := range messages {
		result = append(result, responses.ChatHistoryResponse{
			ID:       msg.ID,
			SenderID: msg.SenderID,
			Message:  msg.Message,
			Time:     msg.CreatedAt.Format("15:04"),
			Date:     formatIndonesianDate(msg.CreatedAt), // Helper untuk "17 Desember 2025"
		})
	}

	// Tandai pesan sebagai terbaca saat riwayat dibuka
	go s.repo.MarkAsRead(receiverID, senderID)

	return result, nil
}

// Helper untuk format tanggal Indonesia
func formatIndonesianDate(t time.Time) string {
	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
}
