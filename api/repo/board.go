package repo

import (
	"log"

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
		Config: dto.Config,
	}
	log.Println(dto.Config)
	
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

func (r *BoardRepo) UpdateBoard(board *model.Board, dto model.Board) (*model.Board, error) {
	board.Name = dto.Name
	board.Code = dto.Code
	board.Description = dto.Description
	board.PostCount = dto.PostCount
	board.Config = dto.Config

	result := r.DB.Save(&board)
	return board, result.Error
}

func (r *BoardRepo) DeleteBoard(board *model.Board) (error) {
	result := r.DB.Delete(&board)
	return result.Error
}

func (r *BoardRepo) GetStatistics() ([]model.BoardStatisticsDTO, error) {
	sql := `
		select boards.id, boards.code, count(distinct threads.id) thread_count, count(distinct posts.id) post_count
		from boards
		join threads on boards.id = threads.board_id
		join posts on boards.id = posts.board_id
		group by boards.id
	;
	`

	var stats []model.BoardStatisticsDTO
	result := r.DB.Raw(sql).Find(&stats)

	return stats, result.Error
}