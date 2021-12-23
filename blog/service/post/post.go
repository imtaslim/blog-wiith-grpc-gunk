package post

import (
	"context"

	"blog-gunk/blog/storage"
	pgp "blog-gunk/gunk/v1/post"
)

type postCoreStore interface{
	CreatePost(context.Context, storage.Post) (int64, error)
	GetPost(context.Context, int64) (storage.Post, error)
	UpdatePost(context.Context, storage.Post) error
	DeletePost(context.Context, int64) error
	ListPost(context.Context) ([]storage.Post, error)
	PaginateSearch(context.Context, int64, int64, string) ([]storage.Post, int64, error)
}

type PostSvc struct{
	pgp.UnimplementedPostServiceServer
	core postCoreStore
}

func NewPostServer(p postCoreStore) *PostSvc {
	return &PostSvc{
		core: p,
	}
}