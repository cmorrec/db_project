package user

import (
	"forums/internal/models"
	"net/http"
)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetUserData(w http.ResponseWriter, r *http.Request)
	UpdateUserData(w http.ResponseWriter, r *http.Request)
}

type UserUsecase interface {
	Create(user models.User) ([]models.User, error)
	GetByNickName(nickname string) (*models.User, error)
	UpdateUserData(user models.User) (*models.User, error)
}

type UserRepo interface {
	Create(newUser models.User) (models.User, error)
	GetByNickName(nickname string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	UpdateUserData(updateUser models.User) (models.User, error)
}
