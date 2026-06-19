package service

import (
	"errors"
	"immy-api/model"
	"immy-api/repo"
)

type BanAppealService struct {
	BanAppealRepo 	*repo.BanAppealRepo
	BanService 		*BanService
	UserService 	*UserService
}

func (s *BanAppealService) ListBanAppeals(offset, limit int) ([]*model.BanAppeal, int64, error) {
	banappeals, err := s.BanAppealRepo.ListBanAppeals(offset, limit)
	if err != nil {
		return banappeals, 0, err
	}
	totalCount, err := s.BanAppealRepo.GetBanAppealCount(true)
	if err != nil {
		return banappeals, 0, err
	}
	return banappeals, totalCount, err
}

func (s *BanAppealService) CreateBanAppeal(dto model.CreateBanAppealDTO) (*model.BanAppeal, error) {
	ban, err := s.BanService.GetBan(dto.BanID)
	if err != nil {
		return nil, err
	}

	canAppeal, err := s.CanAppealBan(dto.BanID)
	if err != nil {
		return nil, err
	}

	if canAppeal == false {
		return nil, errors.New("You cannot appeal this ban anymore")
	}

	return s.BanAppealRepo.CreateBanAppeal(dto, ban)
}

func (s *BanAppealService) GetBanAppeal(banappealId uint) (*model.BanAppeal, error) {
	return s.BanAppealRepo.GetBanAppeal(banappealId)
}

func (s *BanAppealService) UpdateBanAppeal(banappealId uint, dto model.UpdateBanAppealDTO, user *model.User) (*model.BanAppeal, error) {
	banappeal, err := s.GetBanAppeal(banappealId)
	if err != nil {
		return nil, err
	}

	return s.BanAppealRepo.UpdateBanAppeal(banappeal, dto, user)
}

func (s *BanAppealService) DeleteBanAppeal(banappealId uint) (error) {
	banappeal, err := s.GetBanAppeal(banappealId)
	if err != nil {
		return nil
	}

	return s.BanAppealRepo.DeleteBanAppeal(banappeal)
}

func (s *BanAppealService) GetBanAppealsOfBan(banID uint) ([]*model.BanAppeal, error) {
	ban, err := s.BanService.GetBan(banID)
	if err != nil {
		return nil, err
	}
	return s.BanAppealRepo.GetBanAppealsOfBan(ban)
}

func (s *BanAppealService) CanAppealBan(banID uint) (bool, error) {
	ban, err := s.BanService.GetBan(banID)
	if err != nil {
		return false, err
	}
	return s.BanAppealRepo.CanAppealBan(ban)
}

// =====================================================================================

func (s *BanAppealService) censorBanAppeal(banappeal *model.BanAppeal) (*model.BanAppeal) {
	banappeal.ReviewedBy = nil
	return banappeal
}

func (s *BanAppealService) censorBanAppeals(banappeals []*model.BanAppeal) ([]*model.BanAppeal) {
	for _, banappeal := range banappeals {
		banappeal.ReviewedBy = nil
	}
	return banappeals
}

func (s *BanAppealService) censorBanAppeals2(banappeals []*model.BanAppeal, err error) ([]*model.BanAppeal, error) {
	return s.censorBanAppeals(banappeals), err
}

func (s *BanAppealService) censorBanAppeal2(banappeal *model.BanAppeal, err error) (*model.BanAppeal, error) {
	return s.censorBanAppeal(banappeal), err
}