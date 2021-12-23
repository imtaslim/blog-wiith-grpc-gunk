package postgres

import (
	"blog-gunk/blog/storage"
	"context"
)

const insertPost = `
	INSERT INTO posts(
		cat_id,
		title,
		description,
		image
	) VALUES(
		:cat_id,
		:title,
		:description,
		:image
	)
	RETURNING id
`

func (s *Storage) CreatePost(ctx context.Context, c storage.Post) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertPost)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, c); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	var c storage.Post

	if err := s.db.Get(&c, "SELECT posts.id, cat_id, title, description, image, categories.name FROM posts LEFT JOIN categories ON categories.id = cat_id WHERE posts.id = $1", id); err != nil {
		return c, err
	}
	return c, nil
}

func (s *Storage) ListPost(ctx context.Context) ([]storage.Post, error) {
	var c []storage.Post

	if err := s.db.Select(&c, "SELECT posts.id, cat_id, title, description, image, categories.name FROM posts LEFT JOIN categories ON categories.id = cat_id;"); err != nil {
		return c, err
	}
	return c, nil
}

const updatePost = `
	UPDATE posts set 
		cat_id = :cat_id,
		title = :title,
		description = :description,
		image = :image
	Where
		id = :id
	RETURNING id
`

func (s *Storage) UpdatePost(ctx context.Context, req storage.Post) error {
	stmt, err := s.db.PrepareNamed(updatePost)
	if err != nil {
		return err
	}
	var id int64
	if err := stmt.Get(&id, req); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeletePost(ctx context.Context, id int64) error {
	var c storage.Post
	err := s.db.Get(&c, "DELETE FROM posts WHERE id=$1 RETURNING *", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) PaginateSearch(ctx context.Context, offset int64, limit int64, search string) ([]storage.Post, int64, error) {
	var c []storage.Post
	var total int64
	if err := s.db.Select(&c, "SELECT posts.id, cat_id, title, description, image, categories.name FROM posts LEFT JOIN categories ON categories.id = cat_id WHERE title ILIKE '%%' || $1 || '%%' OR name ILIKE '%%' || $1 || '%%' OFFSET $2 LIMIT $3", search, offset, limit); err != nil {
		return c, 0, err
	}
	if err := s.db.Get(&total, "SELECT count(posts.*) FROM posts LEFT JOIN categories ON categories.id = cat_id WHERE title ILIKE '%%' || $1 || '%%' OR name ILIKE '%%' || $1 || '%%'", search); err != nil {
		return c, 0, err
	}
	return c, total, nil
}
