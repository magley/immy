package model

import (
	"time"
)

type Post struct {
	ID 			uint 		`json:"id"`
	ThreadID	uint		`json:"thread_id"`
	BoardID		uint 		`json:"board_id"`
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
	Name 		string 		`json:"name"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	string 		`json:"filename" binding:"required"`
	Options 	string		`json:"options"`
}

type CreatePostDTO struct {
	Name 		string 		`json:"name"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	string 		`json:"filename"`
	Options 	string		`json:"options"`
	
	ThreadID	uint 		`json:"thread_id"`
}

type UpdatePostDTO struct {
	Name 		*string 	`json:"name"`
	Tripcode 	*string 	`json:"tripcode"`
	Sage 		*bool		`json:"sage"`
	Content 	*string		`json:"content"`
	Filename 	*string 	`json:"filename"`
	Html		*string 	`json:"html"`
}