package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// PublicAdHandler handles public HTTP requests for ads
type PublicAdHandler struct {
	publicAdService service.PublicAdService
}

// NewPublicAdHandler constructor untuk PublicAdHandler
func NewPublicAdHandler(publicAdService service.PublicAdService) *PublicAdHandler {
	return &PublicAdHandler{publicAdService: publicAdService}
}

// GetAdsByPage handles GET /ads/:page
// @Summary Get active ads by page (public)
// @Description Get all active advertisement slots for a specific page (no auth required)
// @Tags Public Ads
// @Accept json
// @Produce json
// @Param page path string true "Page name (landing, news, opini, life_at_pmii, islamic, detail_article)"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Router /ads/{page} [get]
func (h *PublicAdHandler) GetAdsByPage(c *gin.Context) {
	page := c.Param("page")

	ads, err := h.publicAdService.GetAdsByPage(c.Request.Context(), page)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data ads berhasil diambil", ads))
}

// GetAvailablePages handles GET /ads/pages
// @Summary Get available ad pages (public)
// @Description Get list of all available ad pages
// @Tags Public Ads
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Router /ads/pages [get]
func (h *PublicAdHandler) GetAvailablePages(c *gin.Context) {
	pages := make([]map[string]string, 0)

	for _, page := range domain.ValidAdPages() {
		ad := &domain.Ad{Page: page}
		pages = append(pages, map[string]string{
			"page":      string(page),
			"page_name": ad.GetPageDisplayName(),
		})
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Daftar halaman ads", pages))
}
