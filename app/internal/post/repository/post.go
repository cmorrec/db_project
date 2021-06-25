package repository

import (
	"context"
	"fmt"

	postModel "github.com/forums/app/internal/post"
	"github.com/forums/app/models"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewPostRepo(db *pgx.ConnPool) postModel.PostRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) GetPostsThread(ctx context.Context, id int) (int, error) {
	query :=
		`
		SELECT thread
		FROM posts
		WHERE id = $1
	`

	var thread int
	err := r.DB.QueryRow(query, id).Scan(
		&thread,
	)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return thread, nil
}

func (r *repo) GetPost(ctx context.Context, id int) (*models.Post, error) {
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message, 
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.id = $1
	`

	post := new(models.Post)
	err := r.DB.QueryRow(query, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *repo) UpdateMessage(ctx context.Context, request *models.MessagePostRequest) error {
	_, err := r.DB.Exec("UPDATE posts SET message = $1, is_edited = true WHERE id = $2", request.Message, request.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) CreatePosts(ctx context.Context, posts *[]models.Post) (*[]models.Post, error) {
	var queryParams []interface{}
	query := "INSERT INTO posts (parent, user_create, message, forum, thread, created) VALUES "

	for i, post := range *posts {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)

		if i != len(*posts)-1 {
			query += ","
		}

		queryParams = append(queryParams, post.Parent, post.Author, post.Message, post.Forum, post.Thread, post.Created)
	}

	query += " returning id, created"

	postsDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}

	i := 0
	for postsDB.Next() {
		err = postsDB.Scan(
			&((*posts)[i].Id),
			&((*posts)[i].Created),
		)

		if err != nil {
			return nil, err
		}
		i++
	}

	if dbErr, ok := postsDB.Err().(pgx.PgError); ok {
		return nil, dbErr
	}

	return posts, nil
}

func (r *repo) CreateForumsUsers(ctx context.Context, posts *[]models.Post) error {
	var queryParams []interface{}
	query := "INSERT INTO forums_users (forum, user_create) VALUES "

	for i, post := range *posts {
		query += fmt.Sprintf("($%d, $%d)",
			i*2+1, i*2+2)

		if i != len(*posts)-1 {
			query += ","
		}

		queryParams = append(queryParams, post.Forum, post.Author)
	}

	query += " ON CONFLICT DO NOTHING"

	_, err := r.DB.Exec(query, queryParams...)
	if err != nil {
		return err
	}

	return nil
}
