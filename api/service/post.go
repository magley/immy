package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"immy-api/model"
	"immy-api/repo"
	"immy-api/util"
	"net/http"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type PostService struct {
	PostRepo 	*repo.PostRepo
	ThreadRepo 	*repo.ThreadRepo
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

func (s *PostService) GetPostsByThread(threadId uint, includeDeleted bool) ([]model.Post, error) {
	// Because thread ID is a redundant field in the posts table,
	// we should check explicitly if the thread exists anyway.
	thread, err := s.ThreadService.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	return s.PostRepo.GetPostsByThread(thread.ID, includeDeleted)
}

// GetNPostsByThread returns the first `n` posts in the thread. If `n` is negative,
// then the last `n` posts are taken. If `n` is bigger than the number of posts in
// the thread (`m`), only `m` posts are returned.
// The posts are ordered by their ID/number/creation date.
func (s *PostService) GetNPostsByThread(threadId uint, n int) ([]model.Post, error) {
	thread, err := s.ThreadService.GetThread(threadId)
	if err != nil {
		return nil, err
	}
	return s.PostRepo.GetNPostsByThread(thread.ID, n)
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
	if post.ThreadNum == post.Num {
		return s.ThreadService.DeleteThread(post.ThreadID)
	}
	
	return s.PostRepo.DeletePost(post)
}

func (s *PostService) DeleteFirstNPostsOfThread(thread *model.Thread, N uint) (error) {
	return s.PostRepo.DeleteFirstNPostsOfThread(thread.ID, thread.PostNum, N)
}

func (s *PostService) CreatePost(dto model.CreatePostDTO, requestIP string, user *model.User) (*model.Post, error) {
	thread, err := s.ThreadService.GetThread(dto.ThreadID)
	if err != nil {
		return nil, err
	}

	if thread.Locked {
		return nil, errors.New("Thread closed. You may not reply at this time.")
	}

	if thread.Archived {
		return nil, errors.New("Thread archived. You may not reply at this time.")
	}

	threadStats, err := s.ThreadService.GetThreadStats(thread)
	if err != nil {
		return nil, err
	}

	board, err := s.BoardService.GetBoard(thread.BoardID)
	if err != nil {
		return nil, err
	}

	if board.Config.Locked {
		return nil, errors.New("Board locked. You may not reply at this time.")
	}

	err = s.validatePost(dto.Filebytes, thread, threadStats, board, false)
	if err != nil {
		return nil, err
	}

	sage := s.hasOption(dto.Options, "sage")
	if threadStats.PostCount >= board.Config.BumpLimit {
		sage = true
	}

	capcode := s.hasOption(dto.Options, "capcode") && user != nil

	md5 := ""
	if (dto.Filebytes != nil) {
		md5 = util.GetFileHashB64(*dto.Filebytes)

		dupPost, err := s.PostRepo.GetPostWithDuplicateFileInThread(board.ID, thread.ID, md5)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// We want this.
			} else {
				return nil, err
			}
		} else {
			if dupPost != nil {
				return nil, errors.New(fmt.Sprintf("Duplicate file at #%d", dupPost.Num))
			}
		}
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

	var publicID *string
	if board.Config.IDsEnabled {
		publicIDstr := util.CreateUserID(requestIP, thread.ID)
		publicID = &publicIDstr
	}

	var userId *uint
	var userRole *model.UserRole
	if user != nil { userId = &user.ID }
	if user != nil { userRole = &user.Role }

	post := &model.Post{
		ThreadID: thread.ID,
		ThreadNum: thread.PostNum,
		BoardID: board.ID,
		Num: board.Meta.PostCount,
		Name: postName,
		Tripcode: postTripcode,
		IPv4: requestIP,
		UserID: userId,
		UserRole: userRole,
		PublicID: publicID,
		Sage: sage,
		Capcode: capcode,
		Content: dto.Content,
		SrcFilename: "",
		Filename: "",
		ImgWidth: 0,
		ImgHeight: 0,
		MD5: "",
		Spoiler: dto.Spoiler,
		Html: "",
	}

	if (dto.Filename != nil && dto.Filebytes != nil) {
		post.SrcFilename = *dto.Filename
		post.Filename = util.GetPostImageFilename(board.Code, post.SrcFilename)
	}

	if (dto.Filename != nil && dto.Filebytes != nil) {
		imgData, err := util.SaveFile(post.Filename, *dto.Filebytes)
		if err != nil {
			return nil, err
		}

		board, err = s.BoardService.IncrementBytesUploaded(board, imgData.SizeImageBytes)
		if err != nil {
			return nil, err
		}

		post.MD5 = md5
		post.Filesize = imgData.SizeImageBytes
		post.ImgWidth = uint(imgData.ImageWidth)
		post.ImgHeight = uint(imgData.ImageHeight)
	}

	post, err = s.PostRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	err = s.ThreadService.UpdateAutoCycleForThread(thread)
	if err != nil {
		return nil, err
	}

	return post, err 	
}

func (s *PostService) CreatePostForThread(dto model.CreatePostForThreadDTO, requestIP string, thread *model.Thread, board *model.Board, user *model.User) (*model.Post, error) {
	if board.Config.Locked {
		return nil, errors.New("Board locked. You may not create threads at this time.")
	}

	board, err := s.BoardService.IncrementBoardPostCount(board)
	if err != nil {
		return nil, err
	}

	threadStats, err := s.ThreadService.GetThreadStats(thread)
	if err != nil {
		return nil, err
	}

	err = s.validatePost(&dto.Filebytes, thread, threadStats, board, true)
	if err != nil {
		return nil, err
	}

	sage := s.hasOption(dto.Options, "sage")
	if threadStats.PostCount >= board.Config.BumpLimit {
		sage = true
	}

	capcode := s.hasOption(dto.Options, "capcode") && user != nil

	md5 := ""
	{
		md5 = util.GetFileHashB64(dto.Filebytes)

		dupPost, err := s.PostRepo.GetOPPostWithDuplicateFileInBoard(board.ID, md5, true)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// We want this.
			} else {
				return nil, err
			}
		} else {
			if dupPost != nil {
				return nil, errors.New(fmt.Sprintf("Duplicate file at #%d", dupPost.Num))
			}
		}
	}

	dto.Name = strings.Trim(dto.Name, " \t")
	dto.Content = strings.Trim(dto.Content, " \t")
	dto.Options = strings.Trim(dto.Options, " \t")
	
	postName, postTripcode := s.createTripcode(dto.Name)
	
	var publicID *string
	if board.Config.IDsEnabled {
		publicIDstr := util.CreateUserID(requestIP, thread.ID)
		publicID = &publicIDstr
	}

	var userId *uint
	var userRole *model.UserRole
	if user != nil { userId = &user.ID }
	if user != nil { userRole = &user.Role }

	post := &model.Post{
		ThreadID: thread.ID,
		ThreadNum: board.Meta.PostCount,
		BoardID: board.ID,
		Num: board.Meta.PostCount,
		Name: postName,
		Tripcode: postTripcode,
		IPv4: requestIP,
		UserID: userId,
		UserRole: userRole,
		PublicID: publicID,
		Sage: sage,
		Capcode: capcode,
		Content: dto.Content,
		SrcFilename: dto.Filename,
		Filename: util.GetPostImageFilename(board.Code, dto.Filename),
		MD5: md5,
		Spoiler: dto.Spoiler,
		Html: "",
	}
	
	imgData, err := util.SaveFile(post.Filename, dto.Filebytes)
	if err != nil {
		return nil, err
	}

	post.Filesize = imgData.SizeImageBytes
	post.ImgWidth = uint(imgData.ImageWidth)
	post.ImgHeight = uint(imgData.ImageHeight)

	board, err = s.BoardService.IncrementBytesUploaded(board, imgData.SizeImageBytes)
	if err != nil {
		return nil, err
	}

	post, err = s.PostRepo.CreatePost(post)

	return post, err 
}

func (s *PostService) validatePost(fileBytes *string, thread *model.Thread, threadStats model.ThreadStats, board *model.Board, opPost bool) error {
	if fileBytes != nil {
		if threadStats.ImageCount >= board.Config.ImageLimit {
			return errors.New("Image limit reached")
		}

		if !board.Config.ReplyFilesAllowed && !opPost {
			return errors.New("Only OP can attach files")
		}

    	data, err := base64.StdEncoding.DecodeString(*fileBytes)
    	if err != nil {
    		return err
    	}

    	if len(data) > int(board.Config.MaxFileSize) {
    		return errors.New("File too large")
    	}

    	mimeType := http.DetectContentType(data[:512])
    	mimeOk := false
    	for _, mime := range board.Config.MimeTypesAllowed {
    		if mime == mimeType {
    			mimeOk = true
    			break
    		}
    	}

    	if !mimeOk {
    		return errors.New(fmt.Sprintf("Unsupported file type: %s", mimeType))
    	}
	}

	return nil
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

func (s *PostService) hasOption(options string, option string) bool {
	parts := strings.Split(options, " ")
	return slices.Index(parts, option) != -1
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
