package handlers

import (
	"context"
	"strings"

	"github.com/garuda-labs-1/pmii-be/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetContextWithRequestInfo creates a context.Context with user ID, IP address and user agent
// extracted from gin.Context. This should be used when calling services that need activity logging.
func GetContextWithRequestInfo(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// Add user ID if available (after auth middleware)
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			ctx = utils.WithUserID(ctx, id)
		}
	}

	// Add IP address
	if ip, exists := c.Get(string(utils.ContextKeyIPAddress)); exists {
		if ipStr, ok := ip.(string); ok {
			ctx = context.WithValue(ctx, utils.ContextKeyIPAddress, ipStr)
		}
	} else {
		// Fallback to ClientIP if middleware not applied
		ctx = context.WithValue(ctx, utils.ContextKeyIPAddress, c.ClientIP())
	}

	// Add user agent
	if ua, exists := c.Get(string(utils.ContextKeyUserAgent)); exists {
		if uaStr, ok := ua.(string); ok {
			ctx = context.WithValue(ctx, utils.ContextKeyUserAgent, uaStr)
		}
	} else {
		// Fallback to header if middleware not applied
		ctx = context.WithValue(ctx, utils.ContextKeyUserAgent, c.GetHeader("User-Agent"))
	}

	return ctx
}

// ValidationErrors mengubah validator.ValidationErrors ke map[string]string
// untuk response yang lebih readable oleh frontend
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			// Convert field name ke snake_case (sesuai form field)
			fieldName := toSnakeCase(fieldError.Field())
			errors[fieldName] = getValidationMessage(fieldError)
		}
	}

	return errors
}

// toSnakeCase convert PascalCase/camelCase ke snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getValidationMessage menghasilkan pesan error yang readable
func getValidationMessage(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())

	switch fe.Tag() {
	case "required":
		return field + " wajib diisi"
	case "email":
		return field + " harus berupa email yang valid"
	case "min":
		return field + " minimal " + fe.Param() + " karakter"
	case "max":
		return field + " maksimal " + fe.Param() + " karakter"
	case "oneof":
		return field + " harus salah satu dari: " + fe.Param()
	case "url":
		return field + " harus berupa URL yang valid"
	case "numeric":
		return field + " harus berupa angka"
	case "alphanum":
		return field + " hanya boleh huruf dan angka"
	default:
		return field + " tidak valid"
	}
}

// getRoleName convert role code ke nama role
func getRoleName(role int) string {
	switch role {
	case 1:
		return "admin"
	case 2:
		return "author"
	default:
		return "author"
	}
}

// getStatusName convert isActive bool ke nama status
func getStatusName(isActive bool) string {
	if isActive {
		return "active"
	}
	return "inactive"
}
