package responses

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
