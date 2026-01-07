package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// 1. GET ALL POSTS (With Pagination & Search)
func (h *PostHandler) GetPosts(c *gin.Context) {
	// Parsing query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	// Memanggil service layer
	data, lastPage, total, err := h.svc.GetAllPosts(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data berita"))
		return
	}

	// Response dengan pagination sesuai base_response.go
	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(
		200,
		"List of posts",
		data,
		page,
		limit,
		total,
		lastPage,
	))
}

// 2. CREATE POST (Multipart/Form-Data)
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req requests.PostCreateRequest

	// Bind form data ke struct request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Input tidak valid: "+err.Error()))
		return
	}

	// Ambil file gambar dari context Gin
	file, _ := c.FormFile("image")
	req.Image = file

	// Eksekusi pembuatan post melalui service
	res, err := h.svc.CreatePost(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal membuat berita"))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Berita berhasil dibuat", res))
}

// 3. GET DETAIL POST (By ID or Slug)
func (h *PostHandler) GetPost(c *gin.Context) {
	identifier := c.Param("slug")
	if identifier == "" {
		identifier = c.Param("id")
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	res, err := h.svc.GetPostDetail(identifier, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Detail berita ditemukan", res))
}

// 4. UPDATE POST
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req requests.PostUpdateRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Input tidak valid"))
		return
	}

	// Ambil file gambar jika ada update foto
	file, _ := c.FormFile("image")
	req.Image = file

	res, err := h.svc.UpdatePost(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal memperbarui berita"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil diupdate", res))
}

// 5. DELETE POST (Soft Delete)
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeletePost(id); err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan atau gagal dihapus"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berita berhasil dihapus", nil))
}
