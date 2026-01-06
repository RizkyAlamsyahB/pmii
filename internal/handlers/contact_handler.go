package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// ContactHandler handles HTTP requests untuk contact info
type ContactHandler struct {
	contactService service.ContactService
}

// NewContactHandler constructor untuk ContactHandler
func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

// Get handles GET /v1/admin/contact
func (h *ContactHandler) Get(c *gin.Context) {
	result, err := h.contactService.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil informasi kontak", result))
}

// Update handles PUT /v1/admin/contact
func (h *ContactHandler) Update(c *gin.Context) {
	var req requests.UpdateContactRequest

	// Bind berdasarkan Content-Type
	contentType := c.ContentType()
	var err error
	if contentType == "application/json" {
		err = c.ShouldBindJSON(&req)
	} else {
		err = c.ShouldBind(&req)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	result, err := h.contactService.Update(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Informasi kontak berhasil diperbarui", result))
}
