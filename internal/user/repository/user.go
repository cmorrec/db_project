package repository

import (
	"database/sql"
	"forums/internal/models"
	"forums/internal/user"
	"strings"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepo {
	return &userRepo{
		DB: db,
	}
}

func (u userRepo) Create(newUser models.User) (models.User, error) {
	query :=
		`
	INSERT INTO users (nickname, fullname, about, email) 
	VALUES ($1, $2, $3, $4)
	`

	u.DB.QueryRow(query, newUser.Nickname, newUser.Fullname, newUser.About, newUser.Email)

	return newUser, nil
}

func (u userRepo) GetByNickName(nickname string) (models.User, error) {
	query :=
		`
	SELECT nickname, fullname, about, email
	FROM users 
	WHERE nickname=$1
	`
	user_ := new(models.User)
	err := u.DB.QueryRow(query, nickname).Scan(
		&user_.Nickname,
		&user_.Fullname,
		&user_.About,
		&user_.Email,
	)
	if !strings.EqualFold(user_.Nickname, nickname) {
		return models.User{}, err
	}
	return *user_, nil
}

func (u userRepo) GetByEmail(email string) (models.User, error) {
	query :=
		`
	SELECT nickname, fullname, about, email
	FROM users 
	WHERE email=$1
	`
	user_ := new(models.User)
	err := u.DB.QueryRow(query, email).Scan(
		&user_.Nickname,
		&user_.Fullname,
		&user_.About,
		&user_.Email,
	)
	if !strings.EqualFold(user_.Email, email) {
		return models.User{}, err
	}
	return *user_, nil
}

func (u userRepo) UpdateUserData(updateUser models.User) (models.User, error) {
	query := `UPDATE users SET fullname = $1, about = $2, email = $3 WHERE nickname = $4`
	_, err := u.DB.Exec(query, updateUser.Fullname, updateUser.About, updateUser.Email, updateUser.Nickname)
	if err != nil {
		return models.User{}, err
	}

	return updateUser, nil
}
