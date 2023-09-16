package models

type ErrorDetailsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Errors  []*ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}