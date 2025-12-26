package handlers

import (
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc service.TagService
}

func NewTagHandler(svc service.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	data, lastPage, total, err := h.svc.GetAll(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data tags"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(200, "List of tags", data, page, limit, total, lastPage))
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req requests.TagRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Nama tag wajib diisi"))
		return
	}

	res, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menyimpan tag"))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Tag berhasil dibuat", res))
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var req requests.TagRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Input tidak valid"))
		return
	}

	res, err := h.svc.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, "Tag tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Tag berhasil diupdate", res))
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal menghapus tag"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Tag berhasil dihapus", nil))
}
