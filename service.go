package postsService

func GetAll(getter PostsGetter) ([]Post, error) {
	posts, err := getter.GetPosts()
	if err != nil {
		return nil, ErrInternal
	}
	return posts, nil
}

func GetOne(getter PostGetter, id int) (Post, error) {
	post, err := getter.GetPost(id)
	if err != nil {
		return Post{}, determineError(err)
	}
	return post, nil
}

func Create(saver PostSaver, post Post) (Post, error) {
	if err := verifyPost(post); err != nil {
		return Post{}, err
	}
	post = formatPost(post)
	return attemptSavePost(saver, post)
}

func Delete(getter PostGetter, deleter PostDeleter, id int) (Post, error) {
	post, err := GetOne(getter, id)
	if err != nil {
		return Post{}, err
	}
	return attemptDelete(deleter, id, post)
}

func Update(getter PostGetter, updater PostUpdater, id int, updateData Post) (Post, error) {
	post, err := getter.GetPost(id)
	if err != nil {
		return Post{}, determineError(err)
	}
	return verifyFieldsAndUpdatePost(post, updater, id, updateData)
}

func determineError(err error) error {
	if err == ErrNotFound {
		return ErrNotFound
	}
	return ErrInternal
}

func verifyPost(post Post) error {
	if !post.hasTitle() {
		return ErrNeedsTitle
	}
	if post.contentIsTooLong() {
		return ErrTooLong
	}
	return nil
}

func formatPost(post Post) Post {
	post.Likes = 0
	post.Dislikes = 0
	return post
}

func attemptSavePost(repo PostSaver, post Post) (Post, error) {
	if err := repo.SavePost(post); err != nil {
		return Post{}, ErrInternal
	}
	return post, nil
}

func attemptDelete(deleter PostDeleter, id int, post Post) (Post, error) {
	if err := deleter.DeletePost(id); err != nil {
		return Post{}, determineError(deleter.DeletePost(id))
	}
	return post, nil
}

func verifyFieldsAndUpdatePost(original Post, updater PostUpdater, id int, updateData Post) (Post, error) {
	err := verifyFields(updateData)
	if err != nil {
		return Post{}, err
	}
	post := updateFields(original, updateData)
	return attemptUpdatePost(updater, post, id)
}

func verifyFields(updateData Post) error {
	if updateData.Likes != 0 || updateData.Dislikes != 0 {
		return ErrCantChangeLikes
	}
	return nil
}

func updateFields(original Post, updateData Post) Post {
	if updateData.Title != "" {
		original.Title = updateData.Title
	}
	if updateData.Content != "" {
		original.Content = updateData.Content
	}
	return original
}

func attemptUpdatePost(updater PostUpdater, post Post, id int) (Post, error) {
	err := updater.UpdatePost(id, post)
	if err != nil {
		return Post{}, determineError(err)
	}
	return post, nil
}
