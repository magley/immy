package repo

import (
	"time"

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

func (r *ThreadRepo) ListThreadsOfBoardOrderByBump(boardId uint, offset int, limit int) ([]model.Thread, error) {
	var threads []model.Thread
	result := r.DB.
		Model(&model.Thread{}).
		Where("threads.archived = ?", false).
		Where("threads.board_id = ?", boardId).
		Joins("LEFT JOIN posts ON posts.thread_id = threads.id AND posts.sage = false").
		Group("threads.id").
		Order("threads.sticky DESC").
		Order("MAX(posts.created_at) DESC").
		Limit(limit).
		Offset(offset).
		Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepo) ListArchivedThreadsOfBoard(boardId uint, offset int, limit int) ([]model.Thread, error) {
	var threads []model.Thread
	result := r.DB.
		Model(&model.Thread{}).
		Where("threads.archived = ?", true).
		Where("threads.board_id = ?", boardId).
		Order("threads.archived_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepo) GetThreadCountPerBoard(boardId uint) (int64, error) {
	var total int64
	result := r.DB.Model(&model.Thread{}).Where("archived = ?", false).Where("board_id = ?", boardId).Count(&total)
	return total, result.Error
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

func (r *ThreadRepo) GetThread(threadId uint) (*model.Thread, error) {
	var thread model.Thread
	result := r.DB.First(&thread, threadId)
	return &thread, result.Error
}

func (r *ThreadRepo) GetThreadByNum(boardId uint, num uint) (*model.Thread, error) {
	var thread model.Thread
	result := r.DB.Where("board_id = ?", boardId).Where("post_num = ?", num).First(&thread)
	return &thread, result.Error
}

func (r *ThreadRepo) UpdateThread(thread *model.Thread, dto model.UpdateThreadDTO) (*model.Thread, error) {
	if dto.Locked != nil { thread.Locked = *dto.Locked }
	if dto.Sticky != nil { thread.Sticky = *dto.Sticky }
	if dto.AutoCycle != nil { thread.AutoCycle = *dto.AutoCycle }
	
	result := r.DB.Save(&thread)
	return thread, result.Error
}

func (r *ThreadRepo) UpdateThreadNum(thread *model.Thread, num uint) (*model.Thread, error) {
	thread.PostNum = num
	result := r.DB.Save(&thread)
	return thread, result.Error
}

func (r *ThreadRepo) DeleteThread(thread *model.Thread) (error) {
	result := r.DB.Delete(&thread)
	return result.Error
}

func (r *ThreadRepo) GetThreadStats(threadId uint) (model.ThreadStats, error) {
	var stats model.ThreadStats
	result := r.DB.
		Model(&model.Post{}).
		Where("thread_id = ?", threadId).
		Select(
			"COUNT(DISTINCT id) AS post_count, "+
			"SUM(CASE WHEN filename != '' THEN 1 ELSE 0 END) AS image_count, "+
			"COUNT(DISTINCT ipv4) AS user_count, "+
			"MAX(CASE WHEN sage = false THEN created_at ELSE NULL END) AS last_bump").
		Scan(&stats)
	return stats, result.Error
}

func (r *ThreadRepo) ArchiveThread(thread *model.Thread) (*model.Thread, error) {
	thread.ArchivedAt = time.Now()
	thread.Archived = true
	result := r.DB.Save(&thread)
	return thread, result.Error
}
