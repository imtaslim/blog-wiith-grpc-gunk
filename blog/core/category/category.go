package category

import (
	"context"

	"blog-gunk/blog/storage"
)

type categoryStore interface {
	Create(context.Context, storage.Category) (int64, error)
	Get(context.Context, int64) (storage.Category, error)
	Update(context.Context, storage.Category) error
	Delete(context.Context, int64) error
	Gets(context.Context) ([]storage.Category, error)
}

type CoreSvc struct {
	store categoryStore
}

func NewCoreSvc(s categoryStore) *CoreSvc {
	return &CoreSvc{
		store: s,
	}
}

func (cs CoreSvc) Create(ctx context.Context, c storage.Category) (int64, error) {
	return cs.store.Create(ctx,c)
}

func (cs CoreSvc) Get(ctx context.Context, id int64) (storage.Category, error) {
	return cs.store.Get(ctx, id)
}

func (cs CoreSvc) Gets(ctx context.Context) ([]storage.Category, error) {
	return cs.store.Gets(ctx)
}

func (cs CoreSvc) Update(ctx context.Context, c storage.Category) error {
	return cs.store.Update(ctx, c)
}

func (cs CoreSvc) Delete(ctx context.Context, id int64) error {
	return cs.store.Delete(ctx, id)
}