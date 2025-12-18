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

func (h *PublicHomeHandler) GetAboutUsSection(c *gin.Context) {
	aboutUsSection, err := h.homeService.GetAboutUsSection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data about us section"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(http.StatusOK, "Berhasil mengambil data about us section", aboutUsSection))
}

func (h *PublicHomeHandler) GetWhySection(c *gin.Context) {
	whySection, err := h.homeService.GetWhySection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data why section"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(http.StatusOK, "Berhasil mengambil data why section", whySection))
}

func (h *PublicHomeHandler) GetTestimonialSection(c *gin.Context) {
	testimonialSection, err := h.homeService.GetTestimonialSection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data testimonial section"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(http.StatusOK, "Berhasil mengambil data testimonial section", testimonialSection))
}
