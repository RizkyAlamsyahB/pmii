package handlers

import (
	"net/http"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// AboutHandler handles HTTP requests untuk about page
type AboutHandler struct {
	aboutService service.AboutService
}

// NewAboutHandler constructor untuk AboutHandler
func NewAboutHandler(aboutService service.AboutService) *AboutHandler {
	return &AboutHandler{aboutService: aboutService}
}

// Get handles GET /v1/admin/about
func (h *AboutHandler) Get(c *gin.Context) {
	about, err := h.aboutService.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil data about", about))
}

// Update handles PUT /v1/admin/about
func (h *AboutHandler) Update(c *gin.Context) {
	var req requests.UpdateAboutRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Get image file (optional)
	imageFile, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Format gambar tidak valid"))
		return
	}

	// Validate image size (max 5MB)
	if imageFile != nil && imageFile.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Ukuran gambar maksimal 5MB"))
		return
	}

	// Call service
	about, err := h.aboutService.Update(c.Request.Context(), req, imageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "About berhasil diupdate", about))
}
