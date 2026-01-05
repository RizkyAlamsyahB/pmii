package handlers

import (
	"net/http"
	"strconv"

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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Logika untuk mendaftarkan koneksi dan mendengarkan pesan masuk
	// Pesan akan disimpan ke PostgreSQL melalui service
}

// GET /v1/inbox
func (h *InboxHandler) GetInboxList(c *gin.Context) {
	// Mengambil user_id dari middleware auth
	userID := c.MustGet("userID").(int)
	data, err := h.svc.GetList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *InboxHandler) GetChatHistory(c *gin.Context) {
	receiverID, _ := strconv.Atoi(c.Param("user_id"))
	senderID := c.MustGet("userID").(int) // Mengambil ID dari token auth

	data, err := h.svc.GetChatHistory(senderID, receiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal memuat riwayat chat"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Success", data))
}
