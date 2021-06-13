package repository

import (
	"database/sql"
	"forums/internal/forum"
	"forums/internal/models"
	"forums/internal/thread"
)

type threadRepo struct {
	DB *sql.DB
	forumRepo forum.ForumRepo
}

func (r threadRepo) GetThreadBySlug(slug string) (models.Thread, error) {
	return r.forumRepo.GetThreadBySlug(slug)
}

func (r threadRepo) GetThreadByID(id int32) (models.Thread, error) {
	query :=
		`
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE id=$1
	`
	thread_ := new(models.Thread)
	err := r.DB.QueryRow(query, id).Scan(
		&thread_.Id,
		&thread_.Title,
		&thread_.Author,
		&thread_.Forum,
		&thread_.Message,
		&thread_.Votes,
		&thread_.Slug,
		&thread_.Created,
	)
	if thread_.Id != id {
		return models.Thread{}, err
	}
	return *thread_, nil
}

func (r threadRepo) GetPostByID(id int64) (models.Post, error) {
	panic("implement me")
}

func (r threadRepo) AddPostsInThreadByID(posts models.Posts, id int32) (models.Posts, error) {
	panic("implement me")
}

func (r threadRepo) AddPostsInThreadBySlug(posts models.Posts, slug string) (models.Posts, error) {
	panic("implement me")
}

func NewThreadRepo(db *sql.DB, repo forum.ForumRepo) thread.ThreadRepo {
	return &threadRepo{
		DB: db,
		forumRepo: repo,
	}
}
