package delivery

import (
	forumModel "forums/internal/forum"
	"forums/internal/models"
	"github.com/labstack/echo/v4"
	"strconv"
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

func (h Handler) CreateThread(c echo.Context) error {
	newThread := new(models.Thread)
	if err := c.Bind(newThread); err != nil {
		// TODO error
		return nil
	}
	forumSlug := c.Param("slug")

	responseThread, err := h.ForumUcase.CreateThread(*newThread, forumSlug)
	if err != nil {
		switch err.Error() {
		case "404":
			return models.SendResponseWithErrorNotFound(c)
		case "409":
			return models.SendResponseWithErrorConflict(c, responseThread)
		}
	}

	return models.SendResponseCreate(c, responseThread)
}

func (h Handler) GetThreadsInForum(c echo.Context) error {
	var limit int
	var since string
	var desc bool
	var err error

	if c.QueryParam("limit") != "" {
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			limit = 100
		}
	} else {
		limit = 100
	}

	since = c.QueryParam("since")

	if c.QueryParam("desc") != "" {
		desc, err = strconv.ParseBool(c.QueryParam("desc"))
		if err != nil {
			desc = false
		}
	} else {
		desc = false
	}

	threads, err := h.ForumUcase.GetThreadsInForum(c.Param("slug"), int32(limit), since, desc)
	if err != nil {
		return models.SendResponseWithErrorNotFound(c)
	}

	return models.SendResponse(c, threads)
}
