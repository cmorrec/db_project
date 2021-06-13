package thread

import (
	"forums/internal/models"
	"github.com/labstack/echo/v4"
)

type ThreadHandler interface {
	AddPosts(c echo.Context) error
}

type ThreadUsecase interface {
	AddPosts(posts models.Posts, slugOrId string) (models.Posts, error)
}

type ThreadRepo interface {
	AddPostsByID(posts models.Posts, id int) (models.Posts, error)
	AddPostsBySlug(posts models.Posts, slug string) (models.Posts, error)
}
