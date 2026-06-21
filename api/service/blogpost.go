package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type BlogpostService struct {
	BlogpostRepo 		*repo.BlogpostRepo
	UserService 	*UserService
}

func (s *BlogpostService) ListBlogposts(offset, limit int) ([]*model.Blogpost, int64, error) {
	blogposts, err := s.BlogpostRepo.ListBlogposts(offset, limit)
	if err != nil {
		return blogposts, 0, err
	}
	totalCount, err := s.BlogpostRepo.GetBlogpostCount(true)
	if err != nil {
		return blogposts, 0, err
	}
	return blogposts, totalCount, err
}

func (s *BlogpostService) CreateBlogpost(dto model.CreateBlogpostDTO, creator *model.User) (*model.Blogpost, error) {
	return s.BlogpostRepo.CreateBlogpost(dto, creator)
}

func (s *BlogpostService) GetBlogpost(blogpostId uint) (*model.Blogpost, error) {
	return s.BlogpostRepo.GetBlogpost(blogpostId)
}

func (s *BlogpostService) UpdateBlogpost(blogpostId uint, dto model.UpdateBlogpostDTO) (*model.Blogpost, error) {
	blogpost, err := s.GetBlogpost(blogpostId)
	if err != nil {
		return nil, err
	}

	return s.BlogpostRepo.UpdateBlogpost(blogpost, dto)
}

func (s *BlogpostService) DeleteBlogpost(blogpostId uint) (error) {
	blogpost, err := s.GetBlogpost(blogpostId)
	if err != nil {
		return nil
	}

	return s.BlogpostRepo.DeleteBlogpost(blogpost)
}

func (s *BlogpostService) ToShort(blogpost *model.Blogpost) (*model.BlogpostShortDTO) {
	return &model.BlogpostShortDTO{
		ID: blogpost.ID,
		Title: blogpost.Title,
		CreatedAt: blogpost.CreatedAt,
	}
}

func (s *BlogpostService) ToShortArr(blogposts []*model.Blogpost) ([]*model.BlogpostShortDTO) {
	var result []*model.BlogpostShortDTO
	for _, blogpost := range blogposts {
		result = append(result, s.ToShort(blogpost))
	}
	return result
}