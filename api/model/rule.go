package model

import (
	"time"

	"gorm.io/gorm"
)


type Rule struct {
	ID 				uint 			`json:"id"`
	Title	 		string 			`json:"title"`
	Description 	string 			`json:"description"`
	IsGlobal 		bool			`json:"is_global"`
	Danger 			int 			`json:"danger"`
	CreatedAt		time.Time 		`json:"created_at"`
	DeletedAt 		gorm.DeletedAt 	`json:"deleted_at"`
}

type RuleBoard struct {
	RuleID			uint			`json:"rule_id"`
	BoardID			uint			`json:"board_id"`
	CreatedAt		time.Time		`json:"created_at"`
	DeletedAt		gorm.DeletedAt	`json:"deleted_at"`
}

func (RuleBoard) TableName() string {
    return "rules_boards"
}

type CreateRuleDTO struct {
	Title	 		string 			`json:"title"`
	Description 	string 			`json:"description"`
	IsGlobal 		bool			`json:"is_global"`
	Danger 			int 			`json:"danger"`
}

type CreateRuleBoardDTO struct {
	RuleID			uint			`json:"rule_id"`
	BoardID			uint			`json:"board_id"`
}

type RuleBoardIdDTO struct {
	RuleID			uint			`json:"rule_id"`
	BoardID			uint			`json:"board_id"`
}

type UpdateRuleDTO struct {
	Title	 		*string 		`json:"title"`
	Description 	*string 		`json:"description"`
	IsGlobal 		*bool			`json:"is_global"`
	Danger 			*int 			`json:"danger"`
}