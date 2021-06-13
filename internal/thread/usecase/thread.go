package usecase

import (
	"forums/internal/models"
	"forums/internal/thread"
)

type threadUsecase struct {
	 threadRepository thread.ThreadRepo
}

func (u threadUsecase) AddPosts(posts models.Posts, slugOrId string) (models.Posts, error) {
	panic("implement me")
}

func NewThreadUsecase(repo thread.ThreadRepo) thread.ThreadUsecase {
	return &threadUsecase{
		threadRepository: repo,
	}
}
