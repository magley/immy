package service

import (
	"errors"
	"fmt"
	"immy-api/model"
	"immy-api/repo"
	"immy-api/util"
)

type UserService struct {
	UserRepo 	*repo.UserRepo
}

func (s *UserService) ListUsers(offset, limit int) ([]model.User, error) {
	return s.UserRepo.ListUsers(offset, limit)
}

func (s *UserService) CreateUser(dto model.CreateUserDTO) (*model.User, error) {
	return s.UserRepo.CreateUser(dto)
}

func (s *UserService) GetUser(userId uint) (*model.User, error) {
	return s.UserRepo.GetUser(userId)
}

func (s *UserService) GetUserByName(username string) (*model.User, error) {
	return s.UserRepo.GetUserByName(username)
}

func (s *UserService) UpdateUser(userId uint, dto model.UpdateUserDTO) (*model.User, error) {
	user, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}
	
	return s.UserRepo.UpdateUser(user, dto)
}

func (s *UserService) DeleteUser(userId uint) (error) {
	user, err := s.GetUser(userId)
	if err != nil {
		return nil
	}
	
	return s.UserRepo.DeleteUser(user)
}

func (s *UserService) LoginUser(dto model.LoginUserDTO) (*model.LoginResponseDTO, error) {
	user, err := s.GetUserByName(dto.Username)
	if err != nil {
		return nil, err
	}

	ok := util.CheckPasswordHash(dto.Password, user.Password)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	
	jwt, err := util.CreateJWT(user.ID, user.Username, string(user.Type))
	if err != nil {
		return nil, err
	}
	
	return &model.LoginResponseDTO{
		ID: user.ID,
		Username: user.Username,
		Type: user.Type,
		JWT: jwt,
	}, nil
}

func (s *UserService) AuthorizeUser(dto model.AuthorizationDTO, jwt *util.JWTClaims) error {
	if dto.Role != nil {
		if jwt.Role != *dto.Role {
			return errors.New("Invalid role")
		}
	}

	user, err := s.GetUser(jwt.Id)
	if err != nil {
		return err
	}

	if jwt.Role != string(user.Type) {
		return errors.New("Invalid role")
	}

	return nil
}