package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forums/internal/models"
	userModel "forums/internal/user"
	"forums/utils"
	"github.com/gorilla/mux"
)

type handler struct {
	UserUcase userModel.UserUsecase
}

func NewUserHandler(usecase userModel.UserUsecase) userModel.UserHandler {
	return &handler{
		UserUcase: usecase,
	}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	newUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	defer r.Body.Close()
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}

	newUser.Nickname = nickname

	responseUser, err := h.UserUcase.Create(*newUser)
	if err != nil {
		utils.NewResponse(http.StatusConflict, responseUser).SendSuccess(w)
		return
	}

	utils.NewResponse(http.StatusCreated, responseUser[0]).SendSuccess(w)
}

func (h handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	defer r.Body.Close()
	user, err := h.UserUcase.GetByNickName(nickname)
	if err != nil {
		message := models.Error{
			Message: fmt.Sprintf("Can't find user with nickname %s\n", nickname),
		}
		utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	utils.NewResponse(http.StatusOK, user).SendSuccess(w)
}

func (h handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	updateUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	defer r.Body.Close()
	if err != nil {
		sendErr := utils.NewError(http.StatusBadRequest, err.Error())
		w.WriteHeader(sendErr.Code())
		return
	}

	updateUser.Nickname = nickname
	responseUser, err := h.UserUcase.UpdateUserData(*updateUser)
	if err != nil {
		switch err.Error() {
		case "404":
			message := models.Error{
				Message: fmt.Sprintf("Can't find user with nickname %s\n", nickname),
			}
			utils.NewResponse(http.StatusNotFound, message).SendSuccess(w)
			return
		case "409":
			message := models.Error{
				Message: fmt.Sprintf("Can't update user with nickname %s\n", nickname),
			}
			utils.NewResponse(http.StatusConflict, message).SendSuccess(w)
			return
		}
	}

	utils.NewResponse(http.StatusOK, responseUser).SendSuccess(w)
}
