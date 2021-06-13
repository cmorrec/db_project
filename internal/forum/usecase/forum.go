package usecase

import (
	"forums/internal/forum"
	"forums/internal/models"
	"forums/utils"
)

type forumUsecase struct {
	forumRepository forum.ForumRepo
}

func NewForumUsecase(repo forum.ForumRepo) forum.ForumUsecase {
	return &forumUsecase{
		forumRepository: repo,
	}
}

func (u forumUsecase) CreateForum(forum_ models.Forum) (*models.Forum, error) {
	// 1 check that not 404
	user, err := u.forumRepository.GetUserByNickName(forum_.User)
	if err != nil || user.Nickname == "" {
		return nil, &utils.CustomError{"404"}
	}
	// 2 check that not 409
	sameSlugForum, err := u.forumRepository.GetForumBySlug(forum_.Slug)
	if err == nil && sameSlugForum.Slug != "" {
		return &sameSlugForum, &utils.CustomError{"409"}
	}

	forum_.User = user.Nickname
	newForum, err := u.forumRepository.CreateForum(forum_)
	if err != nil {
		return &newForum, err
	}

	return &newForum, nil
}

func (u forumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	forum_, err := u.forumRepository.GetForumBySlug(slug)
	if err != nil || forum_.Slug == "" {
		return nil, &utils.CustomError{"404"}
	}

	return &forum_, nil
}

func (u forumUsecase) CreateThread(thread models.Thread, forumSlug string) (*models.Thread, error) {
	// 1 check that not 404
	user, err := u.forumRepository.GetUserByNickName(thread.Author)
	if err != nil || user.Nickname == "" {
		return nil, &utils.CustomError{"404"}
	}
	// 2 check that not 409
	if thread.Slug != "" {
		sameSlugThread, err := u.forumRepository.GetThreadBySlug(thread.Slug)
		if err == nil && sameSlugThread.Title != "" {
			return &sameSlugThread, &utils.CustomError{"409"}
		}
	}

	thread.Author = user.Nickname
	newForum, err := u.forumRepository.CreateThread(thread)
	if err != nil {
		return &newForum, err
	}

	return &newForum, nil
}
