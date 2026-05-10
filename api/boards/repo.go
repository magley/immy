package boards

import (
	"gorm.io/gorm"
)

type BoardRepo struct {
	DB *gorm.DB
}

func (r *BoardRepo) ListBoards(offset int, limit int) ([]Board, error) {
	var boards []Board
	result := r.DB.Limit(limit).Offset(offset).Find(&boards)
	return boards, result.Error
}

func (r *BoardRepo) CreateBoard(dto CreateBoardDTO) (*Board, error) {
	board := Board{
		Name: dto.Name,
		Code: dto.Code,
		Description: dto.Description,
	}
	
	result := r.DB.Create(&board)
	return &board, result.Error
}

func (r *BoardRepo) GetBoard(boardCode string) (*Board, error) {
	var board Board
	result := r.DB.Where("code = ?", boardCode).First(&board)
	return &board, result.Error
}

func (r *BoardRepo) UpdateBoard(boardCode string, dto UpdateBoardDTO) (*Board, error) {
	board, err := r.GetBoard(boardCode)
	if err != nil {
		return nil, err
	}
	
	if dto.Name != nil { board.Name = *dto.Name }
	if dto.Code != nil { board.Code = *dto.Code }
	if dto.Description != nil { board.Description = dto.Description }
	if dto.Locked != nil { board.Locked = *dto.Locked }
	if dto.Hidden != nil { board.Hidden = *dto.Hidden }
	
	result := r.DB.Save(&board)
	return board, result.Error
}

func (r *BoardRepo) DeleteBoard(boardCode string) (error) {
	board, err := r.GetBoard(boardCode)
	if err != nil {
		return err
	}
	
	result := r.DB.Delete(&board)
	
	return result.Error
}