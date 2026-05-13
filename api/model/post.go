package model

import (
	"time"
)

type Post struct {
	ID 			uint 		`json:"id"`
	ThreadID	uint		`json:"thread_id"`
	Num			uint 		`json:"num"`
	Name 		string 		`json:"name"`
	Tripcode 	string 		`json:"tripcode"`
	IPv4 		string		`json:"ipv4"`
	CreatedAt	time.Time	`json:"created_at"`
	Sage 		bool		`json:"sage"`
	Content 	string		`json:"content"`
	Filename 	string 		`json:"filename"`
	Html		string 		`json:"html"`
}

type CreatePostForThreadDTO struct {
	Name 		string 		`json:"name" binding:"required"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	string 		`json:"filename" binding:"required"`
	Options 	*string		`json:"options"`
}

type CreatePostDTO struct {
	Name 		string 		`json:"name" binding:"required"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	string 		`json:"filename" binding:"required"`
	Options 	*string		`json:"options"`
	
	ThreadID	uint 		`json:"thread_id"`
}

type UpdatePostDTO struct {
	Name 		*string 	`json:"name"`
	Code  		*string 	`json:"code"`
	Description *string 	`json:"description"`
	Locked		*bool		`json:"locked"`
	Hidden		*bool		`json:"hidden"`
}