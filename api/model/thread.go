package model

import (
)

type Thread struct {
	ID 			uint 		`json:"id"`
	BoardID		uint		`json:"board_id"`
	Subject 	string		`json:"subject"`
	Locked		bool		`json:"locked"`
	Sticky		bool		`json:"sticky"`
}

type CreateThreadDTO struct {
	BoardCode	string		`json:"board_code" binding:"required"`
	Subject 	string		`json:"subject"`
	Locked		bool		`json:"locked"`
	Sticky		bool		`json:"sticky"`
	Post 		CreatePostForThreadDTO `json:"post"`
}

type UpdateThreadDTO struct {
	Locked		*bool		`json:"locked"`
	Sticky		*bool		`json:"sticky"`
}