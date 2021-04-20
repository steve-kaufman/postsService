package useCases

import (
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
)

func GetOnePost(getter interfaces.PostGetter, id int) (entities.Post, error) {
	post, err := getter.GetPost(id)
	if err != nil {
		return entities.Post{}, determineError(err)
	}
	return post, nil
}

func determineError(err error) error {
	if err == ErrNotFound {
		return ErrNotFound
	}
	return ErrInternal
}
