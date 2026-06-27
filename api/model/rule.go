package model

import (
	"time"

	"gorm.io/gorm"
)


type Rule struct {
	ID 				uint 			`json:"id"`
	Title	 		string 			`json:"title"`
	Description 	string 			`json:"description"`
	IsGlobal 		bool			`json:"is_global"`
	Danger 			int 			`json:"danger"`
	CreatedAt		time.Time 		`json:"created_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"deleted_at"`
}

type CreateRuleDTO struct {
	Title	 		string 			`json:"title"`
	Description 	string 			`json:"description"`
	IsGlobal 		bool			`json:"is_global"`
	Danger 			int 			`json:"danger"`
}

type UpdateRuleDTO struct {
	Title	 		*string 		`json:"title"`
	Description 	*string 		`json:"description"`
	IsGlobal 		*bool			`json:"is_global"`
	Danger 			*int 			`json:"danger"`
}