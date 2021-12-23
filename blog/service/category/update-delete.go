package category

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog-gunk/blog/storage"
	cgp "blog-gunk/gunk/v1/category"
)

func (s *Svc) Update(ctx context.Context, req *cgp.UpdateCategoryRequest ) (*cgp.UpdateCategoryResponse, error) {
	category := storage.Category{
		ID: req.GetCategory().ID,
		Name: req.GetCategory().Name,
	}
	err := s.core.Update(context.Background(), category)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to Update category.")
	}
	return  &cgp.UpdateCategoryResponse{}, nil
}

func (s *Svc) Delete(ctx context.Context, req *cgp.DeleteCategoryRequest ) (*cgp.DeleteCategoryResponse, error) {
	err := s.core.Delete(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to Delete category.")
	}
	return  &cgp.DeleteCategoryResponse{}, nil
}