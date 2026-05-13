package service

import (
	"fmt"
	"strings"
	"immy-api/util"
	"immy-api/repo"
	"immy-api/model"
)

type PostService struct {
	PostRepo 	*repo.PostRepo
	BoardService *BoardService
}

func (s *PostService) CreatePostForThread(dto model.CreatePostForThreadDTO, thread *model.Thread, board *model.Board) (*model.Post, error) {
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
		Num: board.PostCount,
		Name: postName,
		Tripcode: postTripcode,
		IPv4: "unknown",
		Sage: s.isSage(dto.Options),
		Content: dto.Content,
		Filename: "unknown or null...",
		Html: "",
	}
	
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