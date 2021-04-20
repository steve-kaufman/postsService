package useCases_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/useCases"
)

func TestUpdate_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(db.BadRepository)
	_, err := useCases.UpdatePost(repo, repo, 1, entities.Post{Title: "Foo"})

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != useCases.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
}

func TestUpdate_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	badIDs := []int{-10, -1, 0, 4, 5, 10}

	for _, id := range badIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			_, err := useCases.UpdatePost(repo, repo, 0, entities.Post{Title: "Foo"})

			if err != useCases.ErrNotFound {
				t.Fatalf("Expected useCases.ErrNotFound; Got: '%v'", err)
			}
		})
	}
}

type UpdateTest struct {
	name          string
	inputID       int
	updateData    entities.Post
	expectedPost  entities.Post
	expectedError error
}

var updateTests = []UpdateTest{
	{
		name:          "Changing Likes Returns ErrCantChangeLikes",
		inputID:       1,
		updateData:    entities.Post{Likes: 3},
		expectedError: useCases.ErrCantChangeLikes,
	},
	{
		name:          "Changing Dislikes Returns ErrCantChangeLikes",
		inputID:       1,
		updateData:    entities.Post{Dislikes: 2},
		expectedError: useCases.ErrCantChangeLikes,
	},
	{
		name:       "Change title on post 1",
		inputID:    1,
		updateData: entities.Post{Title: "Foo"},
		expectedPost: entities.Post{
			ID:       1,
			Title:    "Foo",
			Content:  "Content of Post 1",
			Likes:    2,
			Dislikes: 1,
		},
	},
	{
		name:       "Change content on post 2",
		inputID:    2,
		updateData: entities.Post{Content: "Bar"},
		expectedPost: entities.Post{
			ID:       2,
			Title:    "Post 2",
			Content:  "Bar",
			Likes:    5,
			Dislikes: 2,
		},
	},
}

func TestUpdate_WithGoodRepo(t *testing.T) {
	for _, tc := range updateTests {
		t.Run(tc.name, func(t *testing.T) {
			repo := db.NewGoodRepository(examplePosts)
			post, err := useCases.UpdatePost(repo, repo, tc.inputID, tc.updateData)

			if err != tc.expectedError {
				t.Fatalf("Expected error '%v'; Got: '%v'", tc.expectedError, err)
			}

			if diff := cmp.Diff(tc.expectedPost, post); diff != "" {
				t.Fatalf("Expected posts to match: \n%s", diff)
			}

			if diff := cmp.Diff(tc.expectedPost, repo.UpdatedPost); diff != "" {
				t.Fatalf("Expected post to be updated: \n%s", diff)
			}
		})
	}
}
