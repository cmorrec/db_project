package delivery

import (
	"encoding/json"
	"fmt"
	"forums/internal/models"
	threadModel "forums/internal/thread"
	"forums/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	threadUcase threadModel.ThreadUsecase
}

func (h Handler) Vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slugOrId"]
	vote := new(models.Vote)
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()

	threadVote, err := h.threadUcase.Vote(*vote, slugOrId)
	if err != nil {
		switch err.Error() {
		case "404":
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	utils.NewResponse(http.StatusOK, threadVote).SendSuccess(w)
}

func (h Handler) AddPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slugOrId"]
	posts := new(models.Posts)
	err := json.NewDecoder(r.Body).Decode(&posts.Posts)
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	if len(posts.Posts) == 0 {
		utils.NewResponse(http.StatusCreated, posts.Posts).SendSuccess(w)
		return
	}
	responsePosts, err := h.threadUcase.AddPosts(*posts, slugOrId)

	if err != nil {
		switch err.Error() {
		case "404":
			w.WriteHeader(http.StatusNotFound)
			return
		case "409":
			utils.NewResponse(http.StatusConflict, responsePosts.Posts).SendSuccess(w)
			return
		}
	}

	utils.NewResponse(http.StatusCreated, responsePosts.Posts).SendSuccess(w)
}

func (h Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrID := vars["slugOrId"]
	defer r.Body.Close()

	thread, err := h.threadUcase.GetThread(slugOrID)
	if err != nil {
		message := models.Error{
			Message: fmt.Sprintf("Can't find thread with slug %s\n", slugOrID),
		}
		utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	utils.NewResponse(http.StatusOK, thread).SendSuccess(w)
}

func NewThreadHandler(threadUcase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	handler := &Handler{
		threadUcase: threadUcase,
	}

	return handler
}
