package usecase

import (
	"forums/internal/models"
	"forums/internal/thread"
	"forums/utils"
	"strconv"
	"strings"
)

type threadUsecase struct {
	threadRepository thread.ThreadRepo
}

func (u threadUsecase) GetThread(slugOrId string) (models.Thread, error) {
	var thread_ models.Thread
	var err error
	slug, id, isId := isID(slugOrId)
	if isId {
		thread_, err = u.GetThreadByID(id)
	} else {
		thread_, err = u.GetThreadBySlug(slug)
	}
	if err != nil || thread_.Slug == "" {
		return models.Thread{}, &utils.CustomError{"404"}
	}

	return thread_, nil
}

func (u threadUsecase) Vote(vote models.Vote, slugOrId string) (models.Thread, error) {
	var thread_ models.Thread
	var err error
	// check 404 thread, user
	slug, id, isId := isID(slugOrId)
	if isId {
		thread_, err = u.GetThreadByID(id)
	} else {
		thread_, err = u.GetThreadBySlug(slug)
	}
	if err != nil {
		return models.Thread{}, err
	}
	_, err = u.threadRepository.GetUserByNickName(vote.Nickname)
	if err != nil {
		return models.Thread{}, &utils.CustomError{"404"}
	}

	voice, was, err := u.threadRepository.GetVoteInThread(thread_.Id, vote.Nickname)
	if err == nil && was {
		err = u.threadRepository.UpdateVoteInThreadByID(thread_.Id, vote.Nickname, vote.Voice)
		if err != nil {
			return models.Thread{}, err
		}
		thread_.Votes += vote.Voice - voice
		err = u.threadRepository.UpdateThreadVotes(thread_.Id, thread_.Votes)
		if err != nil {
			return models.Thread{}, err
		}

		return thread_, nil
	}

	err = u.threadRepository.VoteInThreadByID(thread_.Id, vote.Nickname, vote.Voice)
	if err != nil {
		return models.Thread{}, err
	}

	thread_.Votes += vote.Voice
	err = u.threadRepository.UpdateThreadVotes(thread_.Id, thread_.Votes)
	if err != nil {
		return models.Thread{}, err
	}

	return thread_, nil
}

func (u threadUsecase) checkParentPosts(posts models.Posts) (bool, error) {
	soughtIDs := make([]int64, 0)
	for _, post := range posts.Posts {
		haveId := false
		for _, id := range soughtIDs {
			if id == post.Id {
				haveId = true
				break
			}
		}
		if !haveId && post.Id != 0 {
			soughtIDs = append(soughtIDs, post.Id)
		}
	}

	num, err := u.threadRepository.GetNumOfCorrectPosts(soughtIDs)
	if err != nil {
		return false, err
	}
	if num != len(soughtIDs) {
		return false, nil
	}

	return true, nil
}

func (u threadUsecase) GetThreadByID(id int32) (models.Thread, error) {
	thread_, err := u.threadRepository.GetThreadByID(id)
	if err != nil || thread_.Id != id {
		return models.Thread{}, &utils.CustomError{"404"}
	}

	return thread_, nil
}

func (u threadUsecase) GetThreadBySlug(slug string) (models.Thread, error) {
	thread_, err := u.threadRepository.GetThreadBySlug(slug)
	if err != nil || !strings.EqualFold(thread_.Slug, slug) {
		return models.Thread{}, &utils.CustomError{"404"}
	}

	return thread_, nil
}

func (u threadUsecase) AddPosts(posts models.Posts, slugOrId string) (models.Posts, error) {
	var thread_ models.Thread
	var err error
	slug, id, isId := isID(slugOrId)
	if isId {
		thread_, err = u.GetThreadByID(id)
	} else {
		thread_, err = u.GetThreadBySlug(slug)
	}
	if err != nil {
		return models.Posts{}, err
	}

	// 409 posts
	ok, err := u.checkParentPosts(posts)
	if !ok || err != nil {
		return models.Posts{}, &utils.CustomError{"409"}
	}
	posts_, err := u.threadRepository.AddPostsInThreadByID(posts, thread_.Id, thread_.Forum)
	if err != nil {
		return models.Posts{}, err
	}

	return posts_, nil
}

func isID(slugOrId string) (string, int32, bool) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		return slugOrId, -1, false
	}
	return "", int32(id), true
}

func NewThreadUsecase(repo thread.ThreadRepo) thread.ThreadUsecase {
	return &threadUsecase{
		threadRepository: repo,
	}
}
