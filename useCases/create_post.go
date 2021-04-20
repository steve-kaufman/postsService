package useCases

import (
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
)

func CreatePost(saver interfaces.PostSaver, post entities.Post) (entities.Post, error) {
	post, err := entities.FormatAndValidateNewPost(post)
	if err != nil {
		return entities.Post{}, err
	}
	return attemptSavePost(saver, post)
}

func attemptSavePost(saver interfaces.PostSaver, post entities.Post) (entities.Post, error) {
	if err := saver.SavePost(post); err != nil {
		return entities.Post{}, ErrInternal
	}
	return post, nil
}
