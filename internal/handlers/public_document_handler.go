package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// PublicDocumentHandler handles public HTTP requests untuk document
type PublicDocumentHandler struct {
	publicDocumentService service.PublicDocumentService
}

// NewPublicDocumentHandler constructor untuk PublicDocumentHandler
func NewPublicDocumentHandler(publicDocumentService service.PublicDocumentService) *PublicDocumentHandler {
	return &PublicDocumentHandler{publicDocumentService: publicDocumentService}
}

// GetAllPublic handles GET /v1/documents
// Returns all documents grouped by type
func (h *PublicDocumentHandler) GetAllPublic(c *gin.Context) {
	result, err := h.publicDocumentService.GetAllPublic(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data dokumen"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil data dokumen", result))
}

// GetByTypePublic handles GET /v1/documents/:type
// Returns documents by specific type
func (h *PublicDocumentHandler) GetByTypePublic(c *gin.Context) {
	fileType := c.Param("type")

	result, err := h.publicDocumentService.GetByTypePublic(c.Request.Context(), fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data dokumen"))
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Jenis dokumen tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil data dokumen", result))
}
