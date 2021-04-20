package useCases_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/useCases"
)

func TestDelete_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(db.BadRepository)
	deletedPost, err := useCases.DeletePost(repo, repo, 1)

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != useCases.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
	if (deletedPost != entities.Post{}) {
		t.Fatalf("Expected post to be empty; Got: '%v'", deletedPost)
	}
}

func TestDelete_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	badIDs := []int{-10, -1, 0, 4, 5, 10}

	for _, id := range badIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			_, err := useCases.DeletePost(repo, repo, id)

			if err != useCases.ErrNotFound {
				t.Fatalf("Expected useCases.ErrNotFound; Got: '%v'", err)
			}
		})
	}
}

func TestDelete_DeletesCorrectPost_FromGoodRepo(t *testing.T) {
	goodIDs := []int{1, 2, 3}

	for _, id := range goodIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			post, err := useCases.DeletePost(repo, repo, id)

			if err != nil {
				t.Fatalf("Expected no error; Got: '%v'", err)
			}

			if repo.DeletedPostID != id {
				t.Fatalf("Expected post with id %d to be deleted; Got: %d", id, repo.DeletedPostID)
			}

			expectedPost := examplePosts[id-1]
			if diff := cmp.Diff(expectedPost, post); diff != "" {
				t.Fatal("Expected returned post to be deleted post; Got:", diff)
			}
		})
	}
}
