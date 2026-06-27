package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type RuleService struct {
	RuleRepo 		*repo.RuleRepo
}

func (s *RuleService) ListRules(offset, limit int) ([]*model.Rule, int64, error) {
	Rules, err := s.RuleRepo.ListRules(offset, limit)
	if err != nil {
		return []*model.Rule{}, 0, err
	}
	totalCount, err := s.RuleRepo.GetRuleCount(true)
	if err != nil {
		return []*model.Rule{}, 0, err
	}

	return Rules, totalCount, err
}

func (s *RuleService) CreateRule(dto model.CreateRuleDTO) (*model.Rule, error) {
	return s.RuleRepo.CreateRule(dto)
}

func (s *RuleService) GetRule(RuleId uint) (*model.Rule, error) {
	return s.RuleRepo.GetRule(RuleId)
}

func (s *RuleService) UpdateRule(RuleId uint, dto model.UpdateRuleDTO) (*model.Rule, error) {
	Rule, err := s.GetRule(RuleId)
	if err != nil {
		return nil, err
	}

	return s.RuleRepo.UpdateRule(Rule, dto)
}

func (s *RuleService) DeleteRule(RuleId uint) (error) {
	Rule, err := s.GetRule(RuleId)
	if err != nil {
		return nil
	}

	return s.RuleRepo.DeleteRule(Rule)
}