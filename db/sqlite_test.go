package db_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/steve-kaufman/postsService/db"
	"github.com/steve-kaufman/postsService/service"

	_ "github.com/mattn/go-sqlite3"
)

var testDBPath = "./test.db"

func setup() (*db.SqliteRepo, *sql.DB) {
	os.Remove(testDBPath)
	repo := db.NewSqliteRepo(testDBPath)
	conn, _ := sql.Open("sqlite3", testDBPath)
	return repo, conn
}

func TestInstantiatingRepo_CreatesDBFile(t *testing.T) {
	setup()

	_, err := os.Stat(testDBPath)

	if err != nil {
		t.Fatal("Didn't create db file")
	}
}

func TestInstantiatingRepo_GeneratesTable(t *testing.T) {
	_, conn := setup()

	row, err := conn.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name='posts';`)
	if err != nil {
		t.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		t.Fatal("Expected table 'posts' to exist")
	}
}

func TestAfterInstantiatingRepo_CanInsertPost(t *testing.T) {
	_, conn := setup()

	_, err := conn.Exec(`INSERT INTO posts (title, content, likes, dislikes) VALUES('Foo', 'Bar', 2, 1);`)
	if err != nil {
		t.Fatalf("Expected no error inserting post; Got: '%v'", err)
	}
}

var examplePosts = []service.Post{
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

func TestGetPosts_ReturnsAllPosts(t *testing.T) {
	repo, conn := setup()

	insertExamplePosts(conn)

	posts, err := repo.GetPosts()
	if err != nil {
		t.Fatalf("Expected no error; Got: '%v'", err)
	}

	if diff := cmp.Diff(examplePosts, posts); diff != "" {
		t.Fatalf("Expected example posts; Got: '%s'", diff)
	}
}

func TestGetPost_ReturnsNotFound_IfNoMatch(t *testing.T) {
	repo, conn := setup()

	insertExamplePosts(conn)

	_, err := repo.GetPost(4)

	if err != service.ErrNotFound {
		t.Fatalf("Expected ErrNotFound; Got: '%v'", err)
	}
}

func TestGetPost_ReturnsCorrectPost_WithGoodID(t *testing.T) {
	repo, conn := setup()

	insertExamplePosts(conn)

	post, err := repo.GetPost(2)

	if err != nil {
		t.Fatalf("Expected no error; Got: '%v'", err)
	}

	if diff := cmp.Diff(examplePosts[1], post); diff != "" {
		t.Fatalf("Expected post 2; Got: \n%s", diff)
	}
}

func TestSavePost_InsertsPostInDB(t *testing.T) {
	repo, conn := setup()

	insertExamplePosts(conn)

	err := repo.SavePost(service.Post{
		Title:    "Foo",
		Content:  "Bar",
		Likes:    1,
		Dislikes: 2,
	})

	if err != nil {
		t.Fatalf("Expected no error; Got: '%v'", err)
	}

	var post service.Post
	row := conn.QueryRow(`SELECT id, title, content, likes, dislikes FROM posts WHERE id=4`)
	row.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes)

	expectedPost := service.Post{
		ID:       4,
		Title:    "Foo",
		Content:  "Bar",
		Likes:    1,
		Dislikes: 2,
	}
	if diff := cmp.Diff(expectedPost, post); diff != "" {
		t.Fatalf("Expected post to be inserted: \n%s", diff)
	}
}

func insertPost(conn *sql.DB, post service.Post) {
	conn.Exec(`INSERT INTO posts (title, content, likes, dislikes) VALUES(?, ?, ?, ?);`,
		post.Title,
		post.Content,
		post.Likes,
		post.Dislikes,
	)
}

func insertExamplePosts(conn *sql.DB) {
	for _, post := range examplePosts {
		insertPost(conn, post)
	}
}
