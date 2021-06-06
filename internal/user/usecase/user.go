package usecase

import (
	"forums/internal/models"
	"forums/internal/user"
)

type userUsecase struct {
	userRepository user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) user.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u userUsecase) Create(user_ models.User) (*models.User, error) {
	newUser, err := u.userRepository.Create(user_)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (u userUsecase) GetByNickName(nickname string) (*models.User, error) {
	user_, err := u.userRepository.GetByNickName(nickname)
	if err != nil {
		return nil, err
	}

	return &user_, nil
}
