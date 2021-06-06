package repository

import (
	"database/sql"
	"forums/internal/models"
	"forums/internal/user"
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
	// TODO check user with this data

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
	DBuser, err := u.DB.Query(query, nickname)
	if err != nil {
		// TODO return real error. NOT this
		return models.User{}, nil
	}
	user_ := new(models.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user_.Nickname,
			&user_.Fullname,
			&user_.About,
			&user_.Email,
		)
		if err != nil {
			// TODO return real error. NOT this
			return models.User{}, nil
		}
	}
	return *user_, nil
}
