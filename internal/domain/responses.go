package domain

type AddSongResponse struct {
	ID int `json:"id"`
}

type GetSongVersesResponse struct {
	Verses []string `json:"verses"`
	Page   int      `json:"page"`
	Size   int      `json:"size"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(code int, message, details string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e ErrorResponse) Error() string {
	return e.Message
}
