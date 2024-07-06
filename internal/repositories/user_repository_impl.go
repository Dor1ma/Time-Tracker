package repositories

import (
	"fmt"
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepositoryImpl(db *gorm.DB, logger *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	r.logger.Infof("Create: creating user in database with PassportNumber %s", user.PassportNumber)
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Errorf("Create: failed to create user in database: %v", err)
		return err
	}

	r.logger.Infof("Create: user created in database successfully with ID %d", user.ID)
	return nil
}

func (r *UserRepositoryImpl) GetById(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		r.logger.Errorf("GetById: failed to get user from database with ID %d: %v", id, result.Error)
		return nil, result.Error
	}

	r.logger.Infof("GetById: successfully retrieved user from database with ID %d", id)
	return &user, nil
}

func (r *UserRepositoryImpl) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		r.logger.Errorf("GetAll: failed to fetch all users from database: %v", result.Error)
		return nil, result.Error
	}

	r.logger.Infof("GetAll: successfully fetched %d users from database", len(users))
	return users, nil
}

func (r *UserRepositoryImpl) GetAllWithFiltersAndPagination(filters map[string]interface{}, page int, pageSize int) ([]models.User, error) {
	var users []models.User
	query := r.db.Model(&models.User{})

	r.logger.Debugf("GetAllWithFiltersAndPagination: original filters: %v", filters)
	delete(filters, "page")
	delete(filters, "pageSize")

	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	r.logger.Debugf("GetAllWithFiltersAndPagination: query with filters: %v", query)
	result := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	if result.Error != nil {
		r.logger.Errorf("GetAllWithFiltersAndPagination: failed to fetch users with filters and pagination from database: %v", result.Error)
		return nil, result.Error
	}

	r.logger.Infof("GetAllWithFiltersAndPagination: successfully fetched %d users with filters and pagination from database", len(users))
	return users, nil
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
	r.logger.Infof("Update: updating user in database with ID %d", user.ID)
	if err := r.db.Save(user).Error; err != nil {
		r.logger.Errorf("Update: failed to update user in database with ID %d: %v", user.ID, err)
		return err
	}

	r.logger.Infof("Update: user with ID %d updated successfully in database", user.ID)
	return nil
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	r.logger.Infof("Delete: deleting user from database with ID %d", id)
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		r.logger.Errorf("Delete: failed to delete user from database with ID %d: %v", id, err)
		return err
	}

	r.logger.Infof("Delete: user with ID %d deleted from database successfully", id)
	return nil
}
