package utils

import "gorm.io/gorm"

func OwnerThis(user_id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", user_id)
	}
}
