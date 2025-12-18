package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// TestimonialHandler handles HTTP requests untuk testimonial
type TestimonialHandler struct {
	testimonialService service.TestimonialService
}

// NewTestimonialHandler constructor untuk TestimonialHandler
func NewTestimonialHandler(testimonialService service.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{testimonialService: testimonialService}
}

// Create handles POST /v1/testimonials
func (h *TestimonialHandler) Create(c *gin.Context) {
	var req requests.CreateTestimonialRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Get photo file
	photoFile, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Format foto tidak valid"))
		return
	}

	// Validate photo size (max 5MB)
	if photoFile != nil && photoFile.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Ukuran foto maksimal 5MB"))
		return
	}

	// Call service
	testimonial, err := h.testimonialService.Create(c.Request.Context(), req, photoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Testimonial berhasil dibuat", testimonial))
}

// GetAll handles GET /v1/testimonials with pagination
func (h *TestimonialHandler) GetAll(c *gin.Context) {
	// Parse query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	testimonials, currentPage, lastPage, total, err := h.testimonialService.GetAll(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(200, "List of testimonials", testimonials, currentPage, limit, total, lastPage))
}

// GetByID handles GET /v1/testimonials/:id
func (h *TestimonialHandler) GetByID(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	testimonial, err := h.testimonialService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil data testimonial", testimonial))
}

// Update handles PUT /v1/testimonials/:id
func (h *TestimonialHandler) Update(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	var req requests.UpdateTestimonialRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Get photo file (optional)
	photoFile, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Format foto tidak valid"))
		return
	}

	// Validate photo size (max 5MB)
	if photoFile != nil && photoFile.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Ukuran foto maksimal 5MB"))
		return
	}

	// Call service
	testimonial, err := h.testimonialService.Update(c.Request.Context(), id, req, photoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Testimonial berhasil diupdate", testimonial))
}

// Delete handles DELETE /v1/testimonials/:id
func (h *TestimonialHandler) Delete(c *gin.Context) {
	// Parse ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	// Call service
	if err := h.testimonialService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Testimonial berhasil dihapus", nil))
}
