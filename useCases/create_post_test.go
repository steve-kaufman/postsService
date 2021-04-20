package useCases_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/entities"
	"github.com/steve-kaufman/postsService/interfaces"
	"github.com/steve-kaufman/postsService/useCases"
)

type CreateTest struct {
	name         string
	repo         interfaces.PostSaver
	inputPost    entities.Post
	expectedErr  error
	expectedPost entities.Post
}

var createTests = []CreateTest{
	{
		name:         "Returns ErrInternal from bad repo",
		repo:         new(db.BadRepository),
		inputPost:    entities.Post{Title: "foo", Content: "bar"},
		expectedErr:  useCases.ErrInternal,
		expectedPost: entities.Post{},
	},
	{
		name:         "Returns ErrNeedsTitle if no title (bad repo)",
		repo:         new(db.BadRepository),
		inputPost:    entities.Post{},
		expectedErr:  entities.ErrNeedsTitle,
		expectedPost: entities.Post{},
	},
	{
		name:         "Returns ErrNeedsTitle if no title (good repo)",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{},
		expectedErr:  entities.ErrNeedsTitle,
		expectedPost: entities.Post{},
	},
	{
		name:         "Returns ErrTooLong if content is longer than 500 characters (bad repo)",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{Title: "Foo", Content: strings.Repeat("a", 501)},
		expectedErr:  entities.ErrTooLong,
		expectedPost: entities.Post{},
	},
	{
		name:         "Returns ErrTooLong if content is longer than 500 characters (good repo)",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{Title: "Foo", Content: strings.Repeat("a", 501)},
		expectedErr:  entities.ErrTooLong,
		expectedPost: entities.Post{},
	},
	{
		name:         "Saves post if title and length of content <= 500",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{Title: "Foo", Content: "Bar"},
		expectedErr:  nil,
		expectedPost: entities.Post{Title: "Foo", Content: "Bar"},
	},
	{
		name:         "Saves post if title and length of content <= 500",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{Title: "Foo", Content: strings.Repeat("a", 500)},
		expectedErr:  nil,
		expectedPost: entities.Post{Title: "Foo", Content: strings.Repeat("a", 500)},
	},
	{
		name:         "Sets likes and dislikes to zero regardless of input",
		repo:         db.NewGoodRepository(examplePosts),
		inputPost:    entities.Post{Title: "Foo", Content: "Bar", Likes: 11, Dislikes: 2},
		expectedErr:  nil,
		expectedPost: entities.Post{Title: "Foo", Content: "Bar"},
	},
}

func TestCreate(t *testing.T) {
	for _, tc := range createTests {
		t.Run(tc.name, func(t *testing.T) {
			post, err := useCases.CreatePost(tc.repo, tc.inputPost)

			if err != tc.expectedErr {
				t.Fatalf("Expected err to be: '%v'; Got: '%v'", tc.expectedErr, err)
			}
			if diff := cmp.Diff(tc.expectedPost, post); diff != "" {
				t.Fatalf("Expected posts to match: %s", diff)
			}

			goodRepo, ok := tc.repo.(*db.GoodRepository)
			if !ok {
				return
			}
			if diff := cmp.Diff(tc.expectedPost, goodRepo.SavedPost); diff != "" {
				t.Fatalf("Expected post to be saved: %s", diff)
			}
		})
	}
}
