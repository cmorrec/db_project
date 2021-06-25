package repository

import (
	"context"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewUserRepo(db *pgx.ConnPool) userModel.UserRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) GetUserByNameAndEmail(ctx context.Context, name, email string) (*[]models.User, error) {
	query := "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2"

	usersDB, err := r.DB.Query(query, name, email)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	for usersDB.Next() {
		user := new(models.User)

		err := usersDB.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return &users, nil
}

func (r *repo) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	user := new(models.User)
	query := "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1"

	err := r.DB.QueryRow(query, name).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	query := "SELECT nickname, fullname, about, email FROM users WHERE email = $1"

	err := r.DB.QueryRow(query, email).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) CreateUser(ctx context.Context, user *models.User) (err error) {

	query := "INSERT INTO users (nickname, fullname, about, email)  VALUES ($1, $2, $3, $4)"

	_, err = r.DB.Exec(query,
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateUser(ctx context.Context, user *models.User) (id int, err error) {
	query := "UPDATE users SET fullname = $1, about = $2, email = $3 WHERE nickname = $4"

	_, err = r.DB.Exec(query, user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		return 0, err
	}

	return id, nil
}
