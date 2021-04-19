package db

import (
	"database/sql"
	"os"

	"github.com/steve-kaufman/postsService/service"
)

type SqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(path string) *SqliteRepo {
	os.Create(path)

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	conn.Exec(`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY,
		title TEXT,
		content TEXT,
		likes INTEGER,
		dislikes INTEGER
	);`)

	repo := new(SqliteRepo)
	repo.conn = conn

	return repo
}

func (repo SqliteRepo) GetPosts() ([]service.Post, error) {
	rows, err := repo.conn.Query(`SELECT id, title, content, likes, dislikes FROM posts;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRowsToPosts(rows)
}

func (repo SqliteRepo) GetPost(id int) (service.Post, error) {
	row := repo.conn.QueryRow(`SELECT id, title, content, likes, dislikes FROM posts WHERE id=?`, id)
	post, err := mapToPost(row)
	if err == sql.ErrNoRows {
		return service.Post{}, service.ErrNotFound
	}
	return post, nil
}

func (repo SqliteRepo) SavePost(post service.Post) error {
	_, err := repo.conn.Exec(`INSERT INTO posts (title, content, likes, dislikes) VALUES (?, ?, ?, ?);`,
		post.Title,
		post.Content,
		post.Likes,
		post.Dislikes,
	)
	return err
}

func mapRowsToPosts(rows *sql.Rows) ([]service.Post, error) {
	posts := []service.Post{}
	for rows.Next() {
		post, err := mapToPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

func mapToPost(row RowScanner) (service.Post, error) {
	var post service.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes)
	if err != nil {
		return service.Post{}, err
	}
	return post, nil
}
