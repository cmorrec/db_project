package repository

import (
	"context"

	serviceModel "github.com/forums/app/internal/service"
	"github.com/forums/app/models"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewServiceRepo(db *pgx.ConnPool) serviceModel.ServiceRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) ClearDb(ctx context.Context) error {
	_, err := r.DB.Exec("TRUNCATE users, forums, threads, posts, forums_users, votes CASCADE")
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) StatusDb(ctx context.Context) (*models.InfoStatus, error) {
	info := new(models.InfoStatus)
	var err error

	info.User, err = r.getUsersNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Forum, err = r.getForumsNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Thread, err = r.getThreadsNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Post, err = r.getPostsNumber(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (r *repo) getUsersNumber(ctx context.Context) (number int, err error) {
	err = r.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (r *repo) getForumsNumber(ctx context.Context) (number int, err error) {
	err = r.DB.QueryRow("SELECT COUNT(*) FROM forums").Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (r *repo) getThreadsNumber(ctx context.Context) (number int, err error) {
	err = r.DB.QueryRow("SELECT COUNT(*) FROM threads").Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (r *repo) getPostsNumber(ctx context.Context) (number int, err error) {
	err = r.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}
