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

	for _, id := range soughtIDs {
		post, err := u.threadRepository.GetPostByID(id)
		if err != nil || post.Id != id {
			return false, err
		}
	}

	return true, nil
}

func (u threadUsecase) addPostsByID(posts models.Posts, id int32) (models.Posts, error) {
	// 404 thread
	thread_, err := u.threadRepository.GetThreadByID(id)
	if err != nil || thread_.Id != id {
		return models.Posts{}, &utils.CustomError{"404"}
	}
	// 409 posts
	ok, err := u.checkParentPosts(posts)
	if !ok || err != nil {
		return models.Posts{}, &utils.CustomError{"409"}
	}
	posts_, err := u.threadRepository.AddPostsInThreadByID(posts, id)
	if err != nil {
		return models.Posts{}, err
	}

	return posts_, nil
}

func (u threadUsecase) addPostsBySlug(posts models.Posts, slug string) (models.Posts, error) {
	// 404 thread
	thread_, err := u.threadRepository.GetThreadBySlug(slug)
	if err != nil || !strings.EqualFold(thread_.Slug, slug) {
		return models.Posts{}, &utils.CustomError{"404"}
	}
	// 409 posts
	ok, err := u.checkParentPosts(posts)
	if !ok || err != nil {
		return models.Posts{}, &utils.CustomError{"409"}
	}
	posts_, err := u.threadRepository.AddPostsInThreadBySlug(posts, slug)
	if err != nil {
		return models.Posts{}, err
	}

	return posts_, nil
}

func (u threadUsecase) AddPosts(posts models.Posts, slugOrId string) (models.Posts, error) {
	slug, id, isId := isID(slugOrId)
	if isId {
		return u.addPostsByID(posts, id)
	} else {
		return u.addPostsBySlug(posts, slug)
	}
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
