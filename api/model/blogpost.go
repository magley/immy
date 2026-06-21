package model

import (
	"time"

	"gorm.io/gorm"
)


type Blogpost struct {
	ID 				uint 			`json:"id"`
	Title 			string			`json:"title"`
	Html 			string			`json:"html"`
	AuthorID		uint	 		`json:"author_id"`
	AuthorName 		string			`json:"author_name"` // Redundant
	CreatedAt 		time.Time 		`json:"created_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"deleted_at"`
}

type BlogpostShortDTO struct {
	ID 				uint 			`json:"id"`
	Title 			string			`json:"title"`
	CreatedAt 		time.Time 		`json:"created_at"`
}

type CreateBlogpostDTO struct {
	Title 			string			`json:"title"`
	Html 			string			`json:"html"`
}

type UpdateBlogpostDTO struct {
	Html 			*string			`json:"html"`
}