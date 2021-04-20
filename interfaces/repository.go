package interfaces

import "github.com/steve-kaufman/postsService/entities"

type PostsGetter interface {
	GetPosts() ([]entities.Post, error)
}

type PostGetter interface {
	GetPost(id int) (entities.Post, error)
}

type PostSaver interface {
	SavePost(post entities.Post) error
}

type PostDeleter interface {
	DeletePost(id int) error
}

type PostUpdater interface {
	UpdatePost(id int, data entities.Post) error
}
