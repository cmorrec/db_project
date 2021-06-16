package thread

import (
	"forums/internal/models"
	"net/http"
)

type ThreadHandler interface {
	AddPosts(w http.ResponseWriter, r *http.Request)
}

type ThreadUsecase interface {
	AddPosts(posts models.Posts, slugOrId string) (models.Posts, error)
}

type ThreadRepo interface {
	AddPostsInThreadByID(posts models.Posts, threadID int32, forumSlug string) (models.Posts, error)
	GetThreadByID(id int32) (models.Thread, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	GetPostByID(id int64) (models.Post, error)
	GetNumOfCorrectPosts(ids []int64) (int, error)
}
