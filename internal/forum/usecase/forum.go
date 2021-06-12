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
	if user.Nickname == "" || err != nil {
		return nil, &utils.CustomError{"404"}
	}
	// 2 check that not 409
	sameSlugForum, err := u.forumRepository.GetBySlug(forum_.Slug)
	if err == nil && sameSlugForum.Slug != "" {
		return &sameSlugForum, &utils.CustomError{"409"}
	}

	newForum, err := u.forumRepository.CreateForum(forum_)
	if err != nil {
		return &newForum, err
	}

	return &newForum, nil
}
