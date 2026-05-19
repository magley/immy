package service

import (
	"immy-api/model"
	"immy-api/repo"
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

func (s *ThreadService) CreateThread(dto model.CreateThreadDTO, requestIP string) (*model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(dto.BoardCode)
	if err != nil {
		return nil, err
	}
	
	thread, err := s.ThreadRepo.CreateThread(dto, board.ID)
	if err != nil {
		return nil, err
	}
	
	post, err := s.PostService.CreatePostForThread(dto.Post, requestIP, thread, board)
	if err != nil {
		err = s.DeleteThread(thread.ID)
		return nil, err
	}

	thread, err = s.ThreadRepo.UpdateThreadNum(thread, post.Num)
	if err != nil {
		err = s.DeleteThread(thread.ID)
		err = s.PostService.DeletePost(post.ID)
		// TODO: What about board number?
	}

	return thread, nil
}

func (s *ThreadService) GetThread(threadId uint) (*model.Thread, error) {
	return s.ThreadRepo.GetThread(threadId)
}

func (s *ThreadService) GetThreadByNum(boardCode string, threadNum uint) (*model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}
	
	return s.ThreadRepo.GetThreadByNum(board.ID, threadNum)
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

func (s *ThreadService) GetFullThreadFrom(thread *model.Thread) (*model.ThreadFullDTO, error) {
	posts, err := s.PostService.GetPostsByThread(thread.ID)
	if err != nil {
		return nil, err
	}
	
	return &model.ThreadFullDTO{
		Thread: thread,
		Posts: posts,
	}, nil
}

func (s *ThreadService) GetFullThread(threadId uint) (*model.ThreadFullDTO, error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	
	return s.GetFullThreadFrom(thread)
}

func (s *ThreadService) GetFullThreadByNum(boardCode string, threadNum uint) (*model.ThreadFullDTO, error) {
	thread, err := s.GetThreadByNum(boardCode, threadNum)
		if err != nil {
		return nil, err
	}
	
	return s.GetFullThreadFrom(thread)
}