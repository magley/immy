package repo

import (
	"gorm.io/gorm"
	
	model "immy-api/model"
)

type BoardRepo struct {
	DB *gorm.DB
}

func (r *BoardRepo) ListBoards(offset int, limit int) ([]model.Board, error) {
	var boards []model.Board
	result := r.DB.Limit(limit).Offset(offset).Order("id").Find(&boards)
	return boards, result.Error
}

func (r *BoardRepo) CreateBoard(dto model.CreateBoardDTO) (*model.Board, error) {
	board := model.Board{
		Name: dto.Name,
		Code: dto.Code,
		Description: dto.Description,
	}
	
	result := r.DB.Create(&board)
	return &board, result.Error
}

func (r *BoardRepo) GetBoardByCode(boardCode string) (*model.Board, error) {
	var board model.Board
	result := r.DB.Where("code = ?", boardCode).First(&board)
	return &board, result.Error
}

func (r *BoardRepo) GetBoard(boardId uint) (*model.Board, error) {
	var board model.Board
	result := r.DB.First(&board, boardId)
	return &board, result.Error
}

func (r *BoardRepo) UpdateBoard(board *model.Board, dto model.UpdateBoardDTO) (*model.Board, error) {
	if dto.Name != nil { board.Name = *dto.Name }
	if dto.Code != nil { board.Code = *dto.Code }
	if dto.Description != nil { board.Description = dto.Description }
	if dto.Locked != nil { board.Locked = *dto.Locked }
	if dto.Hidden != nil { board.Hidden = *dto.Hidden }
	if dto.PostCount != nil { board.PostCount = *dto.PostCount }

	result := r.DB.Save(&board)
	return board, result.Error
}

func (r *BoardRepo) DeleteBoard(board *model.Board) (error) {
	result := r.DB.Delete(&board)
	return result.Error
}