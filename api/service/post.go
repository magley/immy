package service

import (
	"fmt"
	"immy-api/model"
	"immy-api/repo"
	"immy-api/util"
	"strings"
)

type PostService struct {
	PostRepo 	*repo.PostRepo
	BoardService *BoardService
	ThreadService *ThreadService
}

func (s *PostService) ListPosts(offset, limit int) ([]model.Post, error) {
	return s.PostRepo.ListPosts(offset, limit)
}

func (s *PostService) GetPost(postId uint) (*model.Post, error) {
	return s.PostRepo.GetPost(postId)
}

func (s *PostService) GetPostByNum(boardCode string, postNum uint) (*model.Post, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}
	return s.PostRepo.GetPostByNum(board.ID, postNum)
}

func (s *PostService) GetPostsByThread(threadId uint) ([]model.Post, error) {
	// Because thread ID is a redundant field in the posts table,
	// we should check explicitly if the thread exists anyway.
	thread, err := s.ThreadService.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	return s.PostRepo.GetPostsByThread(thread.ID)
}

func (s *PostService) UpdatePost(postId uint, dto model.UpdatePostDTO) (*model.Post, error) {
	post, err := s.GetPost(postId)
	if err != nil {
		return nil, err
	}
	
	return s.PostRepo.UpdatePost(post, dto)
}

func (s *PostService) DeletePost(postId uint) (error) {
	post, err := s.GetPost(postId)
	if err != nil {
		return err
	}
	
	return s.PostRepo.DeletePost(post)
}

func (s *PostService) CreatePost(dto model.CreatePostDTO, requestIP string) (*model.Post, error) {
	thread, err := s.ThreadService.GetThread(dto.ThreadID)
	if err != nil {
		return nil, err
	}
	
	board, err := s.BoardService.GetBoard(thread.BoardID)
	if err != nil {
		return nil, err
	}
	
	// TODO: This is identical to CreatePostForThread, so make
	// a common method for both of them. Maybe in the future,
	// CreatePostForThread will actually behave differently
	// (like having additional fields in the Post to signify
	// that this is an OP post). One exception is the ThreadNum
	// field.
	
	board, err = s.BoardService.IncrementBoardPostCount(board)
	if err != nil {
		return nil, err
	}
	
	dto.Name = strings.Trim(dto.Name, " \t")
	dto.Content = strings.Trim(dto.Content, " \t")
	dto.Options = strings.Trim(dto.Options, " \t")
	
	postName, postTripcode := s.createTripcode(dto.Name)

	post := &model.Post{
		ThreadID: thread.ID,
		ThreadNum: thread.PostNum,
		BoardID: board.ID,
		Num: board.PostCount,
		Name: postName,
		Tripcode: postTripcode,
		IPv4: requestIP,
		Sage: s.isSage(dto.Options),
		Content: dto.Content,
		SrcFilename: "",
		Filename: "",
		Html: "",
	}
	if (dto.Filename != nil && dto.Filebytes != nil) {
		post.SrcFilename = *dto.Filename
		post.Filename = util.GetPostImageFilename(board.Code, post.SrcFilename)
	}

	if (dto.Filename != nil && dto.Filebytes != nil) {
		bytesImg, _, err := util.SaveImage(post.Filename, *dto.Filebytes)
		if err != nil {
			return nil, err
		}

		post.Filesize = bytesImg
	}

	post, err = s.PostRepo.CreatePost(post)

	return post, err 	
}

func (s *PostService) CreatePostForThread(dto model.CreatePostForThreadDTO, requestIP string, thread *model.Thread, board *model.Board) (*model.Post, error) {
	board, err := s.BoardService.IncrementBoardPostCount(board)
	if err != nil {
		return nil, err
	}
	
	dto.Name = strings.Trim(dto.Name, " \t")
	dto.Content = strings.Trim(dto.Content, " \t")
	dto.Options = strings.Trim(dto.Options, " \t")
	
	postName, postTripcode := s.createTripcode(dto.Name)
	
	post := &model.Post{
		ThreadID: thread.ID,
		ThreadNum: board.PostCount,
		BoardID: board.ID,
		Num: board.PostCount,
		Name: postName,
		Tripcode: postTripcode,
		IPv4: requestIP,
		Sage: s.isSage(dto.Options),
		Content: dto.Content,
		SrcFilename: dto.Filename,
		Filename: util.GetPostImageFilename(board.Code, dto.Filename),
		Html: "",
	}
	
	bytesImg, _, err := util.SaveImage(post.Filename, dto.Filebytes)
	if err != nil {
		return nil, err
	}

	post.Filesize = bytesImg

	post, err = s.PostRepo.CreatePost(post)

	return post, err 
}

func (s *PostService) createTripcode(fullName string) (string, string) {
	parts, secure := splitCustom(fullName)
	
	
	if len(parts) < 2 {
		return fullName, ""
	}
	if len(parts[1]) == 0 {
		return fullName, ""
	}
	fmt.Printf("%s,%s", parts[0],parts[1])
	
	return parts[0], util.CreateTripcode(parts[1], secure)
}

func (s *PostService) isSage(options string) bool {
	return options == "sage"
}


// splitCustom splits a string into two parts: the username and the tripcode password.
// The two are separated by a '#'. If they are separated by two '#', then the tripcode
// is meant to be secure, and the second return value is true. Otherwise, the tripcode
// is meant to be insecure, so the second return value is false.
//
//
// name          -> ([name], false)
// name#pass     -> ([name, pass], false)
// name##pass    -> ([name, pass], true)
// name###pass   -> (name, #pass], true)
//
func splitCustom(s string) ([]string, bool) {
    idx := strings.Index(s, "#")
    if idx == -1 {
        return []string{s}, false
    }

    count := 1
    for i := idx + 1; i < len(s) && s[i] == '#'; i++ {
        count++
    }

    switch count {
    case 1:
        return []string{s[:idx], s[idx+1:]}, false
    case 2:
        return []string{s[:idx], s[idx+2:]}, true
    default:
        return []string{s[:idx], s[idx+count-2:]}, true
    }
}