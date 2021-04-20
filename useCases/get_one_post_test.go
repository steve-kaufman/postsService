package useCases_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/useCases"
)

func TestGetOne_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	testIDs := []int{-10, -1, 0, 1, 2, 3, 4, 5, 10, 100}

	for _, id := range testIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := new(db.BadRepository)
			post, err := useCases.GetOnePost(repo, id)

			if err == nil {
				t.Fatal("Expected an error")
			}
			if err != useCases.ErrInternal {
				t.Fatalf("Expected ErrInternal; Got: '%v'", err)
			}
			if (post != entities.Post{}) {
				t.Fatalf("Expected empty post; Got: '%v'", post)
			}
		})
	}
}

func TestGetOne_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	outOfBoundsIDs := []int{-10, -1, 0, 4, 5, 10, 100}

	for _, id := range outOfBoundsIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			post, err := useCases.GetOnePost(repo, id)

			if err == nil {
				t.Fatal("Expected an error")
			}
			if err != useCases.ErrNotFound {
				t.Fatalf("Expected useCases.ErrNotFound; Got: '%v'", err)
			}
			if (post != entities.Post{}) {
				t.Fatalf("Expected empty post; Got: '%v'", post)
			}
		})
	}
}

func TestGetOne_ReturnsCorrrectPost_FromGoodRepo(t *testing.T) {
	testIDs := []int{1, 2, 3}

	for _, id := range testIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			post, err := useCases.GetOnePost(repo, id)

			if err != nil {
				t.Fatalf("Expected no error; Got: '%v'", err)
			}

			if diff := cmp.Diff(post, examplePosts[id-1]); diff != "" {
				t.Fatalf("Expected posts to match; Got:\nDiff: %s", diff)
			}
		})
	}
}
