package repository

import (
	"context"
	"strconv"

	"github.com/jackc/pgx"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/forums/app/models"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewThreadRepo(db *pgx.ConnPool) threadModel.ThreadRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateThread(ctx context.Context, thread *models.Thread) (id int, err error) {
	var query string
	var queryParams []interface{}

	queryParams = append(queryParams,
		thread.Title,
		thread.Author,
		thread.Message,
		thread.Forum,
		thread.Slug,
	)
	if thread.Created != nil {
		query = "INSERT INTO threads (title, user_create, message, forum, slug, created) VALUES ($1, $2, $3, $4, $5, $6) returning id"
		queryParams = append(queryParams, thread.Created)
	} else {
		query = "INSERT INTO threads (title, user_create, message, forum, slug) VALUES ($1, $2, $3, $4, $5) returning id"
	}

	err = r.DB.QueryRow(query, queryParams...).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetThreadBySlugOrId(ctx context.Context, slugOrId string) (*models.Thread, error) {

	thread := new(models.Thread)
	query :=
		`
		SELECT th.id, th.title, th.user_create, th.forum, 
		th.message, th.slug, th.created, th.votes
		FROM threads as th
	`

	if _, err := strconv.Atoi(slugOrId); err == nil {
		query += " WHERE th.id = $1"
	} else {
		query += " WHERE th.slug = $1"
	}

	err := r.DB.QueryRow(query, slugOrId).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Slug,
		&thread.Created,
		&thread.Votes,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (r *repo) UpdateThreadBySlug(ctx context.Context, thread *models.Thread) error {
	query := "UPDATE threads SET title = $1, message = $2 WHERE slug = $3"

	_, err := r.DB.Exec(query, thread.Title, thread.Message, thread.Slug)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateVote(ctx context.Context, vote *models.Vote) error {
	query := "UPDATE votes SET voice = $1 WHERE user_create = $2 AND thread = $3"

	_, err := r.DB.Exec(query, vote.Voice, vote.User, vote.Thread)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) AddVote(ctx context.Context, vote *models.Vote) error {
	id := new(int)

	query := "INSERT INTO votes (user_create, thread, voice) VALUES ($1, $2, $3) returning id"
	err := r.DB.QueryRow(query, vote.User, vote.Thread, vote.Voice).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repo) pathSort(ctx context.Context, threadPosts *models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message,
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.thread = $1
	`

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Desc {
		if threadPosts.Since != "" {
			query += " AND p.path < (SELECT p2.path from posts AS p2 WHERE p2.id = $2)"
			queryParams = append(queryParams, threadPosts.Since)
		}

		query += " ORDER BY p.path DESC"
	} else {
		if threadPosts.Since != "" {
			query += " AND p.path > (SELECT p2.path from posts AS p2 WHERE p2.id = $2)"
			queryParams = append(queryParams, threadPosts.Since)
		}
		query += " ORDER BY p.path"
	}

	if threadPosts.Limit != "" {
		query += " LIMIT " + threadPosts.Limit
	}

	return query, queryParams
}

const selectParentPathLimitAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.root_id IN (
		SELECT p2.root_id
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent is NULL
		ORDER BY p2.id
		LIMIT $2
	)
	ORDER BY p.path
`

const selectParentPathLimitDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.root_id IN (
		SELECT p2.root_id
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent is NULL
		ORDER BY p2.id DESC
		LIMIT $2
	)
	ORDER BY p.root_id DESC, p.path ASC
`

const selectParentPathSinceLimitAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.root_id IN (
		SELECT p2.root_id
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent is NULL AND p2.root_id > (SELECT p3.root_id from posts p3 where p3.id = $2)
		ORDER BY p2.id
		LIMIT $3
	)
	ORDER BY p.path
`

const selectParentPathSinceLimitDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.root_id IN (
		SELECT p2.root_id
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent is NULL AND p2.root_id < (SELECT p3.root_id from posts p3 where p3.id = $2)
		ORDER BY p2.id DESC
		LIMIT $3
	)
	ORDER BY p.root_id DESC, p.path ASC
`

const selectParentPathAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1
	ORDER BY p.path
`

const selectParentPathDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1
	ORDER BY p.root_id DESC, p.path ASC
`

const selectParentPathSinceAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.root_id > (SELECT p3.root_id from posts p3 where p3.id = $2)
	ORDER BY p.path
`

const selectParentPathSinceDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.root_id < (SELECT p3.root_id from posts p3 where p3.id = $2)
	ORDER BY p.root_id DESC, p.path ASC
`

func (r *repo) parentPathSort(ctx context.Context, threadPosts *models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	var query string

	queryParams = append(queryParams, threadPosts.ThreadId)

	switch threadPosts.Limit {
	case "":
		if threadPosts.Desc {
			if threadPosts.Since != "" {
				query = selectParentPathSinceDesc
				queryParams = append(queryParams, threadPosts.Since)
			} else {
				query = selectParentPathDesc
			}
		} else {
			if threadPosts.Since != "" {
				query = selectParentPathSinceAsc
				queryParams = append(queryParams, threadPosts.Since)
			} else {
				query = selectParentPathAsc
			}
		}
	default:
		if threadPosts.Desc {
			if threadPosts.Since != "" {
				query = selectParentPathSinceLimitDesc
				queryParams = append(queryParams, threadPosts.Since)
			} else {
				query = selectParentPathLimitDesc
			}
		} else {
			if threadPosts.Since != "" {
				query = selectParentPathSinceLimitAsc
				queryParams = append(queryParams, threadPosts.Since)
			} else {
				query = selectParentPathLimitAsc
			}
		}
		queryParams = append(queryParams, threadPosts.Limit)
	}

	return query, queryParams
}

func (r *repo) flatSort(ctx context.Context, threadPosts *models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message,
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.thread = $1
	`

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Desc {
		if threadPosts.Since != "" {
			query += " AND p.id < $2"
			queryParams = append(queryParams, threadPosts.Since)
		}

		query += " ORDER BY p.id DESC"
	} else {
		if threadPosts.Since != "" {
			query += " AND p.id > $2"
			queryParams = append(queryParams, threadPosts.Since)
		}
		query += " ORDER BY p.id"
	}

	if threadPosts.Limit != "" {
		query += " LIMIT " + threadPosts.Limit
	}

	return query, queryParams
}

func (r *repo) GetPosts(ctx context.Context, threadPosts *models.ThreadPosts) (*[]models.Post, error) {
	queryParams := make([]interface{}, 0)
	var query string

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Sort == "" {
		threadPosts.Sort = "flat"
	}
	switch threadPosts.Sort {
	case "tree":
		query, queryParams = r.pathSort(ctx, threadPosts)

	case "parent_tree":
		query, queryParams = r.parentPathSort(ctx, threadPosts)

	case "flat":
		query, queryParams = r.flatSort(ctx, threadPosts)
	}

	threadsDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0)
	for threadsDB.Next() {
		post := new(models.Post)
		err := threadsDB.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}

	return &posts, nil
}
