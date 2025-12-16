package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// SiteSettingHandler handles HTTP requests untuk site settings
type SiteSettingHandler struct {
	siteSettingService service.SiteSettingService
}

// NewSiteSettingHandler constructor untuk SiteSettingHandler
func NewSiteSettingHandler(siteSettingService service.SiteSettingService) *SiteSettingHandler {
	return &SiteSettingHandler{siteSettingService: siteSettingService}
}

// Get handles GET /v1/admin/settings
func (h *SiteSettingHandler) Get(c *gin.Context) {
	result, err := h.siteSettingService.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil pengaturan situs", result))
}

// Update handles PUT /v1/admin/settings
func (h *SiteSettingHandler) Update(c *gin.Context) {
	var req requests.UpdateSiteSettingRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Get optional file uploads
	favicon, _ := c.FormFile("favicon")
	logoHeader, _ := c.FormFile("logo_header")
	logoBig, _ := c.FormFile("logo_big")

	result, err := h.siteSettingService.Update(c.Request.Context(), req, favicon, logoHeader, logoBig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Pengaturan situs berhasil diperbarui", result))
}
