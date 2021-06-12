package repository

import (
	"database/sql"
	"forums/internal/forum"
	"forums/internal/models"
	"forums/internal/user"
)

type forumRepo struct {
	DB *sql.DB
	userRepo user.UserRepo
}

func NewForumRepo(db *sql.DB, userRepo_ user.UserRepo) forum.ForumRepo {
	return &forumRepo{
		DB: db,
		userRepo: userRepo_,
	}
}

func (u forumRepo) CreateForum(newForum models.Forum) (models.Forum, error) {
	query :=
		`
	INSERT INTO forums (title, userNickname, slug)
	VALUES ($1, $2, $3)
	`

	u.DB.QueryRow(query, newForum.Title, newForum.User, newForum.Slug)

	return newForum, nil
}

func (u forumRepo) GetBySlug(slug string) (models.Forum, error) {
	query :=
		`
	SELECT title, userNickname, slug, posts, threads
	FROM forums 
	WHERE slug=$1
	`
	forum_ := new(models.Forum)
	err := u.DB.QueryRow(query, slug).Scan(
		&forum_.Title,
		&forum_.User,
		&forum_.Slug,
		&forum_.Posts,
		&forum_.Threads,
	)
	if forum_.Slug != slug {
		return models.Forum{}, err
	}
	return *forum_, nil
}

func (u forumRepo) GetUserByNickName(nickname string) (models.User, error) {
	return u.userRepo.GetByNickName(nickname)
}
