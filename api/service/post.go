package service

import (
	"immy-api/repo"
	"immy-api/model"
)

type PostService struct {
	PostRepo 	*repo.PostRepo
}

func (s *PostService) CreatePostForThread(dto model.CreatePostForThreadDTO, thread *model.Thread) (*model.Post, error) {
	// At this point, the thread has been created.
	// We need to obtain the board's post count.
	
	return nil, nil
}