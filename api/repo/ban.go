package repo

import (
	util "immy-api/util"

	"gorm.io/gorm"

	model "immy-api/model"
)

type BanRepo struct {
	DB *gorm.DB
}

func (r *BanRepo) ListBans(offset int, limit int) ([]model.Ban, error) {
	var bans []model.Ban
	result := r.DB.Limit(limit).Offset(offset).Find(&bans)
	return bans, result.Error
}

func (r *BanRepo) ListBansForAdmin(offset int, limit int) ([]model.Ban, error) {
	var bans []model.Ban
	result := r.DB.
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(&bans).
		Order("created_at desc").
		Order("expires_at desc")
	return bans, result.Error
}

func (r *BanRepo) GetBansOfIp(ip string) ([]model.Ban, error) {
	ipBinary := util.IPv4toUint64(ip)
	var bans []model.Ban
	result := r.DB.
		Model(&model.Ban{}).
		Where("(expires_at is null or (expires_at > now() or seen = false))").
		Where("((ip_end is null and ip_start = ?) or (? between ip_start and ip_end))", ipBinary, ipBinary).
		Order("created_at desc").
		Order("expires_at desc").
		Find(&bans)
	return bans, result.Error
}

func (r *BanRepo) CreateBan(dto model.CreateBanDTO, creator *model.User) (*model.Ban, error) {
	ban := model.Ban{
		IpStart: util.IPv4toUint64(dto.IpStart),
		// IpEnd: ,
		ExpiresAt: dto.ExpiresAt,
		BoardID: dto.BoardID,
		CreatorID: creator.ID,
		Reason: dto.Reason,
		Warning: dto.Warning,
	}
	if (ban.Warning) {
		ban.ExpiresAt = nil
	}

	if dto.IpEnd != nil {
		ipEnd := util.IPv4toUint64(*dto.IpEnd)
		ban.IpEnd = &ipEnd
	}

	result := r.DB.Create(&ban)
	return &ban, result.Error
}

func (r *BanRepo) GetBan(banId uint) (*model.Ban, error) {
	var ban model.Ban
	result := r.DB.First(&ban, banId)
	return &ban, result.Error
}

func (r *BanRepo) GetBanForAdmin(banId uint) (*model.Ban, error) {
	var ban model.Ban
	result := r.DB.Unscoped().First(&ban, banId)
	return &ban, result.Error
}

func (r *BanRepo) UpdateBan(ban *model.Ban, dto model.UpdateBanDTO) (*model.Ban, error) {
	if dto.Seen != nil { ban.Seen = *dto.Seen}

	result := r.DB.Save(&ban)
	return ban, result.Error
}

func (r *BanRepo) DeleteBan(ban *model.Ban) (error) {
	result := r.DB.Delete(&ban)
	return result.Error
}