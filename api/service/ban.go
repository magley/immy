package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type BanService struct {
	BanRepo 		*repo.BanRepo
	BoardService 	*BoardService
	UserService 	*UserService
}

func (s *BanService) ListBans(offset, limit int) ([]*model.Ban, int64, error) {
	bans, err := s.censorBans2(s.BanRepo.ListBans(offset, limit))
	if err != nil {
		return bans, 0, err
	}
	totalCount, err := s.BanRepo.GetBanCount(true)
	if err != nil {
		return bans, 0, err
	}
	return bans, totalCount, err
}

func (s *BanService) ListBansForAdmin(offset, limit int) ([]*model.BanExtDTO, int64, error) {
	bans, err := s.BanRepo.ListBansForAdmin(offset, limit)
	var result []*model.BanExtDTO
	if err != nil {
		return result, 0, err
	}
	totalCount, err := s.BanRepo.GetBanCount(true)
	if err != nil {
		return result, 0, err
	}

	for _, ban := range bans {
		var boardCode *string = nil

		if ban.BoardID != nil {
			board, err := s.BoardService.GetBoard(*ban.BoardID)
			if err != nil { return result, 0, err }
			boardCode = &board.Code
		}

		user, err := s.UserService.GetUser(ban.CreatorID)
		if err != nil { return result, 0, err }

		ban2 := &model.BanExtDTO{
			Ban: *ban,
			BoardCode: boardCode,
			CreatorUsername: user.Username,
		}
	 	result = append(result, ban2)
	}

	return result, totalCount, err
}

func (s *BanService) GetBansOfIp(ip string) ([]*model.Ban, error) {
	return s.censorBans2(s.BanRepo.GetBansOfIp(ip))
}

func (s *BanService) CreateBan(dto model.CreateBanDTO, creator *model.User) (*model.Ban, error) {
	return s.BanRepo.CreateBan(dto, creator)
}

func (s *BanService) GetBan(banId uint) (*model.Ban, error) {
	return s.censorBan2(s.BanRepo.GetBan(banId))
}

func (s *BanService) GetBanForAdmin(banId uint) (*model.Ban, error) {
	return s.BanRepo.GetBanForAdmin(banId)
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

// =====================================================================================

func (s *BanService) censorBan(ban *model.Ban) (*model.Ban) {
	ban.CreatorID = 0
	ban.IpStart = 0
	ban.IpEnd = nil
	return ban
}

func (s *BanService) censorBans(bans []*model.Ban) ([]*model.Ban) {
	for _, ban := range bans {
		ban.CreatorID = 0
		ban.IpStart = 0
		ban.IpEnd = nil
	}
	return bans
}

func (s *BanService) censorBans2(bans []*model.Ban, err error) ([]*model.Ban, error) {
	return s.censorBans(bans), err
}


func (s *BanService) censorBan2(ban *model.Ban, err error) (*model.Ban, error) {
	return s.censorBan(ban), err
}