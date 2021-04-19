package service

type PostsGetter interface {
	GetPosts() ([]Post, error)
}

type PostGetter interface {
	GetPost(id int) (Post, error)
}

type PostSaver interface {
	SavePost(post Post) error
}

type PostDeleter interface {
	DeletePost(id int) error
}

type PostUpdater interface {
	UpdatePost(id int, data Post) error
}
