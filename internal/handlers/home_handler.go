package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type PublicHomeHandler struct {
	homeService service.HomeService
}

func NewPublicHomeHandler(homeService service.HomeService) *PublicHomeHandler {
	return &PublicHomeHandler{homeService: homeService}
}

func (h *PublicHomeHandler) GetHeroSection(c *gin.Context) {
	heroSection, err := h.homeService.GetHeroSection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data hero section"))
		return
	}

	// Return empty array [] if no posts exist (not an error)
	c.JSON(http.StatusOK, responses.SuccessResponse(http.StatusOK, "Berhasil mengambil data hero section", heroSection))
}

func (h *PublicHomeHandler) GetLatestNewsSection(c *gin.Context) {
	latestNewsSection, err := h.homeService.GetLatestNewsSection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data latest news section"))
		return
	}

	// Return empty array [] if no posts exist (not an error)
	c.JSON(http.StatusOK, responses.SuccessResponse(http.StatusOK, "Berhasil mengambil data latest news section", latestNewsSection))
}
