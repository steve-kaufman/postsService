package useCases

import (
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
)

func DeletePost(getter interfaces.PostGetter, deleter interfaces.PostDeleter, id int) (entities.Post, error) {
	post, err := GetOnePost(getter, id)
	if err != nil {
		return entities.Post{}, err
	}
	return attemptDelete(deleter, id, post)
}

func attemptDelete(deleter interfaces.PostDeleter, id int, post entities.Post) (entities.Post, error) {
	if err := deleter.DeletePost(id); err != nil {
		return entities.Post{}, determineError(deleter.DeletePost(id))
	}
	return post, nil
}
