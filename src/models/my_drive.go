package models

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/response"
	"myapp/pkg/utils"
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

type MyDriveRenameInput struct {
	Name string `json:"name" form:"name" validate:"required,min=4"`
}

func (md *MyDrive) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (md MyDrive) MarshalJSON() ([]byte, error) {
	type Alias MyDrive
	var link string = md.Link
	if !strings.HasPrefix(link, "http") {
		link = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), md.Link)
	}

	// thumbnail allow null
	var thumbnail *string = nil

	// thumbnail for image file
	if md.FileType == ImageFile {
		thumbnail = utils.GetThumbnail(md.Link)
		if thumbnail != nil {
			*thumbnail = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), *thumbnail)
		}
	}

	// thumbnail for video file
	if md.FileType == VideoFile {
		thumbnail = utils.GetThumbnailVideo(md.Link)
		if thumbnail != nil {
			*thumbnail = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), *thumbnail)
		}
	}

	aux := struct {
		Alias
		Link      string  `json:"link"`
		Thumbnail *string `json:"thumbnail"`
	}{
		Alias:     (Alias)(md),
		Link:      link,
		Thumbnail: thumbnail,
	}
	return json.Marshal(aux)
}

type MyDriveUsecase interface {
	// USECASE
	MyDrive(c *fiber.Ctx) (*response.Pagination, *fiber.Error)
	Get(c *fiber.Ctx) (MyDrive, *fiber.Error)
	Create(c *fiber.Ctx) ([]*MyDrive, *fiber.Error)
	Update(c *fiber.Ctx, payload MyDriveRenameInput) (MyDrive, *fiber.Error)
	Delete(c *fiber.Ctx) *fiber.Error
}

type MyDriveRepository interface {
	// FUNTIONS

	// REPOS
	MyDrive(user User, param response.ParamsPagination) (*response.Pagination, *fiber.Error)
	Get(id string) (MyDrive, *fiber.Error)
	Create(obj MyDrive) (MyDrive, *fiber.Error)
	Update(obj MyDrive) (MyDrive, *fiber.Error)
	Delete(obj MyDrive) *fiber.Error
}
