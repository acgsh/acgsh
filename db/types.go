package db

type ShPubReplyPost struct {
	//SigUserpost string `json:"sig_userpost"`
	Msg   string `json:"msg"`
	K     int64  `json:"k"`
	Lastk int64  `json:"lastk"`
	N     string `json:"n"`
	//Height int64  `json:"height"`
	Time uint64 `json:"time"`
}
type ShPubReplyPosts []ShPubReplyPost

type ShPost struct {
	//SigUserpost string `json:"sig_userpost"`
	Msg   string `json:"msg"`
	N     string `json:"n"`
	K     int64  `json:"k"`
	Lastk int64  `json:"lastk"`
	//Height int64  `json:"height"`

	Time     uint64 `json:"time"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Magnet   string `json:"magnet"`
	Size     uint64 `json:"size"`
	Team     string `json:"team"`
}

type ShPosts []ShPost
