package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type BoardService struct {
	BoardRepo *repo.BoardRepo
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

func (s *BoardService) GetBoard(boardId uint) (*model.Board, error) {
	return s.BoardRepo.GetBoard(boardId)
}

func (s *BoardService) UpdateBoard(boardCode string, dto model.Board) (*model.Board, error) {
	board, err := s.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}

	return s.BoardRepo.UpdateBoard(board, dto)
}

func (s *BoardService) IncrementBoardPostCount(board *model.Board) (*model.Board, error) {
	board.Meta.PostCount += 1
	return s.BoardRepo.UpdateBoard(board, *board)
}

func (s *BoardService) DecrementBoardPostCount(board *model.Board) (*model.Board, error) {
	board.Meta.PostCount -= 1
	return s.BoardRepo.UpdateBoard(board, *board)
}

func (s *BoardService) IncrementBytesUploaded(board *model.Board, bytes uint) (*model.Board, error) {
	board.Meta.BytesUploaded += bytes
	return s.BoardRepo.UpdateBoard(board, *board)
}

func (s *BoardService) DeleteBoard(boardCode string) error {
	board, err := s.GetBoardByCode(boardCode)
	if err != nil {
		return err
	}

	return s.BoardRepo.DeleteBoard(board)
}

func (s *BoardService) GetStatistics() ([]model.BoardStatisticsDTO, error) {
	return s.BoardRepo.GetStatistics()
}
