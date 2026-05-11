package users

import (
	"time"
)

type UserType string

const (
	UserTypeAdmin = "admin"
	UserTypeModerator = "moderator"
	UserTypeJanitor = "janitor"
)

type User struct {
	ID 			uint 		`json:"id"`
	Username 	string		`json:"username"`
	Password 	string		`json:"password"`
	Type 		UserType 	`json:"type"`
	CreatedAt	time.Time	`json:"created_at"`
}

type CreateUserDTO struct {
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
	Type 		UserType   	`json:"type"`
}

type UpdateUserDTO struct {
	Username 	*string 	`json:"username"`
	Type 		*UserType   `json:"type"`
}

type LoginUserDTO struct {
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
}

type LoginResponseDTO struct {
	ID 			uint 		`json:"id"`
	Username 	string		`json:"username"`
	Type 		UserType 	`json:"type"`
	JWT 		string		`json:"jwt"`
}