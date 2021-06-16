package forum

import (
	"forums/internal/models"
	"net/http"
)

type ForumHandler interface {
	CreateForum(w http.ResponseWriter, r *http.Request)
	GetForumBySlug(w http.ResponseWriter, r *http.Request)
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetThreadsInForum(w http.ResponseWriter, r *http.Request)
}

type ForumUsecase interface {
	CreateForum(forum models.Forum) (*models.Forum, error)
	GetForumBySlug(slug string) (*models.Forum, error)
	CreateThread(thread models.Thread, forumSlug string) (*models.Thread, error)
	GetThreadsInForum(forumSlug string, limit int32, since string, desc bool) ([]models.Thread, error)
}

type ForumRepo interface {
	CreateForum(newForum models.Forum) (models.Forum, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetUserByNickName(nickname string) (models.User, error)
	CreateThread(newThread models.Thread) (models.Thread, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	GetThreadsInForum(forumSlug string, limit int32, since string, desc bool) ([]models.Thread, error)
}
