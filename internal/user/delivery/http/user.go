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
	fmt.Println("create user handler", c.QueryParams(), c.ParamValues(), c.Request().Body)
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		// TODO error
		return nil
	}

	responseUser, err := h.UserUcase.Create(*newUser)
	if err != nil {
		// TODO error
		return nil
	}

	return models.SendResponse(c, responseUser)
}

func (h Handler) GetUserData(c echo.Context) error {
	fmt.Println("get user handler", c.QueryParams(), c.ParamValues(), c.Request().Body)
	return models.SendResponse(c, nil)
}
