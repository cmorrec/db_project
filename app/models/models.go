package models

import "time"

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type ForumUsers struct {
	Slug  string `json:"slug"`
	Limit string `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type ForumThreads struct {
	Slug  string `json:"slug"`
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type Message struct {
	Message string `json:"message"`
}

type Post struct {
	Id       int64     `json:"id"`
	Parent   *int64    `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
}

type RequestPost struct {
	Id      int    `json:"id"`
	Related string `json:"related"`
}

type InfoPost struct {
	Post   *Post   `json:"post"`
	User   *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Thread *Thread `json:"thread"`
}

type MessagePostRequest struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type Nesting struct {
	Parent []int64
	Last   []int64
}

type InfoStatus struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

type Thread struct {
	Id      int        `json:"id"`
	Title   string     `json:"title"`
	Author  string     `json:"author"`
	Forum   string     `json:"forum"`
	Message string     `json:"message"`
	Votes   int        `json:"votes"`
	Slug    string     `json:"slug"`
	Created *time.Time `json:"created"`
}

type ThreadPosts struct {
	SlugOrId string `json:"slug"`
	Limit    string `json:"limit"`
	Since    string `json:"since"`
	Sort     string `json:"sort"`
	Desc     bool   `json:"desc"`
	ThreadId int
}

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Vote struct {
	Id     int    `json:"id"`
	User   string `json:"nickname"`
	Thread int    `json:"thread"`
	Voice  int    `json:"voice"`
}
