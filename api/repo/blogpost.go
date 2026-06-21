package repo

import (
	"gorm.io/gorm"

	model "immy-api/model"
)

type BlogpostRepo struct {
	DB *gorm.DB
}

func (r *BlogpostRepo) ListBlogposts(offset int, limit int) ([]*model.Blogpost, error) {
	var blogposts []*model.Blogpost
	result := r.DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&blogposts)
	return blogposts, result.Error
}

func (r *BlogpostRepo) GetBlogpostCount(includeDeleted bool) (int64, error) {
	cnt := int64(0)

	query := r.DB
	if (includeDeleted) {
		query = query.Unscoped()
	}
	result := query.Model(&model.Blogpost{}).Count(&cnt)
	return cnt, result.Error
}

func (r *BlogpostRepo) CreateBlogpost(dto model.CreateBlogpostDTO, creator *model.User) (*model.Blogpost, error) {
	blogpost := model.Blogpost{
		Title: dto.Title,
		Html: dto.Html,
		AuthorID: creator.ID,
		AuthorName: creator.Username,
	}

	result := r.DB.Create(&blogpost)
	return &blogpost, result.Error
}

func (r *BlogpostRepo) GetBlogpost(blogpostId uint) (*model.Blogpost, error) {
	var blogpost model.Blogpost
	result := r.DB.First(&blogpost, blogpostId)
	return &blogpost, result.Error
}

func (r *BlogpostRepo) UpdateBlogpost(blogpost *model.Blogpost, dto model.UpdateBlogpostDTO) (*model.Blogpost, error) {
	if dto.Html != nil { blogpost.Html = *dto.Html }
	result := r.DB.Save(&blogpost)
	return blogpost, result.Error
}

func (r *BlogpostRepo) DeleteBlogpost(blogpost *model.Blogpost) (error) {
	result := r.DB.Delete(&blogpost)
	return result.Error
}