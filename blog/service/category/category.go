package category

import (
	"context"

	"blog-gunk/blog/storage"
	cgp "blog-gunk/gunk/v1/category"
)

type categoryCoreStore interface{
	Create(context.Context, storage.Category) (int64, error)
	Get(context.Context, int64) (storage.Category, error)
	Update(context.Context, storage.Category) error
	Delete(context.Context, int64) error
	Gets(context.Context) ([]storage.Category, error)
}

type Svc struct{
	cgp.UnimplementedCategoryServiceServer
	core categoryCoreStore
}

func NewCategoryServer(c categoryCoreStore) *Svc {
	return &Svc{
		core: c,
	}
}