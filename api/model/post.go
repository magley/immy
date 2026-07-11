package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID 			uint 		`json:"id"`
	ThreadID	uint		`json:"thread_id"`
	ThreadNum	uint		`json:"thread_num"` // Redundant
	BoardID		uint 		`json:"board_id"`   // Redundant
	UserID		*uint	     //`json:"user_id"`
	UserRole	*UserRole	`json:"user_role"`  // Redundant
	Num			uint 		`json:"num"`
	Name 		string 		`json:"name"`
	PublicID	*string		`json:"public_id"`
	Tripcode 	string 		`json:"tripcode"`
	IPv4 		string		`json:"ipv4"`
	CreatedAt	time.Time	`json:"created_at"`
	DeletedAt	gorm.DeletedAt `json:"deleted_at"`
	Sage 		bool		`json:"sage"`
	Capcode 	bool		`json:"capcode"`
	Content 	string		`json:"content"`
	Filename 	string 		`json:"filename"`
	Filesize 	uint 		`json:"filesize"`
	ImgWidth	uint		`json:"img_width"`
	ImgHeight	uint		`json:"img_height"`
	MD5			string		`json:"md5"`        // Base64 of file
	SrcFilename string 		`json:"src_filename"`
	Spoiler 	bool 		`json:"spoiler"`
	Html		string 		`json:"html"`
}

type CreatePostForThreadDTO struct {
	Name 		string 		`json:"name"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	string 		`json:"filename" binding:"required"`
	Filebytes 	string    	`json:"filebytes" binding:"required"`
	Options 	string		`json:"options"`
	Spoiler		bool		`json:"spoiler"`
}

type CreatePostDTO struct {
	Name 		string 		`json:"name"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	*string 	`json:"filename"`
	Filebytes 	*string    	`json:"filebytes"`
	Options 	string		`json:"options"`
	Spoiler		bool		`json:"spoiler"`
	
	ThreadID	uint 		`json:"thread_id"`
}

type UpdatePostDTO struct {
	Name 		*string 	`json:"name"`
	Tripcode 	*string 	`json:"tripcode"`
	Sage 		*bool		`json:"sage"`
	Content 	*string		`json:"content"`
	Filename 	*string 	`json:"filename"`
	Html		*string 	`json:"html"`
	Spoiler		*bool		`json:"spoiler"`
}

// Internal use only
type CreatePostCommonDTO struct {
	Name 		string 		`json:"name"`
	Content 	string		`json:"content" binding:"required"`
	Filename 	*string 	`json:"filename"`
	Filebytes 	*string    	`json:"filebytes"`
	Options 	string		`json:"options"`
	Spoiler		bool		`json:"spoiler"`
	ThreadID	*uint 		`json:"thread_id"`
}