package model

import (
	"time"
)

type Board struct {
	ID 			uint 		`json:"id"`
	Name 		string 		`json:"name"`
	Code  		string 		`json:"code"`
	Description *string 	`json:"description"`
	CreatedAt	time.Time	`json:"created_at"`
	Locked		bool		`json:"locked"`
	Hidden		bool		`json:"hidden"`
	PostCount	uint 		`json:"post_count"`	
}

type CreateBoardDTO struct {
	Name 		string 		`json:"name" binding:"required"`
	Code  		string 		`json:"code" binding:"required"`
	Description *string 	`json:"description"`
}

type UpdateBoardDTO struct {
	Name 		*string 	`json:"name"`
	Code  		*string 	`json:"code"`
	Description *string 	`json:"description"`
	Locked		*bool		`json:"locked"`
	Hidden		*bool		`json:"hidden"`
}