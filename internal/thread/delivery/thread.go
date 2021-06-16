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

func (h Handler) AddPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slugOrId"]
	posts := new(models.Posts)
	err := json.NewDecoder(r.Body).Decode(&posts.Posts)
	if err != nil {
		fmt.Println(1)
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	fmt.Println(2)
	responsePosts, err := h.threadUcase.AddPosts(*posts, slugOrId)
	fmt.Println(3)
	if err != nil {
		fmt.Println(4)
		switch err.Error() {
		case "404":
			fmt.Println(5)
			w.WriteHeader(http.StatusNotFound)
			return
		case "409":
			fmt.Println(6)
			utils.NewResponse(http.StatusConflict, responsePosts).SendSuccess(w)
			return
		}
	}
	fmt.Println(7)
	utils.NewResponse(http.StatusOK, responsePosts)
	return
}

func NewThreadHandler(threadUcase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	handler := &Handler{
		threadUcase: threadUcase,
	}

	return handler
}
