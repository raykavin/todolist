package dto

// Response is the base response structure
type Response struct {
	Success bool       `json:"success"                 example:"true"`
	Message string     `json:"message,omitempty"       example:"Operation successful"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string         `json:"code"              example:"VALIDATION_ERROR"`
	Message string         `json:"message"           example:"Invalid input data"`
	Details map[string]any `json:"details,omitempty"`
}

// PaginatedResponse is the paginated response structure
type PaginatedResponse struct {
	Success    bool       `json:"success"    example:"true"`
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"        example:"1"`
	PageSize   int   `json:"page_size"   example:"20"`
	Total      int64 `json:"total"       example:"100"`
	TotalPages int   `json:"total_pages" example:"5"`
}

// SuccessResponse creates a success response
func SuccessResponse(data any, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code, message string, details map[string]any) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// PaginatedSuccessResponse creates a paginated success response
func PaginatedSuccessResponse(data any, page, pageSize int, total int64) PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// QueryParams common query parameters
type QueryParams struct {
	Page     int    `form:"page"      binding:"min=0"`
	PageSize int    `form:"page_size" binding:"min=0,max=100"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`
	OrderBy  string `form:"order_by"`
	Search   string `form:"search"`
}

// SetDefaults sets default values for query params
func (q *QueryParams) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	if q.OrderDir == "" {
		q.OrderDir = "asc"
	}
}

// GetOffset calculates the offset for pagination
func (q *QueryParams) GetOffset() int {
	return (q.Page - 1) * q.PageSize
}
