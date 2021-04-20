package useCases

import (
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
)

func GetAllPosts(getter interfaces.PostsGetter) ([]entities.Post, error) {
	posts, err := getter.GetPosts()
	if err != nil {
		return nil, ErrInternal
	}
	return posts, nil
}
