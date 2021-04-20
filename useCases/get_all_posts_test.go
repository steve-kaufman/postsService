package useCases_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/useCases"
)

var examplePosts = []entities.Post{
	{
		ID:       1,
		Title:    "Post 1",
		Content:  "Content of Post 1",
		Likes:    2,
		Dislikes: 1,
	},
	{
		ID:       2,
		Title:    "Post 2",
		Content:  "Content of Post 2",
		Likes:    5,
		Dislikes: 2,
	},
	{
		ID:       3,
		Title:    "Post 3",
		Content:  "Content of Post 3",
		Likes:    0,
		Dislikes: 10,
	},
}

func TestGetAll_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(db.BadRepository)
	posts, err := useCases.GetAllPosts(repo)

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != useCases.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
	if posts != nil {
		t.Fatal("Expected no posts")
	}
}
func TestGetAll_ReturnsPosts_FromGoodRepo(t *testing.T) {
	repo := db.NewGoodRepository(examplePosts)
	posts, err := useCases.GetAllPosts(repo)

	if err != nil {
		t.Fatalf("Expected no error; Got: '%s'", err)
	}
	if diff := cmp.Diff(posts, examplePosts); diff != "" {
		t.Fatalf("Expected posts from database:\nDiff: %s", diff)
	}
}
