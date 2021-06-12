package forum

import (
	"forums/internal/models"
	"github.com/labstack/echo/v4"
)

type ForumHandler interface {
	CreateForum(c echo.Context) error
}

type ForumUsecase interface {
	CreateForum(user models.Forum) (*models.Forum, error)
}

type ForumRepo interface {
	CreateForum(newUser models.Forum) (models.Forum, error)
	GetBySlug(slug string) (models.Forum, error)
	GetUserByNickName(nickname string) (models.User, error)
}
