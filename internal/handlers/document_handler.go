package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// DocumentHandler handles HTTP requests untuk document
type DocumentHandler struct {
	documentService service.DocumentService
}

// NewDocumentHandler constructor untuk DocumentHandler
func NewDocumentHandler(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

// Create handles POST /v1/admin/documents
func (h *DocumentHandler) Create(c *gin.Context) {
	// Set timeout 30 menit untuk upload file besar
	ctx, cancel := context.WithTimeout(GetContextWithRequestInfo(c), 30*time.Minute)
	defer cancel()

	var req requests.CreateDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		errors := FormatValidationErrors(err)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(errors))
			return
		}
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Nama dan jenis file wajib diisi"))
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "File wajib diupload"))
		return
	}

	result, err := h.documentService.Create(ctx, req, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Dokumen berhasil ditambahkan", result))
}

// GetAll handles GET /v1/admin/documents
func (h *DocumentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	fileType := c.Query("file_type")
	search := c.Query("search")

	documents, currentPage, lastPage, total, err := h.documentService.GetAll(c.Request.Context(), page, limit, fileType, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(200, "List of documents", documents, currentPage, limit, total, lastPage))
}

// GetByID handles GET /v1/admin/documents/:id
func (h *DocumentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	result, err := h.documentService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Berhasil mengambil dokumen", result))
}

// Update handles PUT /v1/admin/documents/:id
func (h *DocumentHandler) Update(c *gin.Context) {
	// Set timeout 2 menit untuk upload file besar
	ctx, cancel := context.WithTimeout(GetContextWithRequestInfo(c), 2*time.Minute)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	var req requests.UpdateDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		errors := FormatValidationErrors(err)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse(errors))
			return
		}
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Get optional file from form
	file, _ := c.FormFile("file")

	result, err := h.documentService.Update(ctx, id, req, file)
	if err != nil {
		if err.Error() == "dokumen tidak ditemukan" {
			c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Dokumen berhasil diperbarui", result))
}

// Delete handles DELETE /v1/admin/documents/:id
func (h *DocumentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID tidak valid"))
		return
	}

	if err := h.documentService.Delete(GetContextWithRequestInfo(c), id); err != nil {
		if err.Error() == "dokumen tidak ditemukan" {
			c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Dokumen berhasil dihapus", nil))
}

// GetTypes handles GET /v1/admin/documents/types
func (h *DocumentHandler) GetTypes(c *gin.Context) {
	types := h.documentService.GetDocumentTypes()
	c.JSON(http.StatusOK, responses.SuccessResponse(200, "List of document types", types))
}
