package postsService_test

import (
	"errors"
	"fmt"
	"postsService"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var examplePosts = []postsService.Post{
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

// BadRepository indiscriminately returns an obscure error for every method
type BadRepository struct{}

var ErrBad = errors.New("something bad went wrong")

func (BadRepository) GetPosts() ([]postsService.Post, error) {
	return nil, ErrBad
}

func (BadRepository) GetPost(id int) (postsService.Post, error) {
	return postsService.Post{}, ErrBad
}

func (BadRepository) SavePost(post postsService.Post) error {
	return ErrBad
}

func (BadRepository) DeletePost(id int) error {
	return ErrBad
}

func (BadRepository) UpdatePost(id int, data postsService.Post) error {
	return ErrBad
}

// GoodRepository is a quasi-functional in-memory repository for the service
type GoodRepository struct {
	savedPost     postsService.Post
	deletedPostID int
	updatedPost   postsService.Post
}

func (GoodRepository) GetPosts() ([]postsService.Post, error) {
	return examplePosts, nil
}

func (GoodRepository) GetPost(id int) (postsService.Post, error) {
	if id < 1 || id > len(examplePosts) {
		return postsService.Post{}, postsService.ErrNotFound
	}
	return examplePosts[id-1], nil
}

func (repo *GoodRepository) SavePost(post postsService.Post) error {
	repo.savedPost = post
	return nil
}

func (repo *GoodRepository) DeletePost(id int) error {
	if id < 1 || id > len(examplePosts) {
		return postsService.ErrNotFound
	}
	repo.deletedPostID = id
	return nil
}

func (repo *GoodRepository) UpdatePost(id int, post postsService.Post) error {
	if id < 1 || id > len(examplePosts) {
		return postsService.ErrNotFound
	}

	repo.updatedPost = post
	return nil
}

/** Tests **/

func TestGetAll_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(BadRepository)
	posts, err := postsService.GetAll(repo)

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != postsService.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
	if posts != nil {
		t.Fatal("Expected no posts")
	}
}
func TestGetAll_ReturnsPosts_FromGoodRepo(t *testing.T) {
	repo := new(GoodRepository)
	posts, err := postsService.GetAll(repo)

	if err != nil {
		t.Fatalf("Expected no error; Got: '%s'", err)
	}
	if diff := cmp.Diff(posts, examplePosts); diff != "" {
		t.Fatalf("Expected posts from database:\nDiff: %s", diff)
	}
}

func TestGetOne_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	testIDs := []int{-10, -1, 0, 1, 2, 3, 4, 5, 10, 100}

	for _, id := range testIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := new(BadRepository)
			post, err := postsService.GetOne(repo, id)

			if err == nil {
				t.Fatal("Expected an error")
			}
			if err != postsService.ErrInternal {
				t.Fatalf("Expected ErrInternal; Got: '%v'", err)
			}
			if (post != postsService.Post{}) {
				t.Fatalf("Expected empty post; Got: '%v'", post)
			}
		})
	}
}

func TestGetOne_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	outOfBoundsIDs := []int{-10, -1, 0, 4, 5, 10, 100}

	for _, id := range outOfBoundsIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := new(GoodRepository)
			post, err := postsService.GetOne(repo, id)

			if err == nil {
				t.Fatal("Expected an error")
			}
			if err != postsService.ErrNotFound {
				t.Fatalf("Expected ErrNotFound; Got: '%v'", err)
			}
			if (post != postsService.Post{}) {
				t.Fatalf("Expected empty post; Got: '%v'", post)
			}
		})
	}
}

func TestGetOne_ReturnsCorrrectPost_FromGoodRepo(t *testing.T) {
	testIDs := []int{1, 2, 3}

	for _, id := range testIDs {
		t.Run(fmt.Sprintf("With ID '%d'", id), func(t *testing.T) {
			repo := new(GoodRepository)
			post, err := postsService.GetOne(repo, id)

			if err != nil {
				t.Fatalf("Expected no error; Got: '%v'", err)
			}

			if diff := cmp.Diff(post, examplePosts[id-1]); diff != "" {
				t.Fatalf("Expected posts to match; Got:\nDiff: %s", diff)
			}
		})
	}
}

type CreateTest struct {
	name         string
	repo         postsService.PostSaver
	inputPost    postsService.Post
	expectedErr  error
	expectedPost postsService.Post
}

var createTests = []CreateTest{
	{
		name:         "Returns ErrInternal from bad repo",
		repo:         new(BadRepository),
		inputPost:    postsService.Post{Title: "foo", Content: "bar"},
		expectedErr:  postsService.ErrInternal,
		expectedPost: postsService.Post{},
	},
	{
		name:         "Returns ErrNeedsTitle if no title (bad repo)",
		repo:         new(BadRepository),
		inputPost:    postsService.Post{},
		expectedErr:  postsService.ErrNeedsTitle,
		expectedPost: postsService.Post{},
	},
	{
		name:         "Returns ErrNeedsTitle if no title (good repo)",
		repo:         new(GoodRepository),
		inputPost:    postsService.Post{},
		expectedErr:  postsService.ErrNeedsTitle,
		expectedPost: postsService.Post{},
	},
	{
		name:         "Returns ErrTooLong if content is longer than 500 characters (bad repo)",
		repo:         new(BadRepository),
		inputPost:    postsService.Post{Title: "Foo", Content: strings.Repeat("a", 501)},
		expectedErr:  postsService.ErrTooLong,
		expectedPost: postsService.Post{},
	},
	{
		name:         "Returns ErrTooLong if content is longer than 500 characters (good repo)",
		repo:         new(GoodRepository),
		inputPost:    postsService.Post{Title: "Foo", Content: strings.Repeat("a", 501)},
		expectedErr:  postsService.ErrTooLong,
		expectedPost: postsService.Post{},
	},
	{
		name:         "Saves post if title and length of content <= 500",
		repo:         new(GoodRepository),
		inputPost:    postsService.Post{Title: "Foo", Content: "Bar"},
		expectedErr:  nil,
		expectedPost: postsService.Post{Title: "Foo", Content: "Bar"},
	},
	{
		name:         "Saves post if title and length of content <= 500",
		repo:         new(GoodRepository),
		inputPost:    postsService.Post{Title: "Foo", Content: strings.Repeat("a", 500)},
		expectedErr:  nil,
		expectedPost: postsService.Post{Title: "Foo", Content: strings.Repeat("a", 500)},
	},
	{
		name:         "Sets likes and dislikes to zero regardless of input",
		repo:         new(GoodRepository),
		inputPost:    postsService.Post{Title: "Foo", Content: "Bar", Likes: 11, Dislikes: 2},
		expectedErr:  nil,
		expectedPost: postsService.Post{Title: "Foo", Content: "Bar"},
	},
}

func TestCreate(t *testing.T) {
	for _, tc := range createTests {
		t.Run(tc.name, func(t *testing.T) {
			post, err := postsService.Create(tc.repo, tc.inputPost)

			if err != tc.expectedErr {
				t.Fatalf("Expected err to be: '%v'; Got: '%v'", tc.expectedErr, err)
			}
			if diff := cmp.Diff(tc.expectedPost, post); diff != "" {
				t.Fatalf("Expected posts to match: %s", diff)
			}

			goodRepo, ok := tc.repo.(*GoodRepository)
			if !ok {
				return
			}
			if diff := cmp.Diff(tc.expectedPost, goodRepo.savedPost); diff != "" {
				t.Fatalf("Expected post to be saved: %s", diff)
			}
		})
	}
}

func TestDelete_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(BadRepository)
	deletedPost, err := postsService.Delete(repo, repo, 1)

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != postsService.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
	if (deletedPost != postsService.Post{}) {
		t.Fatalf("Expected post to be empty; Got: '%v'", deletedPost)
	}
}

func TestDelete_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	badIDs := []int{-10, -1, 0, 4, 5, 10}

	for _, id := range badIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := new(GoodRepository)
			_, err := postsService.Delete(repo, repo, id)

			if err != postsService.ErrNotFound {
				t.Fatalf("Expected ErrNotFound; Got: '%v'", err)
			}
		})
	}
}

func TestDelete_DeletesCorrectPost_FromGoodRepo(t *testing.T) {
	goodIDs := []int{1, 2, 3}

	for _, id := range goodIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := new(GoodRepository)
			post, err := postsService.Delete(repo, repo, id)

			if err != nil {
				t.Fatalf("Expected no error; Got: '%v'", err)
			}

			if repo.deletedPostID != id {
				t.Fatalf("Expected post with id %d to be deleted; Got: %d", id, repo.deletedPostID)
			}

			expectedPost := examplePosts[id-1]
			if diff := cmp.Diff(expectedPost, post); diff != "" {
				t.Fatal("Expected returned post to be deleted post; Got:", diff)
			}
		})
	}
}

func TestUpdate_ReturnsErrInternal_FromBadRepo(t *testing.T) {
	repo := new(BadRepository)
	_, err := postsService.Update(repo, repo, 1, postsService.Post{Title: "Foo"})

	if err == nil {
		t.Fatal("Expected an error")
	}
	if err != postsService.ErrInternal {
		t.Fatalf("Expected ErrInternal; Got: '%v'", err)
	}
}

func TestUpdate_ReturnsErrNotFound_FromGoodRepoWithBadID(t *testing.T) {
	badIDs := []int{-10, -1, 0, 4, 5, 10}

	for _, id := range badIDs {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			repo := new(GoodRepository)
			_, err := postsService.Update(repo, repo, 0, postsService.Post{Title: "Foo"})

			if err != postsService.ErrNotFound {
				t.Fatalf("Expected ErrNotFound; Got: '%v'", err)
			}
		})
	}
}

type UpdateTest struct {
	name          string
	inputID       int
	updateData    postsService.Post
	expectedPost  postsService.Post
	expectedError error
}

var updateTests = []UpdateTest{
	{
		name:          "Changing Likes Returns ErrCantChangeLikes",
		inputID:       1,
		updateData:    postsService.Post{Likes: 3},
		expectedError: postsService.ErrCantChangeLikes,
	},
	{
		name:          "Changing Dislikes Returns ErrCantChangeLikes",
		inputID:       1,
		updateData:    postsService.Post{Dislikes: 2},
		expectedError: postsService.ErrCantChangeLikes,
	},
	{
		name:       "Change title on post 1",
		inputID:    1,
		updateData: postsService.Post{Title: "Foo"},
		expectedPost: postsService.Post{
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
		updateData: postsService.Post{Content: "Bar"},
		expectedPost: postsService.Post{
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
			repo := new(GoodRepository)
			post, err := postsService.Update(repo, repo, tc.inputID, tc.updateData)

			if err != tc.expectedError {
				t.Fatalf("Expected error '%v'; Got: '%v'", tc.expectedError, err)
			}

			if diff := cmp.Diff(tc.expectedPost, post); diff != "" {
				t.Fatalf("Expected posts to match: \n%s", diff)
			}

			if diff := cmp.Diff(tc.expectedPost, repo.updatedPost); diff != "" {
				t.Fatalf("Expected post to be updated: \n%s", diff)
			}
		})
	}
}
