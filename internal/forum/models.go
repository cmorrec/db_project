package forum

import (
	"forums/internal/models"
	"github.com/labstack/echo/v4"
)

type ForumHandler interface {
	CreateForum(c echo.Context) error
	GetForumBySlug(c echo.Context) error
	CreateThread(c echo.Context) error
	GetThreadsInForum(c echo.Context) error
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
