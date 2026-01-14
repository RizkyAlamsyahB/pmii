package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// PublicSiteSettingHandler handles public HTTP requests untuk site settings
type PublicSiteSettingHandler struct {
	publicSiteSettingService service.PublicSiteSettingService
}

// NewPublicSiteSettingHandler constructor untuk PublicSiteSettingHandler
func NewPublicSiteSettingHandler(publicSiteSettingService service.PublicSiteSettingService) *PublicSiteSettingHandler {
	return &PublicSiteSettingHandler{publicSiteSettingService: publicSiteSettingService}
}

// Get handles GET /v1/settings
// Returns public site settings (logo, social links, etc.)
func (h *PublicSiteSettingHandler) Get(c *gin.Context) {
	result, err := h.publicSiteSettingService.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil pengaturan situs"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil pengaturan situs", result))
}
