package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	svc service.NewsService
}

func NewNewsHandler(svc service.NewsService) *NewsHandler {
	return &NewsHandler{svc: svc}
}

func (h *NewsHandler) GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	data, lastPage, total, err := h.svc.FetchPublicNews(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data berita"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(
		200, "Berita berhasil dimuat", data, page, limit, total, lastPage,
	))
}

func (h *NewsHandler) GetNewsDetail(c *gin.Context) {
	slug := c.Param("slug")
	data, err := h.svc.FetchNewsDetail(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Berita tidak ditemukan"))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Detail berita ditemukan", data))
}
