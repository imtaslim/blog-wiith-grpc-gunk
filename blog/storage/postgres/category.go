package postgres

import (
	"blog-gunk/blog/storage"
	"context"
)

const insertCategory = `
	INSERT INTO categories(
		name
	) VALUES(
		:name
	)
	RETURNING id
`

func (s *Storage) Create(ctx context.Context, c storage.Category) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertCategory)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, c); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) Get(ctx context.Context, id int64) (storage.Category, error) {
	var c storage.Category
	
	if err := s.db.Get(&c, "SELECT * FROM categories WHERE id = $1", id); err != nil{
		return c, err
	}
	return c, nil
}

func (s *Storage) Gets(ctx context.Context) ([]storage.Category, error) {
	var c []storage.Category
	
	if err := s.db.Select(&c, "SELECT * FROM categories"); err != nil{
		return c, err
	}
	return c, nil
}

const updateCategory = `
	UPDATE categories set 
		name = :name
	Where
		id = :id
	RETURNING id
`

func (s *Storage) Update(ctx context.Context, req storage.Category) error {
	stmt, err := s.db.PrepareNamed(updateCategory)
	if err != nil {
		return err
	}
	var id int64
	if err := stmt.Get(&id, req); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	var c storage.Category
	err := s.db.Get(&c, "DELETE FROM categories WHERE id=$1 RETURNING *", id)
	if err != nil {
		return err
	}
	return nil
}