package models

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/response"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileType string

const (
	ImageFile FileType = "I"
	FileFile  FileType = "F"
	VideoFile FileType = "V"
)

type MyDrive struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;unique;type:uuid;column:id;default:uuid_generate_v4()"`
	Name     string    `json:"name"`
	FileType FileType  `gorm:"column:file_type;size:1;" json:"file_type"`
	Link     string    `json:"link"`
	// foreignkey User
	UserID uint
	// User   User `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (md *MyDrive) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (md MyDrive) MarshalJSON() ([]byte, error) {
	type Alias MyDrive
	var link string = md.Link
	if !strings.HasPrefix(link, "http") {
		link = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), md.Link)
	}

	aux := struct {
		Alias
		Link string `json:"link"`
	}{
		Alias: (Alias)(md),
		Link:  link,
	}
	return json.Marshal(aux)
}

type MyDriveUsecase interface {
	// USECASE
	MyDrive(c *fiber.Ctx) (*response.Pagination, *fiber.Error)
	// Get(c *fiber.Ctx) (MyDrive, *fiber.Error)
	Create(c *fiber.Ctx) ([]*MyDrive, *fiber.Error)
	// Update(c *fiber.Ctx, id uint) (MyDrive, *fiber.Error)
	// Delete(c *fiber.Ctx, id uint) *fiber.Error

}

type MyDriveRepository interface {
	// FUNTIONS

	// REPOS
	MyDrive(user User, param response.ParamsPagination) (*response.Pagination, *fiber.Error)
	Create(md MyDrive) (MyDrive, *fiber.Error)
}
