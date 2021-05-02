package models

type PostFull struct {
	Post   Post
	Author User
	Thread Thread
	Forum  Forum
}
