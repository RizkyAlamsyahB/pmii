package repository

import (
	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

type InboxRepository interface {
	GetLatestMessagesPerUser(userID int) ([]domain.Inbox, error)
	GetMessagesBetweenUsers(userID1, userID2 int) ([]domain.Inbox, error)
	CountUnread(senderID, receiverID int) int
	MarkAsRead(senderID, receiverID int) error
	Create(message *domain.Inbox) error
}

type inboxRepository struct{}

func NewInboxRepository() InboxRepository {
	return &inboxRepository{}
}

// GetLatestMessagesPerUser mengambil pesan terakhir dari setiap percakapan unik
func (r *inboxRepository) GetLatestMessagesPerUser(userID int) ([]domain.Inbox, error) {
	var messages []domain.Inbox

	// Query ini mengelompokkan pesan berdasarkan pasangan chat dan mengambil ID terbaru
	subQuery := config.DB.Table("inbox").
		Select("MAX(id)").
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)")

	err := config.DB.Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&messages).Error

	return messages, err
}

// GetMessagesBetweenUsers mengambil riwayat chat lengkap antara dua user
func (r *inboxRepository) GetMessagesBetweenUsers(u1, u2 int) ([]domain.Inbox, error) {
	var messages []domain.Inbox
	err := config.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", u1, u2, u2, u1).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *inboxRepository) CountUnread(senderID, receiverID int) int {
	var count int64
	config.DB.Model(&domain.Inbox{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Count(&count)
	return int(count)
}

func (r *inboxRepository) MarkAsRead(senderID, receiverID int) error {
	return config.DB.Model(&domain.Inbox{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Update("is_read", true).Error
}

func (r *inboxRepository) Create(message *domain.Inbox) error {
	return config.DB.Create(message).Error
}
