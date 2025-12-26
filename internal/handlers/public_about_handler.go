package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// PublicAboutHandler handles public HTTP requests untuk about page
type PublicAboutHandler struct {
	publicAboutService service.PublicAboutService
}

// NewPublicAboutHandler constructor untuk PublicAboutHandler
func NewPublicAboutHandler(publicAboutService service.PublicAboutService) *PublicAboutHandler {
	return &PublicAboutHandler{publicAboutService: publicAboutService}
}

// GetAboutPage handles GET /v1/about
// Returns about info + members grouped by department (first page each)
func (h *PublicAboutHandler) GetAboutPage(c *gin.Context) {
	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))
	search := c.Query("search")

	result, err := h.publicAboutService.GetAboutPage(c.Request.Context(), limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data about"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil data about", result))
}

// GetMembersByDepartment handles GET /v1/about/members/:department
// Returns members by department with pagination & search
func (h *PublicAboutHandler) GetMembersByDepartment(c *gin.Context) {
	department := c.Param("department")

	// Validate department
	validDepartment := false
	for _, d := range domain.ValidDepartments() {
		if string(d) == department {
			validDepartment = true
			break
		}
	}
	if !validDepartment {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Department tidak valid. Pilih: pengurus_harian, kabid, wasekbid, wakil_bendahara"))
		return
	}

	// Parse query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))
	search := c.Query("search")

	members, currentPage, lastPage, total, err := h.publicAboutService.GetMembersByDepartment(c.Request.Context(), department, page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data members"))
		return
	}

	// Get department label
	deptLabel := domain.MemberDepartment(department).GetLabel()

	// Build data response
	data := gin.H{
		"department":      department,
		"departmentLabel": deptLabel,
		"members":         members,
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(
		200,
		"List of "+deptLabel,
		data,
		currentPage,
		limit,
		total,
		lastPage,
	))
}

// GetDepartments handles GET /v1/about/departments
// Returns list of available departments
func (h *PublicAboutHandler) GetDepartments(c *gin.Context) {
	departments := make([]gin.H, 0, 4)
	for _, d := range domain.ValidDepartments() {
		departments = append(departments, gin.H{
			"value": string(d),
			"label": d.GetLabel(),
		})
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "List of departments", departments))
}
