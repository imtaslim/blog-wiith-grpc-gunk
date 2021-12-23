package storage

type Category struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	Status bool   `db:"status"`
}

type Post struct {
	ID          int64  `db:"id"`
	CatID       int64  `db:"cat_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
	CatName     string `db:"name"`
}
