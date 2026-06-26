package model

import (
	"time"

	"gorm.io/gorm"
)


type Ban struct {
	ID 				uint 			`json:"id"`
	// IP Addresses are saved as numbers in the DB for range queries.
	IpStart 		uint64 			`json:"ip_start"`
	IpEnd 			*uint64 		`json:"ip_end"`
	CreatedAt 		time.Time 		`json:"created_at"`
	ExpiresAt		*time.Time 		`json:"expires_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"deleted_at"`
	BoardID 		*uint 			`json:"board_id"`
	BoardCode 		string 			`json:"board_code"`			// Redundant
	CreatorID 		uint 			`json:"creator_id"`
	CreatorUsername string 			`json:"creator_username"`	// Redundant
	Reason 			string 			`json:"reason"`
	Warning 		bool 			`json:"warning"`
	Seen			bool			`json:"seen"`
}

type CreateBanDTO struct {
	IpStart 		string 			`json:"ip_start"`
	IpEnd 			*string 		`json:"ip_end"`
	ExpiresAt 		*time.Time 		`json:"expires_at"`
	BoardID 		*uint 			`json:"board_id"`
	Reason 			string 			`json:"reason"`
	Warning 		bool 			`json:"warning"`
}

type UpdateBanDTO struct {
	Seen 			*bool			`json:"seen"`
}