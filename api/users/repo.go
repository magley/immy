package users

import (
	"gorm.io/gorm"
	util "immy-api/util"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) ListUsers(offset int, limit int) ([]User, error) {
	var users []User
	result := r.DB.Limit(limit).Offset(offset).Find(&users)
	return users, result.Error
}

func (r *UserRepo) CreateUser(dto CreateUserDTO) (*User, error) {
	hashedPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := User{
		Username: dto.Username,
		Password: hashedPassword,
		Type: dto.Type,
	}
	
	result := r.DB.Create(&user)
	return &user, result.Error
}

func (r *UserRepo) GetUser(userId int) (*User, error) {
	var user User
	result := r.DB.First(&user, userId)
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(userId int, dto UpdateUserDTO) (*User, error) {
	user, err := r.GetUser(userId)
	if err != nil {
		return nil, err
	}
	
	if dto.Username != nil { user.Username = *dto.Username }
	if dto.Type != nil { user.Type = *dto.Type }
	
	result := r.DB.Save(&user)
	return user, result.Error
}

func (r *UserRepo) DeleteUser(userId int) (error) {
	user, err := r.GetUser(userId)
	if err != nil {
		return err
	}
	
	result := r.DB.Delete(&user)
	return result.Error
}

func (r *UserRepo) GetUserByName(username string) (*User, error) {
	var user User
	result := r.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}