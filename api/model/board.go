package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Board struct {
	ID 			uint 			`json:"id"`
	Name 		string 			`json:"name"`
	Code  		string 			`json:"code"`
	Description *string 		`json:"description"`
	CreatedAt	time.Time		`json:"created_at"`
	DeletedAt	gorm.DeletedAt `json:"deleted_at"`
	Config		BoardConfig		`json:"config" gorm:"embedded"`
	Meta		BoardMeta		`json:"meta" gorm:"embedded"`
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
	MaxThreads  uint			`json:"max_threads"`
	AllowSpoilers bool			`json:"allow_spoilers"`
	SpoilerImage string			`json:"spoiler_image"`
}

type BoardMeta struct {
	PostCount		uint 		`json:"post_count"`
	BytesUploaded	uint 		`json:"bytes_uploaded"`
}

type CreateBoardDTO struct {
	Name 		string 		`json:"name" binding:"required"`
	Code  		string 		`json:"code" binding:"required"`
	Description *string 	`json:"description"`
	Config		BoardConfig		`json:"config" gorm:"embedded"`
}

type BoardStatisticsDTO struct {
	ID 				uint 	`json:"id"`
	Code 			string 	`json:"code"`
	ThreadCount 	uint 	`json:"thread_count"`
	PostCount 		uint 	`json:"post_count"`
	BytesUploaded 	uint	`json:"bytes_uploaded"`
}
