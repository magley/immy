package service

import (
	"immy-api/repo"
	"immy-api/model"
)

type BoardService struct {
	BoardRepo 	*repo.BoardRepo
}

func (s *BoardService) ListBoards(offset, limit int) ([]model.Board, error) {
	return s.BoardRepo.ListBoards(offset, limit)
}

func (s *BoardService) CreateBoard(dto model.CreateBoardDTO) (*model.Board, error) {
	return s.BoardRepo.CreateBoard(dto)
}

func (s *BoardService) GetBoardByCode(code string) (*model.Board, error) {
	return s.BoardRepo.GetBoardByCode(code)
}

func (s *BoardService) UpdateBoard(boardCode string, dto model.UpdateBoardDTO) (*model.Board, error) {
	board, err := s.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}
	
	return s.BoardRepo.UpdateBoard(board, dto)
}

func (s *BoardService) DeleteBoard(boardCode string) (error) {
	board, err := s.GetBoardByCode(boardCode)
	if err != nil {
		return err
	}
	
	return s.BoardRepo.DeleteBoard(board)
}