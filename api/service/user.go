package service

import (
	"immy-api/repo"
	_"immy-api/model"
)

type UserService struct {
	UserRepo 	*repo.UserRepo
}
