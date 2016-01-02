package rpc

type TwisterPost struct {
	SigUserpost string `json:"sig_userpost"`
	Userpost    struct {
		Height int64  `json:"height"`
		K      int64  `json:"k"`
		Lastk  *int64 `json:"lastk"`
		Msg    string `json:"msg"`
		N      string `json:"n"`
		Reply  *struct {
			K int64  `json:"k"`
			N string `json:"n"`
		} `json:"reply"`
		Rt *struct {
			Height int64  `json:"height"`
			K      int64  `json:"k"`
			Lastk  int64  `json:"lastk"`
			Msg    string `json:"msg"`
			N      string `json:"n"`
			Reply  struct {
				K int64  `json:"k"`
				N string `json:"n"`
			} `json:"reply"`
			Time int64 `json:"time"`
		} `json:"rt"`
		SigRt string `json:"sig_rt"`
		Time  uint64 `json:"time"`
	} `json:"userpost"`
}

type TwisterPosts []TwisterPost
