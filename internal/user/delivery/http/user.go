package http

import (
	"fmt"
	"forums/internal/models"
	userModel "forums/internal/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserUcase userModel.UserUsecase
}

func NewUserHandler(userUcase userModel.UserUsecase) userModel.UserHandler {
	handler := &Handler{
		UserUcase: userUcase,
	}

	return handler
}

func (h Handler) Create(c echo.Context) error {
	newUser := new(models.User)
	newUser.Nickname = c.Param("nickname")
	if err := c.Bind(newUser); err != nil {
		// TODO error
		fmt.Println("no bind")
		return nil
	}

	responseUser, err := h.UserUcase.Create(*newUser)
	if err != nil {
		return models.SendResponseWithErrorConflict(c, responseUser)
	}

	return models.SendResponseCreate(c, responseUser[0])
}

func (h Handler) GetUserData(c echo.Context) error {
	user, err := h.UserUcase.GetByNickName(c.Param("nickname"))
	if err != nil {
		return models.SendResponseWithErrorNotFound(c)
	}
	return models.SendResponse(c, user)
}

func (h Handler) UpdateUserData(c echo.Context) error {
	updateUser := new(models.User)
	updateUser.Nickname = c.Param("nickname")
	if err := c.Bind(updateUser); err != nil {
		// TODO error
		return nil
	}

	responseUser, err := h.UserUcase.UpdateUserData(*updateUser)
	if err != nil {
		switch err.Error() {
		case "404":
			return models.SendResponseWithErrorNotFound(c)
		case "409":
			return models.SendResponseWithErrorConflictMessage(c)
		}
	}

	return models.SendResponse(c, responseUser)
}
