package models

type ErrorDetailsResponse struct {
	Message string           `json:"message"`
	Errors  []*ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
