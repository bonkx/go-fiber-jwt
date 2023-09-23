package models

type ErrorDetailsResponse struct {
	Field   string
	Tag     string
	Message string
}

// ResponseHTTP represents response body of this API
type ResponseHTTP struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Error   *string                 `json:"error,omitempty"`
	Errors  []*ErrorDetailsResponse `json:"errors,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
