package models

type ErrorDetailsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Errors  []*ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Field   string
	Tag     string
	Message string
}

type ValueError struct {
	Code int
	Err  error
}
