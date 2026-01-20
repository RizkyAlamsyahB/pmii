package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// AdHandler handles HTTP requests for ads management
type AdHandler struct {
	adService service.AdService
}

// NewAdHandler constructor untuk AdHandler
func NewAdHandler(adService service.AdService) *AdHandler {
	return &AdHandler{adService: adService}
}

// GetAllAds handles GET /admin/ads
// @Summary Get all ads grouped by page
// @Description Get all advertisement slots grouped by page for admin management
// @Tags Ads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /admin/ads [get]
func (h *AdHandler) GetAllAds(c *gin.Context) {
	ctx := GetContextWithRequestInfo(c)

	ads, err := h.adService.GetAllAds(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data ads"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data ads berhasil diambil", ads))
}

// GetAdByID handles GET /admin/ads/:id
// @Summary Get ad by ID
// @Description Get single advertisement slot by ID
// @Tags Ads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ad ID"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /admin/ads/{id} [get]
func (h *AdHandler) GetAdByID(c *gin.Context) {
	ctx := GetContextWithRequestInfo(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	ad, err := h.adService.GetAdByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data ad berhasil diambil", ad))
}

// GetAdsByPage handles GET /admin/ads/page/:page
// @Summary Get ads by page
// @Description Get all advertisement slots for a specific page
// @Tags Ads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page path string true "Page name (landing, news, opini, life_at_pmii, islamic, detail_article)"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Router /admin/ads/page/{page} [get]
func (h *AdHandler) GetAdsByPage(c *gin.Context) {
	ctx := GetContextWithRequestInfo(c)

	page := c.Param("page")

	ads, err := h.adService.GetAdsByPage(ctx, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data ads berhasil diambil", ads))
}

// UpdateAd handles PUT /admin/ads/:id
// @Summary Update ad image
// @Description Update advertisement slot image
// @Tags Ads
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ad ID"
// @Param image formData file true "Ad image"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /admin/ads/{id} [put]
func (h *AdHandler) UpdateAd(c *gin.Context) {
	ctx := GetContextWithRequestInfo(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	// Get image file (required)
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Gambar wajib diupload"))
		return
	}

	ad, err := h.adService.UpdateAd(ctx, id, image)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Ad berhasil diupdate", ad))
}

// DeleteAdImage handles DELETE /admin/ads/:id/image
// @Summary Delete ad image
// @Description Delete the image of an advertisement slot (set to null)
// @Tags Ads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Ad ID"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /admin/ads/{id}/image [delete]
func (h *AdHandler) DeleteAdImage(c *gin.Context) {
	ctx := GetContextWithRequestInfo(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	ad, err := h.adService.DeleteAdImage(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Gambar ad berhasil dihapus", ad))
}
