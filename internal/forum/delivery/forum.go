package delivery

import (
	"encoding/json"
	"fmt"
	forumModel "forums/internal/forum"
	"forums/internal/models"
	"forums/utils"
	"github.com/gorilla/mux"
	"net/http"
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

func (h Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	newForum := new(models.Forum)
	err := json.NewDecoder(r.Body).Decode(&newForum)
	defer r.Body.Close()
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}

	responseForum, err := h.ForumUcase.CreateForum(*newForum)
	if err != nil {
		switch err.Error() {
		case "404":
			message := models.Error{
				Message: fmt.Sprintf("Can't find user with nickname \n"),
			}
			utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
			return
		case "409":
			utils.NewResponse(http.StatusConflict, responseForum).SendSuccess(w)
			return
		}
	}

	utils.NewResponse(http.StatusCreated, responseForum).SendSuccess(w)
}

func (h Handler) GetForumBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	defer r.Body.Close()

	forum, err := h.ForumUcase.GetForumBySlug(slug)
	if err != nil {
		message := models.Error{
			Message: fmt.Sprintf("Can't find forum with slug %s\n", slug),
		}
		utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	utils.NewResponse(http.StatusOK, forum).SendSuccess(w)
}

func (h Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	forumSlug := vars["slug"]
	newThread := new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(&newThread)
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()

	responseThread, err := h.ForumUcase.CreateThread(*newThread, forumSlug)
	if err != nil {
		switch err.Error() {
		case "404":
			message := models.Error{
				Message: fmt.Sprintf("Can't find forum with slug %s\n", forumSlug),
			}
			utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
			return
		case "409":
			utils.NewResponse(http.StatusConflict, responseThread).SendSuccess(w)
			return
		}
	}

	utils.NewResponse(http.StatusCreated, responseThread).SendSuccess(w)
}

func (h Handler) GetThreadsInForum(w http.ResponseWriter, r *http.Request) {
	var limit int
	var since string
	var desc bool
	var err error


	vars := mux.Vars(r)
	defer r.Body.Close()

	limitQuery :=  r.URL.Query().Get("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			limit = 100
		}
	} else {
		limit = 100
	}

	since = r.URL.Query().Get("since")

	descQuery := r.URL.Query().Get("desc")
	if descQuery != "" {
		desc, err = strconv.ParseBool(descQuery)
		if err != nil {
			desc = false
		}
	} else {
		desc = false
	}

	slug := vars["slug"]
	threads, err := h.ForumUcase.GetThreadsInForum(slug, int32(limit), since, desc)
	if err != nil {
		message := models.Error{
			Message: fmt.Sprintf("Can't find forum with slug %s\n", slug),
		}
		utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	utils.NewResponse(http.StatusOK, threads).SendSuccess(w)
}
