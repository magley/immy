package model

import (
	"time"

	"github.com/lib/pq"
)

type Board struct {
	ID 			uint 			`json:"id"`
	Name 		string 			`json:"name"`
	Code  		string 			`json:"code"`
	Description *string 		`json:"description"`
	CreatedAt	time.Time		`json:"created_at"`
	PostCount	uint 			`json:"post_count"`
	Config		BoardConfig		`json:"config" gorm:"embedded"`
}

type BoardConfig struct {
	Locked		bool			`json:"locked"`
	Hidden		bool			`json:"hidden"`
	MaxFileSize	uint			`json:"max_file_size"`
	ReplyFilesAllowed bool		`json:"reply_files_allowed"`
	MimeTypesAllowed pq.StringArray 	`json:"mime_types_allowed" gorm:"type:text[]"`
	BumpLimit	uint			`json:"bump_limit"`
	ImageLimit	uint			`json:"image_limit"`
	FlagsEnabled bool			`json:"flags_enabled"`
	IDsEnabled	bool			`json:"ids_enabled"`
	CodeEnabled bool 			`json:"code_enabled"`
	MathEnabled bool			`json:"math_enabled"`
}

type CreateBoardDTO struct {
	Name 		string 		`json:"name" binding:"required"`
	Code  		string 		`json:"code" binding:"required"`
	Description *string 	`json:"description"`
	Config		BoardConfig		`json:"config" gorm:"embedded"`
}

type UpdateBoardDTO struct {
	Name 		*string 	`json:"name"`
	Code  		*string 	`json:"code"`
	Description *string 	`json:"description"`
	PostCount	*uint 		/* No json */
}