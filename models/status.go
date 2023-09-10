package models

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name string `json:"name" validate:"required" gorm:"size:50;not null;"`
}

// TableName overrides the table name used by Statuses to `status`
func (Status) TableName() string {
	return "status"
}
