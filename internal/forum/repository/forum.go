package repository

import (
	"database/sql"
	"fmt"
	"forums/internal/forum"
	"forums/internal/models"
	"forums/internal/user"
	"strings"
)

type forumRepo struct {
	DB       *sql.DB
	userRepo user.UserRepo
}

func NewForumRepo(db *sql.DB, userRepo_ user.UserRepo) forum.ForumRepo {
	return &forumRepo{
		DB:       db,
		userRepo: userRepo_,
	}
}

func (u forumRepo) CreateForum(newForum models.Forum) (models.Forum, error) {
	query :=
		`
	INSERT INTO forums (title, userNickname, slug)
	VALUES ($1, $2, $3)
	`

	_, err := u.DB.Exec(query,
		newForum.Title,
		newForum.User,
		newForum.Slug)
	if err != nil {
		return models.Forum{}, err
	}

	return newForum, nil
}

func (u forumRepo) GetForumBySlug(slug string) (models.Forum, error) {
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
	if !strings.EqualFold(forum_.Slug, slug) {
		return models.Forum{}, err
	}
	return *forum_, nil
}

func (u forumRepo) GetUserByNickName(nickname string) (models.User, error) {
	return u.userRepo.GetByNickName(nickname)
}

func (u forumRepo) CreateThread(newThread models.Thread) (models.Thread, error) {
	var id int32
	var query string
	var queryParams []interface{}

	queryParams = append(queryParams,
		newThread.Title,
		newThread.Author,
		newThread.Forum,
		newThread.Message,
		newThread.Slug,
	)
	if newThread.Created != "" {
		query =
			`
		INSERT INTO threads (title, author, forum, message, slug, created)
		VALUES ($1, $2, $3, $4, $5, $6) returning id
	`
		queryParams = append(queryParams, newThread.Created)
	} else {
		query =
			`
		INSERT INTO threads (title, author, forum, message, slug)
		VALUES ($1, $2, $3, $4, $5) returning id
	`
	}
	err := u.DB.QueryRow(query, queryParams...).Scan(&id)
	if err != nil {
		return models.Thread{}, err
	}
	newThread.Id = id

	return newThread, nil
}

func (u forumRepo) GetThreadBySlug(slug string) (models.Thread, error) {
	query :=
		`
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug=$1
	`
	thread_ := new(models.Thread)
	err := u.DB.QueryRow(query, slug).Scan(
		&thread_.Id,
		&thread_.Title,
		&thread_.Author,
		&thread_.Forum,
		&thread_.Message,
		&thread_.Votes,
		&thread_.Slug,
		&thread_.Created,
	)
	if !strings.EqualFold(thread_.Slug, slug) {
		return models.Thread{}, err
	}
	return *thread_, nil
}

func getThreadsInForumQuery(forumSlug string, limit int32, since string, desc bool) string {
	query := fmt.Sprintf(`SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE forum='%s'`, forumSlug)

	if since != "" {
		if desc {
			query += fmt.Sprintf(` AND created <= '%s'`, since)
		} else {
			query += fmt.Sprintf(` AND created >= '%s'`, since)
		}
	}

	if desc {
		query += `
	ORDER BY created DESC
	`
	} else {
		query += `
	ORDER BY created ASC
	`
	}

	query += fmt.Sprintf(`
	LIMIT %d`, limit)

	return query
}

func (u forumRepo) GetThreadsInForum(forumSlug string, limit int32, since string, desc bool) ([]models.Thread, error) {
	threadsList := make([]models.Thread, 0)
	query := getThreadsInForumQuery(forumSlug, limit, since, desc)
	threadsDB, err := u.DB.Query(query)
	if err != nil {
		return []models.Thread{}, err
	}
	for threadsDB.Next() {
		thread := new(models.Thread)
		err = threadsDB.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			return []models.Thread{}, err
		}
		threadsList = append(threadsList, *thread)
	}
	return threadsList, nil
}
