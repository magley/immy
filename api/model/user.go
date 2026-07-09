package model

import (
	"time"
)

type UserRole string

const (
	UserRoleAdmin = "admin"
	UserRoleModerator = "moderator"
	UserRoleJanitor = "janitor"
)

type User struct {
	ID 			uint 		`json:"id"`
	Username 	string		`json:"username"`
	Password 	string		`json:"password"`
	Role 		UserRole 	`json:"role"`
	CreatedAt	time.Time	`json:"created_at"`
}

type CreateUserDTO struct {
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
	Role 		UserRole   	`json:"role"`
}

type CreateFirstAdminDTO struct {
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
}

type UpdateUserDTO struct {
	Username 	*string 	`json:"username"`
	Role 		*UserRole   `json:"role"`
}

type LoginUserDTO struct {
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
}

type LoginResponseDTO struct {
	ID 			uint 		`json:"id"`
	Username 	string		`json:"username"`
	Role 		UserRole 	`json:"role"`
	JWT 		string		`json:"jwt"`
}

type AuthorizationDTO struct {
	RequiredRoles []string	`json:"required_roles"`
}