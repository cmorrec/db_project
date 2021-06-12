package delivery

import (
	forumModel "forums/internal/forum"
	"forums/internal/models"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ForumUcase forumModel.ForumUsecase
}

func NewForumHandler(forumUsecase forumModel.ForumUsecase) forumModel.ForumHandler {
	handler := &Handler{
		ForumUcase: forumUsecase,
	}

	return handler
}

func (h Handler) CreateForum(c echo.Context) error {
	newForum := new(models.Forum)
	if err := c.Bind(newForum); err != nil {
		// TODO error
		return nil
	}

	responseForum, err := h.ForumUcase.CreateForum(*newForum)
	if err != nil {
		switch err.Error() {
		case "404":
			return models.SendResponseWithErrorNotFound(c)
		case "409":
			return models.SendResponseWithErrorConflict(c, responseForum)
		}
	}

	return models.SendResponseCreate(c, responseForum)
}

func (h Handler) GetForumBySlug(c echo.Context) error {
	forum, err := h.ForumUcase.GetForumBySlug(c.Param("slug"))
	if err != nil {
		return models.SendResponseWithErrorNotFound(c)
	}

	return models.SendResponse(c, forum)
}
