package main

type getUserPosts struct {
	data []interface{}
}

func (g *getUserPosts) addUser(username string, max, since int64) {

	g.data = append(g.data, userGetPosts{
		Username: username,
		Max_id:   max,
		Since_id: since,
	})
}
