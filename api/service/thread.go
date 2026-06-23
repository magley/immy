package service

import (
	"immy-api/model"
	"immy-api/repo"
	"log"
	"slices"
)

type ThreadService struct {
	ThreadRepo   *repo.ThreadRepo
	BoardService *BoardService
	PostService  *PostService
}

func (s *ThreadService) ListThreads(offset, limit int) ([]model.Thread, error) {
	return s.ThreadRepo.ListThreads(offset, limit)
}

func (s *ThreadService) ListThreadsOfBoard(boardCode string, offset, limit int) ([]model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}

	return s.ThreadRepo.ListThreadsOfBoard(board.ID, offset, limit)
}

func (s *ThreadService) ListThreadsOfBoardOrderByBump(boardCode string, offset, limit int) ([]model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}

	return s.ThreadRepo.ListThreadsOfBoardOrderByBump(board.ID, offset, limit)
}

func (s *ThreadService) ListArchivedThreadsOfBoard(boardCode string, offset, limit int) ([]model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}

	return s.ThreadRepo.ListArchivedThreadsOfBoard(board.ID, offset, limit)
}

func (s *ThreadService) GetThreadCountPerBoard(boardCode string) (int64, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return 0, err
	}

	return s.ThreadRepo.GetThreadCountPerBoard(board.ID)
}

func (s *ThreadService) CreateThread(dto model.CreateThreadDTO, requestIP string, user *model.User) (*model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(dto.BoardCode)
	if err != nil {
		return nil, err
	}

	thread, err := s.ThreadRepo.CreateThread(dto, board.ID)
	if err != nil {
		return nil, err
	}

	post, err := s.PostService.CreatePostForThread(dto.Post, requestIP, thread, board, user)
	if err != nil {
		err = s.DeleteThread(thread.ID)
		return nil, err
	}

	thread, err = s.ThreadRepo.UpdateThreadNum(thread, post.Num)
	if err != nil {
		err = s.DeleteThread(thread.ID)
		err = s.PostService.DeletePost(post.ID)
		// TODO: What about board number?
	}

	err = s.ArchiveBumpedOffThreads(board)

	return thread, err // Last error isn't major enough to return nil.
}

func (s *ThreadService) GetThread(threadId uint) (*model.Thread, error) {
	return s.ThreadRepo.GetThread(threadId)
}

func (s *ThreadService) GetThreadByNum(boardCode string, threadNum uint) (*model.Thread, error) {
	board, err := s.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		return nil, err
	}

	return s.ThreadRepo.GetThreadByNum(board.ID, threadNum)
}

func (s *ThreadService) UpdateThread(threadId uint, dto model.UpdateThreadDTO) (*model.Thread, error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return nil, err
	}

	return s.ThreadRepo.UpdateThread(thread, dto)
}

func (s *ThreadService) DeleteThread(threadId uint) error {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return err
	}

	return s.ThreadRepo.DeleteThread(thread)
}

func (s *ThreadService) GetFullThreadFrom(thread *model.Thread) (*model.ThreadFullDTO, error) {
	posts, err := s.PostService.GetPostsByThread(thread.ID, thread.Archived)
	if err != nil {
		return nil, err
	}

	return &model.ThreadFullDTO{
		Thread: thread,
		Posts:  posts,
	}, nil
}

func (s *ThreadService) GetThreadsForCatalog(boardCode string) ([]model.ThreadForCatalogDTO, error) {
	// TODO: Hardcoding it like this is bad.
	threads, err := s.ListThreadsOfBoardOrderByBump(boardCode, 0, 1000)
	if err != nil {
		return nil, err
	}

	var res []model.ThreadForCatalogDTO

	for _, thread := range threads {
		post, err := s.PostService.GetPostByNum(boardCode, thread.PostNum)
		if err != nil {
			return nil, err
		}
		stats, err := s.GetThreadStats(&thread)
		if err != nil {
			return nil, err
		}
		lastPos, err := s.PostService.GetNPostsByThread(thread.ID, -1)
		if err != nil {
			return nil, err
		}

		threadWithPost := model.ThreadForCatalogDTO{
			Thread:   thread,
			Post:     *post,
			Stats:    stats,
			LastPost: lastPos[0],
		}
		res = append(res, threadWithPost)
	}

	return res, nil
}

func (s *ThreadService) GetThreadsForArchive(boardCode string) ([]model.ThreadForCatalogDTO, error) {
	// TODO: Hardcoding it like this is bad.
	threads, err := s.ListArchivedThreadsOfBoard(boardCode, 0, 1000)
	if err != nil {
		return nil, err
	}

	var res []model.ThreadForCatalogDTO

	for _, thread := range threads {
		post, err := s.PostService.GetPostByNum(boardCode, thread.PostNum)
		if err != nil {
			return nil, err
		}
		stats, err := s.GetThreadStats(&thread)
		if err != nil {
			return nil, err
		}
		lastPos, err := s.PostService.GetNPostsByThread(thread.ID, -1)
		if err != nil {
			return nil, err
		}

		threadWithPost := model.ThreadForCatalogDTO{
			Thread:   thread,
			Post:     *post,
			Stats:    stats,
			LastPost: lastPos[0],
		}
		res = append(res, threadWithPost)
	}

	return res, nil
}

func (s *ThreadService) GetThreadsForHome(boardCode string, lastNpostsCount int, offset, limit int) ([]model.ThreadForHomeDTO, int64, error) {
	threads, err := s.ListThreadsOfBoardOrderByBump(boardCode, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	totalThreads, err := s.GetThreadCountPerBoard(boardCode)
	if err != nil {
		return nil, 0, err
	}

	var res []model.ThreadForHomeDTO

	for _, thread := range threads {
		post, err := s.PostService.GetPostByNum(boardCode, thread.PostNum)
		if err != nil {
			return nil, 0, err
		}
		stats, err := s.GetThreadStats(&thread)
		if err != nil {
			return nil, 0, err
		}
		lastPosts, err := s.PostService.GetNPostsByThread(thread.ID, -lastNpostsCount)
		if err != nil {
			return nil, 0, err
		}
		slices.Reverse(lastPosts)

		// Only add OP if it's not included in the last N posts.
		if len(lastPosts) > 0 && lastPosts[0].ID != post.ID {
			lastPosts = append([]model.Post{*post}, lastPosts...)
		}

		threadWithPosts := model.ThreadForHomeDTO{
			Thread:   thread,
			Posts:    lastPosts,
			Stats:    stats,
		}
		res = append(res, threadWithPosts)
	}

	return res, totalThreads, nil

}

func (s *ThreadService) GetFullThread(threadId uint) (*model.ThreadFullDTO, error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return nil, err
	}

	return s.GetFullThreadFrom(thread)
}

func (s *ThreadService) GetFullThreadByNum(boardCode string, threadNum uint) (*model.ThreadFullDTO, error) {
	thread, err := s.GetThreadByNum(boardCode, threadNum)
	if err != nil {
		return nil, err
	}

	return s.GetFullThreadFrom(thread)
}

func (s *ThreadService) GetThreadStats(thread *model.Thread) (model.ThreadStats, error) {
	return s.ThreadRepo.GetThreadStats(thread.ID)
}

func (s *ThreadService) ArchiveThread(threadId uint) (*model.Thread, error) {
	thread, err := s.GetThread(threadId)
	if err != nil {
		return nil, err
	}

	return s.archiveThread(thread)
}

func (s *ThreadService) archiveThread(thread *model.Thread) (*model.Thread, error) {
	return s.ThreadRepo.ArchiveThread(thread)
}

func (s *ThreadService) ArchiveBumpedOffThreads(board *model.Board) (error) {
	threads, err := s.ListThreadsOfBoardOrderByBump(board.Code, 0, 1000)
	if err != nil {
		return err
	}

	L := min(board.Config.MaxThreads, uint(len(threads)))
	R := uint(len(threads)) - 1

	for i := L; i <= R; i++ {
		_, err := s.archiveThread(&threads[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ThreadService) UpdateAutoCycleForThread(thread *model.Thread) error {
	if thread.AutoCycle == 0 {
		return nil
	}

	stats, err := s.GetThreadStats(thread)
	if err != nil {
		return err
	}

	N := (int(stats.PostCount) - 1) - int(thread.AutoCycle) // -1 because OP is not counted.

	if N > 0 {
		log.Println("deleting ", N, "....................")
		err = s.PostService.DeleteFirstNPostsOfThread(thread, uint(N))
		return err
	} else {
		return nil
	}
}