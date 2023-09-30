package usecase

import (
	"fmt"
	"myapp/pkg/response"
	"myapp/pkg/utils"
	"myapp/src/models"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
)

type MyDriveUsecase struct {
	dRepo models.MyDriveRepository
	uRepo models.UserRepository
}

func NewMyDriveUsecase(md models.MyDriveRepository, user models.UserRepository) models.MyDriveUsecase {
	return &MyDriveUsecase{
		dRepo: md,
		uRepo: user,
	}
}

// Create implements models.MyDriveUsecase.
func (uc *MyDriveUsecase) Create(c *fiber.Ctx) ([]*models.MyDrive, *fiber.Error) {
	var listObj []*models.MyDrive

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return nil, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// MultipartForm POST
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["files"]

		if files == nil {
			return nil, fiber.NewError(400, "Please select the files to upload")
		}

		// Loop through files:
		for _, file := range files {
			// declare new model
			var obj models.MyDrive
			// fill the owner
			obj.UserID = user.ID

			// get mime type
			_fileType := string(models.FileFile)
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			if strings.HasPrefix(file.Header["Content-Type"][0], "image/") {
				_fileType = string(models.ImageFile)
			}
			if strings.HasPrefix(file.Header["Content-Type"][0], "video/") {
				_fileType = string(models.VideoFile)
			}
			obj.Name = file.Filename
			obj.FileType = models.FileType(_fileType)

			// Save the files to disk:
			imageUrl, errFile := utils.FileUpload(c, file, "drive")
			if errFile != nil {
				return nil, fiber.NewError(500, errFile.Error())
			}
			// update file url
			obj.Link = imageUrl

			// save the data
			obj, err := uc.dRepo.Create(obj)
			if err != nil {
				return nil, err
			}

			// append data to list data
			listObj = append(listObj, &obj)
		}
	}

	return listObj, nil
}

// MyDrive implements models.MyDriveUsecase.
func (uc *MyDriveUsecase) MyDrive(c *fiber.Ctx) (*response.Pagination, *fiber.Error) {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return nil, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// 	Parse the query parameters
	search := c.Query("search")
	sortBy := c.Query("sort", "created_at|desc")
	page := c.Query("page", "1")
	limit := c.Query("per_page", "10")

	// Convert the page and limit to integers
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	sortQuery, errSort := utils.ValidateAndReturnSortQuery(sortBy)
	// log.Print(sortQuery)
	if errSort != nil {
		errD := fiber.NewError(fiber.StatusInternalServerError, errSort.Error())
		return nil, errD
	}

	// make param pagination struct
	pagParam := response.ParamsPagination{
		Page:      pageInt,
		Limit:     limitInt,
		SortQuery: sortQuery,
		Search:    search,
		NoPage:    c.Query("no_page"),
	}

	pagination, err := uc.dRepo.MyDrive(user, pagParam)
	if err != nil {
		return nil, err
	}
	return pagination, nil
}
