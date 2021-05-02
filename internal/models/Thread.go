package models

type Thread struct {
	Id      int32  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Slug    string `json:"slug"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int32  `json:"votes"`
	Created string `json:"created"`
}
