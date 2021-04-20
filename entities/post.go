package entities

type Post struct {
	ID       int
	Title    string
	Content  string
	Likes    int
	Dislikes int
}

func FormatAndValidateNewPost(post Post) (Post, error) {
	if err := validatePost(post); err != nil {
		return Post{}, err
	}
	return formatNewPost(post), nil
}

func validatePost(post Post) error {
	if post.Title == "" {
		return ErrNeedsTitle
	}
	if len(post.Content) > 500 {
		return ErrTooLong
	}
	return nil
}

func formatNewPost(post Post) Post {
	post.Likes = 0
	post.Dislikes = 0
	return post
}
