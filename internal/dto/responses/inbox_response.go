package responses

type InboxListItemResponse struct {
	UserID      int    `json:"user_id"`
	FullName    string `json:"full_name"`
	PhotoURI    string `json:"photo_uri"`
	LastMessage string `json:"last_message"`
	Time        string `json:"time"` // Format: "9.56 PM"
	UnreadCount int    `json:"unread_count"`
}

type ChatHistoryResponse struct {
	ID       int    `json:"id"`
	SenderID int    `json:"sender_id"`
	Message  string `json:"message"`
	Time     string `json:"time"`
	Date     string `json:"date"` // Untuk pemisah "17 Desember 2025"
}
