package responses

import "math"

// Response adalah struktur response standar API
type Response struct {
	Meta   Meta        `json:"meta"`
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

// Meta adalah metadata response
type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PaginationMeta struct untuk meta pagination
type PaginationMeta struct {
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Total    int64 `json:"total"`
	LastPage int   `json:"lastPage"`
}

// MetaWithPagination struct wrapper agar Pagination masuk ke dalam object Meta
type MetaWithPagination struct {
	Meta
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationResponse struct khusus untuk response list dengan pagination
type PaginationResponse struct {
	Meta MetaWithPagination `json:"meta"`
	Data interface{}        `json:"data"`
}

// SuccessResponse membuat response sukses
func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Status:  "success",
			Message: message,
		},
		Data: data,
	}
}

// SuccessPaginationResponse membuat response sukses dengan metadata pagination
func SuccessResponseWithPagination(code int, message string, page int, limit int, total int64, data interface{}) PaginationResponse {
	// Hitung last page (Total data dibagi limit, dibulatkan ke atas)
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	// Jika total 0, lastPage minimal 0 atau 1 (tergantung preferensi, di sini 0 jika kosong)
	if total == 0 {
		lastPage = 1
	}

	return PaginationResponse{
		Meta: MetaWithPagination{
			Meta: Meta{
				Code:    code,
				Status:  "success",
				Message: message,
			},
			Pagination: PaginationMeta{
				Page:     page,
				Limit:    limit,
				Total:    total,
				LastPage: lastPage,
			},
		},
		Data: data,
	}
}

// ErrorResponse membuat response error
func ErrorResponse(code int, message string) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Status:  "error",
			Message: message,
		},
		Data: nil,
	}
}

// ValidationErrorResponse membuat response untuk validation error
func ValidationErrorResponse(errors interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    400,
			Status:  "error",
			Message: "Validasi gagal",
		},
		Errors: errors,
	}
}
