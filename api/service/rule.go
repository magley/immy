package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type RuleService struct {
	RuleRepo 		*repo.RuleRepo
	BoardService 	*BoardService
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

// ===========================================================================
// RULE <--> BOARD

func (s *RuleService) CreateRuleBoard(dto model.CreateRuleBoardDTO) (*model.RuleBoard, error) {
	return s.RuleRepo.CreateRuleBoard(dto)
}

func (s *RuleService) DeleteRuleBoard(boardID uint, ruleID uint) (error) {
	rule, err := s.RuleRepo.GetRuleBoard(boardID, ruleID)
	if err != nil {
		return nil
	}

	return s.RuleRepo.DeleteRuleBoard(rule)
}

func (s *RuleService) ListAllRuleBoards() ([]*model.RuleBoard, int64, error) {
	rules, err := s.RuleRepo.ListAllRuleBoards()
	if err != nil {
		return []*model.RuleBoard{}, 0, err
	}
	totalCount := int64(len(rules))

	return rules, totalCount, err
}

func (s *RuleService) ListAllRulesOfBoard(boardID uint) ([]*model.Rule, int64, error) {
	board, err := s.BoardService.GetBoard(boardID)
	if err != nil {
		return nil, 0, err
	}
	rules, err := s.RuleRepo.ListAllRulesOfBoard(board)
	if err != nil {
		return nil, 0, err
	}
	return rules, int64(len(rules)), err
}

func (s *RuleService) ListAllBoardsOfRule(ruleID uint) ([]*model.Board, int64, error) {
	rule, err := s.GetRule(ruleID)
	if err != nil {
		return nil, 0, err
	}
	boards, err := s.RuleRepo.ListAllBoardsOfRule(rule)
	if err != nil {
		return nil, 0, err
	}
	return boards, int64(len(boards)), err
}