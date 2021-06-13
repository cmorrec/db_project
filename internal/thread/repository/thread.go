package repository

import (
	"database/sql"
	"forums/internal/models"
	"forums/internal/thread"
)

type threadRepo struct {
	DB *sql.DB
}

func (t threadRepo) AddPostsByID(posts models.Posts, id int) (models.Posts, error) {
	panic("implement me")
}

func (t threadRepo) AddPostsBySlug(posts models.Posts, slug string) (models.Posts, error) {
	panic("implement me")
}

func NewThreadRepo(db *sql.DB) thread.ThreadRepo {
	return &threadRepo{
		DB: db,
	}
}
