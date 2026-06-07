package repo

import (
	"gorm.io/gorm"
	model "immy-api/model"
)

type PostRepo struct {
	DB *gorm.DB
}

func (r *PostRepo) ListPosts(offset int, limit int) ([]model.Post, error) {
	var posts []model.Post
	result := r.DB.Limit(limit).Offset(offset).Find(&posts)
	return posts, result.Error
}

func (r *PostRepo) CreatePost(post *model.Post) (*model.Post, error) {
	result := r.DB.Create(post)
	return post, result.Error
}

func (r *PostRepo) GetPost(postId uint) (*model.Post, error) {
	var post model.Post
	result := r.DB.First(&post, postId)
	return &post, result.Error
}

func (r *PostRepo) GetPostsByThread(threadId uint) ([]model.Post, error) {
	var posts []model.Post

	result := r.DB.Where("thread_id = ?", threadId).Find(&posts)
	return posts, result.Error
}

func (r *PostRepo) GetNPostsByThread(threadId uint, n int) ([]model.Post, error) {
	var posts []model.Post

	orderDirection := "id asc"
	absN := n
	if n < 0 {
		orderDirection = "id desc"
		absN *= -1
	}

	result := r.DB.Where("thread_id = ?", threadId).Order(orderDirection).Limit(absN).Find(&posts)
	return posts, result.Error
}

func (r *PostRepo) GetPostWithDuplicateFileInThread(boardId, threadId uint, md5 string) (*model.Post, error) {
	var post model.Post
	result := r.DB.Where("board_id = ?", boardId).Where("thread_id = ?", threadId).Where("md5 = ?", md5).First(&post)
	return &post, result.Error
}

func (r *PostRepo) GetOPPostWithDuplicateFileInBoard(boardId uint, md5 string, ignoreArchivedThreads bool) (*model.Post, error) {
	var post model.Post

	query := r.DB.Model(&model.Post{}).
		Where("posts.board_id = ?", boardId).
		Where("posts.thread_num = posts.num").
		Where("posts.md5 = ?", md5)

	if ignoreArchivedThreads {
		query = query.Joins("LEFT JOIN threads ON threads.id = posts.thread_id").
		Where("threads.archived = false")
	}

	result := query.First(&post)
	return &post, result.Error
}

func (r *PostRepo) GetPostByNum(boardId, postNum uint) (*model.Post, error) {
	var post model.Post
	result := r.DB.Where("num = ?", postNum).Where("board_id = ?", boardId).First(&post)
	return &post, result.Error
}

func (r *PostRepo) UpdatePost(post *model.Post, dto model.UpdatePostDTO) (*model.Post, error) {
	if dto.Name != nil {
		post.Name = *dto.Name
	}
	if dto.Tripcode != nil {
		post.Tripcode = *dto.Tripcode
	}
	if dto.Sage != nil {
		post.Sage = *dto.Sage
	}
	if dto.Content != nil {
		post.Content = *dto.Content
	}
	if dto.Filename != nil {
		post.Filename = *dto.Filename
	}
	if dto.Html != nil {
		post.Html = *dto.Html
	}

	result := r.DB.Save(&post)
	return post, result.Error
}

func (r *PostRepo) DeletePost(post *model.Post) error {
	result := r.DB.Delete(&post)
	return result.Error
}
