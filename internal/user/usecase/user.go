package usecase

import (
	"fmt"
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

func (u userUsecase) Create(user_ models.User) ([]models.User, error) {
	// 1 check that not 409
	userErrors := make([]models.User, 0)
	userNickname, errNickname := u.userRepository.GetByNickName(user_.Nickname)
	if errNickname == nil {
		userErrors = append(userErrors, userNickname)
	}
	userEmail, errEmail := u.userRepository.GetByEmail(user_.Email)
	if errEmail == nil {
		if len(userErrors) == 0 {
			userErrors = append(userErrors, userEmail)
		} else if userErrors[0].Nickname != userEmail.Nickname {
			userErrors = append(userErrors, userEmail)
		}
	}
	if len(userErrors) > 0 {
		return userErrors, &utils.CustomError{"409"}
	}

	newUser, err := u.userRepository.Create(user_)
	if err != nil {
		return []models.User{newUser}, err
	}

	return []models.User{newUser}, nil
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
	soughtUser, err := u.userRepository.GetByNickName(user_.Nickname)
	if err != nil {
		return nil, &utils.CustomError{"404"}
	}
	// 2 check that not 409
	sameEmailUser, err := u.userRepository.GetByEmail(user_.Email)
	if err == nil && sameEmailUser.Email != "" && sameEmailUser.Nickname != user_.Nickname {
		fmt.Println("update usecase ", sameEmailUser, user_)
		return nil, &utils.CustomError{"409"}
	}

	fixData(&soughtUser, &user_)
	newUser, err := u.userRepository.UpdateUserData(user_)
	if err != nil {
		return &newUser, err
	}

	return &newUser, nil
}

func fixData(soughtUser *models.User, newUser *models.User) {
	if newUser.Email == "" {
		newUser.Email = soughtUser.Email
	}
	if newUser.About == "" {
		newUser.About = soughtUser.About
	}
	if newUser.Fullname == "" {
		newUser.Fullname = soughtUser.Fullname
	}
}
