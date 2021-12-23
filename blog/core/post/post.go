package post

import (
	"context"

	"blog-gunk/blog/storage"
)

type postStore interface {
	CreatePost(context.Context, storage.Post) (int64, error)
	GetPost(context.Context, int64) (storage.Post, error)
	UpdatePost(context.Context, storage.Post) error
	DeletePost(context.Context, int64) error
	ListPost(context.Context) ([]storage.Post, error)
	PaginateSearch(context.Context, int64, int64, string) ([]storage.Post, int64, error)
}

type PostCoreSvc struct {
	store postStore
}

func NewCoreSvc(s postStore) *PostCoreSvc {
	return &PostCoreSvc{
		store: s,
	}
}

func (cs PostCoreSvc) CreatePost(ctx context.Context, c storage.Post) (int64, error) {
	return cs.store.CreatePost(ctx,c)
}

func (cs PostCoreSvc) GetPost(ctx context.Context, id int64) (storage.Post, error) {
	return cs.store.GetPost(ctx, id)
}

func (cs PostCoreSvc) ListPost(ctx context.Context) ([]storage.Post, error) {
	return cs.store.ListPost(ctx)
}

func (cs PostCoreSvc) UpdatePost(ctx context.Context, c storage.Post) error {
	return cs.store.UpdatePost(ctx, c)
}

func (cs PostCoreSvc) DeletePost(ctx context.Context, id int64) error {
	return cs.store.DeletePost(ctx, id)
}

func (cs PostCoreSvc) PaginateSearch(ctx context.Context, offset int64, limit int64, search string) ([]storage.Post, int64, error) {
	return cs.store.PaginateSearch(ctx, offset, limit, search)
}