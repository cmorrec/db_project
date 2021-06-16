package repository

import (
	"database/sql"
	"fmt"
	"forums/internal/forum"
	"forums/internal/models"
	"forums/internal/thread"
	"strconv"
	"strings"
)

type threadRepo struct {
	DB        *sql.DB
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

func (r threadRepo) GetNumOfCorrectPosts(ids []int64) (int, error) {
	//select count(id) from users where id in (1,2,3);
	var count int
	if len(ids) == 0 {
		return 0, nil
	}
	idsArr := ""
	for index, id := range ids {
		if index != len(ids)-1 {
			idsArr += fmt.Sprintf("%d, ", id)
		} else {
			idsArr += fmt.Sprintf("%d", id)
		}
	}
	query := fmt.Sprintf("select count(id) from users where id in (%s)", idsArr)
	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r threadRepo) GetPostByID(id int64) (models.Post, error) {
	panic("implement me")
}

func (r threadRepo) AddPostsInThreadByID(posts models.Posts, threadId int32, forumSlug string) (models.Posts, error) {
	// insert into posts (thread, forum, author, parent, created, message, isEdited)
	// select 3, '', unnest([], [], [], [], [])
	if len(posts.Posts) == 0 {
		fmt.Println(15)
		return models.Posts{}, nil
	}
	created := posts.Posts[0].Created
	authors := make([]string, 0)
	parents := make([]string, 0)
	messages := make([]string, 0)
	fmt.Println(16)
	for _, post := range posts.Posts {
		authors = append(authors, post.Author)
		parents = append(parents, strconv.FormatInt(post.Parent, 10))
		messages = append(messages, post.Message)
	}
	fmt.Println(17)
	authorsString := strings.Join(authors, ", ")
	parentsString := strings.Join(parents, ", ")
	messagesString := strings.Join(messages, ", ")

	query := fmt.Sprintf(`
	insert into posts (thread, forum, created, author, parent, message)
	select %d, '%s', %s, unnest([%s], [%s], [%s]) returning id
	`, threadId, forumSlug, created, authorsString, parentsString, messagesString)
	fmt.Println(18, query)
	postsDB, err := r.DB.Query(query)
	if err != nil {
		fmt.Println(19, err)
		return models.Posts{}, nil
	}
	i := 0
	for postsDB.Next() {
		err := postsDB.Scan(&posts.Posts[i].Id)
		fmt.Println(20, i)
		if err != nil {
			return models.Posts{}, nil
		}
		i++
	}
	fmt.Println(21)
	return posts, nil
}

func NewThreadRepo(db *sql.DB, repo forum.ForumRepo) thread.ThreadRepo {
	return &threadRepo{
		DB:        db,
		forumRepo: repo,
	}
}
