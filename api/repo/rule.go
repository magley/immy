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