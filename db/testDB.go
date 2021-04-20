package db

import (
	"errors"

	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/useCases"
)

// BadRepository indiscriminately returns an obscure error for every method
type BadRepository struct{}

var ErrBad = errors.New("something bad went wrong")

func (BadRepository) GetPosts() ([]entities.Post, error) {
	return nil, ErrBad
}

func (BadRepository) GetPost(id int) (entities.Post, error) {
	return entities.Post{}, ErrBad
}

func (BadRepository) SavePost(post entities.Post) error {
	return ErrBad
}

func (BadRepository) DeletePost(id int) error {
	return ErrBad
}

func (BadRepository) UpdatePost(id int, data entities.Post) error {
	return ErrBad
}

// GoodRepository is a quasi-functional in-memory repository for the useCases
type GoodRepository struct {
	posts         []entities.Post
	SavedPost     entities.Post
	DeletedPostID int
	UpdatedPost   entities.Post
}

func NewGoodRepository(posts []entities.Post) *GoodRepository {
	repo := new(GoodRepository)
	repo.posts = posts
	return repo
}

func (repo GoodRepository) GetPosts() ([]entities.Post, error) {
	return repo.posts, nil
}

func (repo GoodRepository) GetPost(id int) (entities.Post, error) {
	if id < 1 || id > len(repo.posts) {
		return entities.Post{}, useCases.ErrNotFound
	}
	return repo.posts[id-1], nil
}

func (repo *GoodRepository) SavePost(post entities.Post) error {
	repo.SavedPost = post
	return nil
}

func (repo *GoodRepository) DeletePost(id int) error {
	if id < 1 || id > len(repo.posts) {
		return useCases.ErrNotFound
	}
	repo.DeletedPostID = id
	return nil
}

func (repo *GoodRepository) UpdatePost(id int, post entities.Post) error {
	if id < 1 || id > len(repo.posts) {
		return useCases.ErrNotFound
	}

	repo.UpdatedPost = post
	return nil
}
