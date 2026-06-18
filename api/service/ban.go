package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type BanService struct {
	BanRepo 	*repo.BanRepo
}

func (s *BanService) ListBans(offset, limit int) ([]model.Ban, error) {
	return s.BanRepo.ListBans(offset, limit)
}

func (s *BanService) ListBansForAdmin(offset, limit int) ([]model.Ban, error) {
	return s.BanRepo.ListBansForAdmin(offset, limit)
}

func (s *BanService) GetBansOfIp(ip string) ([]model.Ban, error) {
	return s.BanRepo.GetBansOfIp(ip)
}

func (s *BanService) CreateBan(dto model.CreateBanDTO, creator *model.User) (*model.Ban, error) {
	return s.BanRepo.CreateBan(dto, creator)
}

func (s *BanService) GetBan(banId uint) (*model.Ban, error) {
	return s.BanRepo.GetBan(banId)
}

func (s *BanService) UpdateBan(banId uint, dto model.UpdateBanDTO) (*model.Ban, error) {
	ban, err := s.GetBan(banId)
	if err != nil {
		return nil, err
	}

	return s.BanRepo.UpdateBan(ban, dto)
}

func (s *BanService) DeleteBan(banId uint) (error) {
	ban, err := s.GetBan(banId)
	if err != nil {
		return nil
	}

	return s.BanRepo.DeleteBan(ban)
}