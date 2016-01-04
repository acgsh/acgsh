package main

type userGetPosts struct {
	Username string `json:"username"`
	Max_id   int64  `json:"max_id"`
	Since_id int64  `json:"since_id"`
}
