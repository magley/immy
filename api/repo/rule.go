package repo

import (
	"gorm.io/gorm"

	model "immy-api/model"
)

type RuleRepo struct {
	DB *gorm.DB
}

func (r *RuleRepo) ListRules(offset int, limit int) ([]*model.Rule, error) {
	var Rules []*model.Rule
	result := r.DB.Limit(limit).Offset(offset).Find(&Rules)
	return Rules, result.Error
}

func (r *RuleRepo) GetRuleCount(includeDeleted bool) (int64, error) {
	cnt := int64(0)

	query := r.DB
	if (includeDeleted) {
		query = query.Unscoped()
	}
	result := query.Model(&model.Rule{}).Count(&cnt)
	return cnt, result.Error
}

func (r *RuleRepo) CreateRule(dto model.CreateRuleDTO) (*model.Rule, error) {
	Rule := model.Rule{
		Title: dto.Title,
		Description: dto.Description,
		IsGlobal: dto.IsGlobal,
		Danger: dto.Danger,
	}

	result := r.DB.Create(&Rule)
	return &Rule, result.Error
}

func (r *RuleRepo) GetRule(RuleId uint) (*model.Rule, error) {
	var Rule model.Rule
	result := r.DB.First(&Rule, RuleId)
	return &Rule, result.Error
}

func (r *RuleRepo) UpdateRule(Rule *model.Rule, dto model.UpdateRuleDTO) (*model.Rule, error) {
	if dto.Title != nil { Rule.Title = *dto.Title }
	if dto.Description != nil { Rule.Description = *dto.Description }
	if dto.IsGlobal != nil { Rule.IsGlobal = *dto.IsGlobal }
	if dto.Danger != nil { Rule.Danger = *dto.Danger }
	result := r.DB.Save(&Rule)
	return Rule, result.Error
}

func (r *RuleRepo) DeleteRule(Rule *model.Rule) (error) {
	result := r.DB.Delete(&Rule)
	return result.Error
}

// ===========================================================================
// RULE <--> BOARD

// -------------- Low level ------------------------------------------

func (r *RuleRepo) ListAllRuleBoards() ([]*model.RuleBoard, error) {
	var Rules []*model.RuleBoard
	result := r.DB.Find(&Rules)
	return Rules, result.Error
}

func (r *RuleRepo) ListAllRuleBoardsOfBoard(boardID uint) ([]*model.RuleBoard, error) {
	var Rules []*model.RuleBoard
	result := r.DB.Where("board_id = ?", boardID).Find(&Rules)
	return Rules, result.Error
}

func (r *RuleRepo) ListAllRuleBoardsOfRule(ruleID uint) ([]*model.RuleBoard, error) {
	var Rules []*model.RuleBoard
	result := r.DB.Where("rule_id = ?", ruleID).Find(&Rules)
	return Rules, result.Error
}

func (r *RuleRepo) CreateRuleBoard(dto model.CreateRuleBoardDTO) (*model.RuleBoard, error) {
	Rule := model.RuleBoard{
		RuleID: dto.RuleID,
		BoardID: dto.BoardID,
	}

	result := r.DB.Create(&Rule)
	return &Rule, result.Error
}

func (r *RuleRepo) GetRuleBoard(boardID uint, ruleID uint) (*model.RuleBoard, error) {
	var rule model.RuleBoard
	result := r.DB.Where("board_id = ?", boardID).Where("rule_id = ?", ruleID).First(&rule)
	return &rule, result.Error
}

func (r *RuleRepo) DeleteRuleBoard(ruleBoard *model.RuleBoard) (error) {
	result := r.DB.Delete(&ruleBoard)
	return result.Error
}

// -------------- High level ------------------------------------------

func (r *RuleRepo) ListAllRulesOfBoard(board *model.Board) ([]*model.Rule, error) {
	var rules []*model.Rule
	result := r.DB.Model(&model.Rule{}).
		Joins("rules_boards ON rules_boards.rule_id = rules.id").
		Where("rules_boards.board_id = ? or rules.is_global = TRUE", board.ID).
		Find(&rules)
	return rules, result.Error
}

func (r *RuleRepo) ListAllBoardsOfRule(rule *model.Rule) ([]*model.Board, error) {
	var boards []*model.Board
	var result *gorm.DB

	if rule.IsGlobal {
		result = r.DB.Model(&model.Board{}).Find(&boards)
	} else {
		result = r.DB.Model(&model.Board{}).
			Joins("rules_boards ON rules_boards.board_id = boards.id").
			Where("rules_boards.rule_id = ?", rule.ID).
			Find(&boards)
	}

	return boards, result.Error
}