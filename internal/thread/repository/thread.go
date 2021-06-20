package repository

import (
	"database/sql"
	"fmt"
	"forums/internal/forum"
	"forums/internal/models"
	"forums/internal/thread"
)

type threadRepo struct {
	DB        *sql.DB
	forumRepo forum.ForumRepo
}

func (r threadRepo) VoteInThreadByID(threadID int32, nickname string, voice int32) error {
	query := `
				INSERT INTO votes (author, thread, voice) values ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, nickname, threadID, voice)
	if err != nil {
		return err
	}

	return nil
}

func (r threadRepo) UpdateThreadVotes(threadID int32, votes int32) error {
	query := `
				UPDATE threads SET votes=$2 WHERE id=$1
	`
	_, err := r.DB.Exec(query, threadID, votes)
	if err != nil {
		return err
	}

	return nil
}

func (r threadRepo) UpdateVoteInThreadByID(threadID int32, nickname string, voice int32) error {
	query := `
				UPDATE votes SET voice=$3 WHERE author=$1 AND thread=$2
	`
	_, err := r.DB.Exec(query, nickname, threadID, voice)
	if err != nil {
		return err
	}

	return nil
}

func (r threadRepo) GetVoteInThread(threadID int32, nickname string) (int32, bool, error) {
	var voice int32
	var err error
	query :=
		`
	SELECT voice
	FROM votes
	WHERE thread=$1 and author=$2
	`
	err = r.DB.QueryRow(query, threadID, nickname).Scan(&voice)
	if err != nil {
		return voice, false, err
	}

	return voice, true, nil
}

func (r threadRepo) GetThreadBySlug(slug string) (models.Thread, error) {
	return r.forumRepo.GetThreadBySlug(slug)
}

func (r threadRepo) GetUserByNickName(nickname string) (models.User, error) {
	return r.forumRepo.GetUserByNickName(nickname)
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
		return models.Posts{}, nil
	}
	var created = ""
	//authors := make([]string, 0)
	//parents := make([]string, 0)
	//messages := make([]string, 0)
	for _, post := range posts.Posts {
		//authors = append(authors, post.Author)
		//parents = append(parents, strconv.FormatInt(post.Parent, 10))
		//messages = append(messages, post.Message)
		if created == "" && post.Created != "" {
			created = post.Created
		}
	}
	//authorsString := strings.Join(authors, "', '")
	//parentsString := strings.Join(parents, ", ")
	//messagesString := strings.Join(messages, "', '")

	var query string
	//if created != "" {
	//	query = fmt.Sprintf(`
	//insert into posts (thread, forum, created, author, parent, message)
	//select %d, '%s', '%s', unnest(array['%s'], array[%s], array['%s']) returning id
	//`, threadId, forumSlug, created, authorsString, parentsString, messagesString)
	//} else {
	//	query = fmt.Sprintf(`
	//insert into posts (thread, forum, author, parent, message)
	//select %d, '%s', unnest(array['%s'], array[%s], array['%s']) returning id
	//`, threadId, forumSlug, authorsString, parentsString, messagesString)
	//}
	if created != "" {
		query = `
					insert into posts (thread, forum, created, author, parent, message) values`
		for index, post := range posts.Posts {
			query += fmt.Sprintf(`
					(%d, '%s', '%s', '%s', %d, '%s')`,
				threadId, forumSlug, post.Created, post.Author, post.Parent, post.Message)
			if index != len(posts.Posts)-1 {
				query += ", "
			}
		}
		query += ` returning id, created`
	} else {
		query = `
					insert into posts (thread, forum, author, parent, message) values`
		for index, post := range posts.Posts {
			query += fmt.Sprintf(`
					(%d, '%s', '%s', %d, '%s')`,
				threadId, forumSlug, post.Author, post.Parent, post.Message)
			if index != len(posts.Posts)-1 {
				query += ", "
			}
		}
		query += ` returning id`
	}

	postsDB, err := r.DB.Query(query)
	if err != nil {
		return models.Posts{}, nil
	}
	i := 0
	for postsDB.Next() {
		if created != "" {
			err = postsDB.Scan(&posts.Posts[i].Id, &posts.Posts[i].Created)
		} else {
			err = postsDB.Scan(&posts.Posts[i].Id)
		}
		if err != nil {
			return models.Posts{}, err
		}
		posts.Posts[i].Forum = forumSlug
		posts.Posts[i].Thread = threadId
		i++
	}

	return posts, nil
}

func NewThreadRepo(db *sql.DB, repo forum.ForumRepo) thread.ThreadRepo {
	return &threadRepo{
		DB:        db,
		forumRepo: repo,
	}
}
