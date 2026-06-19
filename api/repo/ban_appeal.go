package repo

import (
	"time"

	"gorm.io/gorm"

	model "immy-api/model"
)

type BanAppealRepo struct {
	DB *gorm.DB
}

func (r *BanAppealRepo) ListBanAppeals(offset int, limit int) ([]*model.BanAppeal, error) {
	var banappeals []*model.BanAppeal
	result := r.DB.Limit(limit).Offset(offset).Find(&banappeals)
	return banappeals, result.Error
}

func (r *BanAppealRepo) GetBanAppealCount(includeDeleted bool) (int64, error) {
	cnt := int64(0)

	query := r.DB
	if (includeDeleted) {
		query = query.Unscoped()
	}
	result := query.Model(&model.BanAppeal{}).Count(&cnt)
	return cnt, result.Error
}

func (r *BanAppealRepo) GetBanAppeal(banappealId uint) (*model.BanAppeal, error) {
	var banappeal model.BanAppeal
	result := r.DB.First(&banappeal, banappealId)
	return &banappeal, result.Error
}

func (r *BanAppealRepo) CreateBanAppeal(dto model.CreateBanAppealDTO, ban *model.Ban) (*model.BanAppeal, error) {
	banappeal := model.BanAppeal{
		BanID: ban.ID,
		Message: dto.Message,
	}

	result := r.DB.Create(&banappeal)
	return &banappeal, result.Error
}

func (r *BanAppealRepo) UpdateBanAppeal(banappeal *model.BanAppeal, dto model.UpdateBanAppealDTO, user *model.User) (*model.BanAppeal, error) {
	if dto.Status != nil {
		now := time.Now()

		banappeal.Status = *dto.Status
		banappeal.ReviewedBy = &user.ID
		banappeal.ReviewedAt = &now
	}

	result := r.DB.Save(&banappeal)
	return banappeal, result.Error
}

func (r *BanAppealRepo) DeleteBanAppeal(banappeal *model.BanAppeal) (error) {
	result := r.DB.Delete(&banappeal)
	return result.Error
}

// -------------------------------------

func (r *BanAppealRepo) GetBanAppealsOfBan(ban *model.Ban) ([]*model.BanAppeal, error) {
	var banappeals []*model.BanAppeal
	result := r.DB.Model(&model.BanAppeal{}).Where("ban_id = ?", ban.ID).Find(&banappeals)
	return banappeals, result.Error
}

func (r *BanAppealRepo) CanAppealBan(ban *model.Ban) (bool, error) {
	var finalRejectedCount int64
	result := r.DB.Model(&model.BanAppeal{}).
		Where("ban_id = ?", ban.ID).
		Where("status = ?", model.BanAppealStatusRejectedFinal).
		Count(&finalRejectedCount)
	return finalRejectedCount == 0, result.Error
}