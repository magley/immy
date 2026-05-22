package model

import "time"

type Thread struct {
	ID 			uint 		`json:"id"`
	BoardID		uint		`json:"board_id"`
	PostNum		uint  		`json:"post_num"`
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

type ThreadFullDTO struct {
	Thread 		*Thread  	`json:"thread"`
	Posts		[]Post 		`json:"posts"`
}

type ThreadForCatalogDTO struct {
	Thread		Thread	`json:"thread"`
	Post		Post	`json:"post"`
	LastPost 	Post 	`json:"last_post"`
	Stats 		ThreadStats	`json:"stats"`
}

type ThreadStats struct {
	PostCount 	uint 		`json:"post_count"`
	ImageCount 	uint 		`json:"image_count"`
	UserCount 	uint 		`json:"user_count"`
	LastBump 	time.Time 	`json:"last_bump"`
}