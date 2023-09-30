package repository

import (
	"myapp/pkg/response"
	"myapp/pkg/utils"
	"myapp/src/models"

	fiber "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MyDriveRepository struct {
	DB *gorm.DB
}

func NewMyDriveRepository(Conn *gorm.DB) models.MyDriveRepository {
	return &MyDriveRepository{Conn}
}

// Delete implements models.MyDriveRepository.
func (r *MyDriveRepository) Delete(obj models.MyDrive) *fiber.Error {
	// delete permanent obj
	err := r.DB.Unscoped().Delete(&obj).Error
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return nil
}

// Get implements models.MyDriveRepository.
func (r *MyDriveRepository) Get(id string) (models.MyDrive, *fiber.Error) {
	var obj models.MyDrive
	result := r.DB.First(&obj, "id = ?", id)
	if result.RowsAffected == 0 {
		return obj, fiber.NewError(404, utils.ERR_DATA_NOT_FOUND)
	}
	return obj, nil
}

// Update implements models.MyDriveRepository.
func (r *MyDriveRepository) Update(obj models.MyDrive) (models.MyDrive, *fiber.Error) {
	err := r.DB.UpdateColumns(&obj).Error
	if err != nil {
		return obj, fiber.NewError(500, err.Error())
	}

	return obj, nil
}

// Create implements models.MyDriveRepository.
func (r *MyDriveRepository) Create(obj models.MyDrive) (models.MyDrive, *fiber.Error) {
	result := r.DB.Create(&obj)
	if result.Error != nil {
		return obj, fiber.NewError(500, result.Error.Error())
	}

	return obj, nil
}

// MyDrive implements models.MyDriveRepository.
func (r *MyDriveRepository) MyDrive(user models.User, param response.ParamsPagination) (*response.Pagination, *fiber.Error) {
	var data []*models.MyDrive
	// var count int64
	var pagination response.Pagination

	db := r.DB.Where("user_id = ?", user.ID)

	if param.Search != "" {
		// search data based on title, description
		db = db.Where("title ILIKE ?", "%"+param.Search+"%").
			Or("description ILIKE ?", "%"+param.Search+"%")
	}

	// 	fill all params pagination
	pagination.Sort = param.SortQuery
	pagination.Page = param.Page
	pagination.Limit = param.Limit

	err := db.Scopes(response.Paginate(data, &pagination, db)).Find(&data).Error
	if err != nil {
		return nil, fiber.NewError(500, db.Error.Error())
	}

	pagination.Data = data

	return &pagination, nil
}
