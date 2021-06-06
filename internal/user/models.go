package user

import (
	"forums/internal/models"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Create(c echo.Context) error
	GetUserData(c echo.Context) error
}

type UserUsecase interface {
	Create(user models.User) (*models.User, error)
	GetByNickName(nickname string) (*models.User, error)
}

type UserRepo interface {
	Create(newUser models.User) (models.User, error)
	GetByNickName(nickname string) (models.User, error)
}
