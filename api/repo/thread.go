package repo


import (
	"gorm.io/gorm"
	
	model "immy-api/model"
)

type ThreadRepo struct {
	DB *gorm.DB
}

func (r *ThreadRepo) ListThreads(offset int, limit int) ([]model.Thread, error) {
	var threads []model.Thread
	result := r.DB.Limit(limit).Offset(offset).Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepo) ListThreadsOfBoard(boardId uint, offset int, limit int) ([]model.Thread, error) {
	var threads []model.Thread
	result := r.DB.Where("board_id = ?", boardId).Limit(limit).Offset(offset).Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepo) CreateThread(dto model.CreateThreadDTO, boardID uint) (*model.Thread, error) {
	thread := model.Thread{
		BoardID: boardID,
		Subject: dto.Subject,
		Locked: dto.Locked,
		Sticky: dto.Sticky,
	}
	
	result := r.DB.Create(&thread)
	return &thread, result.Error
}

func (r *ThreadRepo) GetThread(threadId string) (*model.Thread, error) {
	var thread model.Thread
	result := r.DB.First(&thread, threadId)
	return &thread, result.Error
}

func (r *ThreadRepo) UpdateThread(threadId string, dto model.UpdateThreadDTO) (*model.Thread, error) {
	thread, err := r.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	
	if dto.Locked != nil { thread.Locked = *dto.Locked }
	if dto.Sticky != nil { thread.Sticky = *dto.Sticky }
	
	result := r.DB.Save(&thread)
	return thread, result.Error
}

func (r *ThreadRepo) DeleteThread(threadId string) (error) {
	thread, err := r.GetThread(threadId)
	if err != nil {
		return err
	}
	
	result := r.DB.Delete(&thread)
	return result.Error
}