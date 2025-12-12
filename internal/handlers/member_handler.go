package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService service.MemberService
}

// NewMemberHandler creates a new MemberHandler
func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// Create handles POST /v1/admin/members
func (h *MemberHandler) Create(c *gin.Context) {
	var req requests.CreateMemberRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Handle social_links JSON if provided
	if socialLinksStr := c.PostForm("social_links"); socialLinksStr != "" {
		var socialLinks map[string]any
		if err := json.Unmarshal([]byte(socialLinksStr), &socialLinks); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Format social_links JSON tidak valid"))
			return
		}
		req.SocialLinks = socialLinks
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
	member, err := h.memberService.Create(c.Request.Context(), req, photoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(201, "Member berhasil dibuat", member))
}

// GetAll handles GET /v1/admin/members
func (h *MemberHandler) GetAll(c *gin.Context) {
	members, err := h.memberService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data member berhasil diambil", members))
}

// GetByID handles GET /v1/admin/members/:id
func (h *MemberHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID member tidak valid"))
		return
	}

	member, err := h.memberService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data member berhasil diambil", member))
}

// Update handles PUT /v1/admin/members/:id
func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID member tidak valid"))
		return
	}

	var req requests.UpdateMemberRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Data tidak valid"))
		return
	}

	// Handle social_links JSON if provided
	if socialLinksStr := c.PostForm("social_links"); socialLinksStr != "" {
		var socialLinks map[string]any
		if err := json.Unmarshal([]byte(socialLinksStr), &socialLinks); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "Format social_links JSON tidak valid"))
			return
		}
		req.SocialLinks = socialLinks
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
	member, err := h.memberService.Update(c.Request.Context(), id, req, photoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Member berhasil diupdate", member))
}

// Delete handles DELETE /v1/admin/members/:id
func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse(400, "ID member tidak valid"))
		return
	}

	if err := h.memberService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Member berhasil dihapus", nil))
}
