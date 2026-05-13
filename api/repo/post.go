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

func (r *PostRepo) CreatePost(dto model.CreatePostDTO) (*model.Post, error) {
	postNumber := uint(0)
	
	post := model.Post{
		ThreadID: dto.ThreadID,
		Num: postNumber,
		Name: dto.Name,
		// Tripcode: ...,
		// IPv4: ...,
		// Sage: ...,
		Content: dto.Content,
		Filename: dto.Filename,
		Html: "",
	}
	
	result := r.DB.Create(&post)
	return &post, result.Error
}


// func (r *PostRepo) CreatePostForThread(dto model.CreatePostForThreadDTO, threadID uint) (*model.Post, error) {
	
	
	
// 	ost := model.Post{
// 		ThreadID: threadID,
// 		Num: postNumber
// 		Name: dto.Name,
// 		Tripcode: ...,
// 		IPv4: ...,
// 		Sage: ...,
// 		Content: dto.Content,
// 		Filename: dto.Filename,
// 		Html: "",
// 	}
	
// 	result := r.DB.Create(&post)
// 	return &post, result.Error
// }


func (r *PostRepo) GetPost(postCode string) (*model.Post, error) {
	var post model.Post
	result := r.DB.Where("code = ?", postCode).First(&post)
	return &post, result.Error
}

func (r *PostRepo) UpdatePost(postCode string, dto model.UpdatePostDTO) (*model.Post, error) {
	post, err := r.GetPost(postCode)
	if err != nil {
		return nil, err
	}
	
	// if dto.Name != nil { post.Name = *dto.Name }
	// if dto.Code != nil { post.Code = *dto.Code }
	// if dto.Description != nil { post.Description = dto.Description }
	// if dto.Locked != nil { post.Locked = *dto.Locked }
	// if dto.Hidden != nil { post.Hidden = *dto.Hidden }
	
	result := r.DB.Save(&post)
	return post, result.Error
}

func (r *PostRepo) DeletePost(postCode string) (error) {
	post, err := r.GetPost(postCode)
	if err != nil {
		return err
	}
	
	result := r.DB.Delete(&post)
	return result.Error
}