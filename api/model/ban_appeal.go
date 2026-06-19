package model

import (
	"time"

	"gorm.io/gorm"
)


type BanAppealStatus string

const (
	BanAppealStatusPending = "pending"
	BanAppealStatusRejected = "rejected"
	BanAppealStatusRejectedFinal = "rejected_final"
	BanAppealStatusApproved = "approved"
)

type BanAppeal struct {
	ID 				uint 			`json:"id"`
	BanID			uint 			`json:"ban_id"`
	Message 		string 			`json:"message"`
	Status			BanAppealStatus	`json:"status"`
	CreatedAt 		time.Time 		`json:"created_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"deleted_at"`
	ReviewedAt		*time.Time 		`json:"reviewed_at"`
	ReviewedBy		*uint			`json:"reviewed_by"`
}

type CreateBanAppealDTO struct {
	BanID			uint 			`json:"ban_id"`
	Message 		string 			`json:"message"`
}

type UpdateBanAppealDTO struct {
	Status			*BanAppealStatus `json:"status"`
}