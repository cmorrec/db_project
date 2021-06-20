package thread

import (
	"forums/internal/models"
	"net/http"
)

type ThreadHandler interface {
	AddPosts(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
}

type ThreadUsecase interface {
	AddPosts(posts models.Posts, slugOrId string) (models.Posts, error)
	Vote(vote models.Vote, slugOrId string) (models.Thread, error)
}

type ThreadRepo interface {
	AddPostsInThreadByID(posts models.Posts, threadID int32, forumSlug string) (models.Posts, error)
	VoteInThreadByID(threadID int32, nickname string, voice int32) error
	UpdateVoteInThreadByID(threadID int32, nickname string, voice int32) error
	UpdateThreadVotes(threadID int32, votes int32) error
	GetVoteInThread(threadID int32, nickname string) (int32, bool, error)
	GetUserByNickName(nickname string) (models.User, error)
	GetThreadByID(id int32) (models.Thread, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	GetPostByID(id int64) (models.Post, error)
	GetNumOfCorrectPosts(ids []int64) (int, error)
}
