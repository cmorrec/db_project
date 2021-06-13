package delivery

import (
	threadModel "forums/internal/thread"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ThreadUcase threadModel.ThreadUsecase
}

func (h Handler) AddPosts(c echo.Context) error {
	panic("implement me")
}

func NewThreadHandler(threadUcase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	handler := &Handler{
		ThreadUcase: threadUcase,
	}

	return handler
}
