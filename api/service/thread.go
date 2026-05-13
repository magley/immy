package service

import (
	"immy-api/repo"
	"immy-api/model"
)

type ThreadService struct {
	ThreadRepo 	*repo.ThreadRepo
	BoardService *BoardService
	PostService *PostService
}

func (s *ThreadService) ListThreads(offset, limit int) ([]model.Thread, error) {
	return s.ThreadRepo.ListThreads(offset, limit)
}

func (s *ThreadService) ListThreadsOfBoard(boardCode string, offset, limit int) ([]model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}
	
	return s.ThreadRepo.ListThreadsOfBoard(board.ID, offset, limit)
}

func (s *ThreadService) CreateThread(dto model.CreateThreadDTO) (*model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(dto.BoardCode)
	if err != nil {
		return nil, err
	}
	
	thread, err := s.ThreadRepo.CreateThread(dto, board.ID)
	if err != nil {
		return nil, err
	}
	
	_, err = s.PostService.CreatePostForThread(dto.Post, thread, board)
	if err != nil {
		err = s.DeleteThread(thread.ID)
		return nil, err
	}

	return thread, nil
}

func (s *ThreadService) GetThread(threadId uint) (*model.Thread, error) {
	return s.ThreadRepo.GetThread(threadId)
}

func (s *ThreadService) UpdateThread(threadId uint, dto model.UpdateThreadDTO) (*model.Thread, error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	
	return s.ThreadRepo.UpdateThread(thread, dto)
}

func (s *ThreadService) DeleteThread(threadId uint) (error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return err
	}
	
	return s.ThreadRepo.DeleteThread(thread)
}