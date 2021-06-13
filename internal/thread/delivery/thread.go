package delivery

import (
	"forums/internal/models"
	threadModel "forums/internal/thread"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	threadUcase threadModel.ThreadUsecase
}

func (h Handler) AddPosts(c echo.Context) error {
	slugOrId := c.Param("slugOrId")
	posts := new(models.Posts)
	if err := c.Bind(posts.Posts); err != nil {
		return nil
	}

	responsePosts, err := h.threadUcase.AddPosts(*posts, slugOrId)
	if err != nil {
		switch err.Error() {
		case "404":
			return models.SendResponseWithErrorNotFound(c)
		case "409":
			return models.SendResponseWithErrorConflictMessage(c)
		}
	}

	return models.SendResponse(c, responsePosts)
}

func NewThreadHandler(threadUcase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	handler := &Handler{
		threadUcase: threadUcase,
	}

	return handler
}
