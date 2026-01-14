package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type InboxHandler struct {
	svc service.InboxService
}

func NewInboxHandler(svc service.InboxService) *InboxHandler {
	return &InboxHandler{svc: svc}
}

// WS /v1/chat/ws
func (h *InboxHandler) HandleWebSocket(c *gin.Context) {
	// 1. Ambil Sender ID dari context (Pastikan kunci "user_id" sesuai middleware)
	senderIDVal, exists := c.Get("user_id")
	if !exists {
		log.Println("WS Error: user_id not found in context")
		return
	}
	senderID := senderIDVal.(int)

	// 2. Upgrade HTTP ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WS Upgrade Error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("User %d connected via WebSocket", senderID)

	// 3. Looping Forever untuk mendengarkan pesan masuk agar koneksi tidak putus
	for {
		var msgReq struct {
			ReceiverID int    `json:"receiver_id"`
			Message    string `json:"message"`
		}

		// Membaca pesan JSON dari Postman
		err := conn.ReadJSON(&msgReq)
		if err != nil {
			log.Printf("WS Read Error (User %d disconnected): %v", senderID, err)
			break // Keluar dari loop jika koneksi ditutup atau error
		}

		// 4. Bungkus ke model domain untuk disimpan
		newInbox := &domain.Inbox{
			SenderID:   senderID,
			ReceiverID: msgReq.ReceiverID,
			Message:    msgReq.Message,
			IsRead:     false,
		}

		// 5. Simpan ke Database melalui Service
		err = h.svc.SendMessage(newInbox)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			conn.WriteJSON(gin.H{"error": "Failed to save message"})
			continue
		}

		// 6. Feedback ke pengirim (Optional)
		conn.WriteJSON(gin.H{
			"status":  "sent",
			"message": "Pesan berhasil disimpan ke database",
		})
	}
}

// GET /v1/inbox
func (h *InboxHandler) GetInboxList(c *gin.Context) {
	// Mengambil user_id dari middleware auth
	userID := c.MustGet("user_id").(int)
	data, err := h.svc.GetList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GET /v1/chat/:user_id
func (h *InboxHandler) GetChatHistory(c *gin.Context) {
	receiverID, _ := strconv.Atoi(c.Param("user_id"))
	senderID := c.MustGet("user_id").(int) // Mengambil ID dari token auth

	data, err := h.svc.GetChatHistory(senderID, receiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal memuat riwayat chat"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Success", data))
}
