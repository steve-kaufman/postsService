package postsService

type Post struct {
	ID       int
	Title    string
	Content  string
	Likes    int
	Dislikes int
}

func (post Post) hasTitle() bool {
	return post.Title != ""
}

func (post Post) contentIsTooLong() bool {
	return len(post.Content) > 500
}
