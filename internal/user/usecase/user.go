package usecase

import (
	"forums/internal/models"
	"forums/internal/user"
	"forums/utils"
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
		return &newUser, err
	}

	return &newUser, nil
}

func (u userUsecase) GetByNickName(nickname string) (*models.User, error) {
	user_, err := u.userRepository.GetByNickName(nickname)
	if err != nil || user_.Nickname == "" {
		return nil, &utils.CustomError{"404"}
	}

	return &user_, nil
}

func (u userUsecase) UpdateUserData(user_ models.User) (*models.User, error) {
	// 1 check that not 404
	_, err := u.userRepository.GetByNickName(user_.Nickname)
	if err != nil {
		return nil, &utils.CustomError{"404"}
	}
	// 2 check that not 409
	sameEmailUser, err := u.userRepository.GetByEmail(user_.Email)
	if err == nil && sameEmailUser.Nickname != user_.Nickname {
		return nil, &utils.CustomError{"409"}
	}

	newUser, err := u.userRepository.UpdateUserData(user_)
	if err != nil {
		return &newUser, err
	}

	return &newUser, nil
}
