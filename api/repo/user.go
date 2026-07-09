package repo

import (
	util "immy-api/util"

	"gorm.io/gorm"

	model "immy-api/model"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) ListUsers(offset int, limit int) ([]model.User, error) {
	var users []model.User
	result := r.DB.Limit(limit).Offset(offset).Find(&users)
	return users, result.Error
}

func (r *UserRepo) GetUserCountOfRole(role model.UserRole) (int64, error) {
	var count int64
	result := r.DB.Model(&model.User{}).Where("role = ?", role).Count(&count)
	return count, result.Error
}

func (r *UserRepo) CreateUser(dto model.CreateUserDTO) (*model.User, error) {
	hashedPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: dto.Username,
		Password: hashedPassword,
		Role:     dto.Role,
	}

	result := r.DB.Create(&user)
	return &user, result.Error
}

func (r *UserRepo) GetUser(userId uint) (*model.User, error) {
	var user model.User
	result := r.DB.First(&user, userId)
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user *model.User, dto model.UpdateUserDTO) (*model.User, error) {
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Role != nil {
		user.Role = *dto.Role
	}

	result := r.DB.Save(&user)
	return user, result.Error
}

func (r *UserRepo) DeleteUser(user *model.User) error {
	result := r.DB.Delete(&user)
	return result.Error
}

func (r *UserRepo) GetUserByName(username string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}
