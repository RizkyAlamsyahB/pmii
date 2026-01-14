package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// 1. GET ALL CATEGORIES (With Pagination & Search)
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// Mengambil query params sesuai index.json
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	data, lastPage, total, err := h.svc.GetAll(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data kategori"))
		return
	}

	// Response menggunakan SuccessResponseWithPagination agar konsisten
	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(
		200,
		"List of categories",
		data,
		page,
		limit,
		total,
		lastPage,
	))
}

// 2. CREATE CATEGORY
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req requests.CategoryRequest
	// Bind dari x-www-form-urlencoded
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Nama kategori wajib diisi"))
		return
	}

	res, err := h.svc.Create(GetContextWithRequestInfo(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan kategori"))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Kategori berhasil dibuat", res))
}

// 3. UPDATE CATEGORY
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req requests.CategoryRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Input tidak valid"))
		return
	}

	res, err := h.svc.Update(GetContextWithRequestInfo(c), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Kategori tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Kategori berhasil diupdate", res))
}

// 4. DELETE CATEGORY
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(GetContextWithRequestInfo(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menghapus kategori. Pastikan tidak ada rilis berita yang menggunakan kategori ini."))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Kategori berhasil dihapus", nil))
}
