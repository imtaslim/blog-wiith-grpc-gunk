package category

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog-gunk/blog/storage"
	cgp "blog-gunk/gunk/v1/category"
)

func (s *Svc) Create(ctx context.Context, req *cgp.CreateCategoryRequest ) (*cgp.CreateCategoryResponse, error) {
	category := storage.Category{
		ID: req.GetCategory().ID,
		Name: req.GetCategory().Name,
	}
	id, err := s.core.Create(context.Background(), category)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to create category.")
	}
	return &cgp.CreateCategoryResponse{
		ID: id,
	}, nil
}

func (s *Svc) Get(ctx context.Context, req *cgp.GetCategoryRequest) (*cgp.GetCategoryResponse, error) {
	var cat storage.Category

	cat, err := s.core.Get(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to get category.")
	}

	return  &cgp.GetCategoryResponse{
		Category : &cgp.Category{
			ID: cat.ID,
			Name: cat.Name,
			Status: cat.Status,
		},
	}, nil
}

func (s *Svc) Gets(ctx context.Context, req *cgp.GetsCategoryRequest) (*cgp.GetsCategoryResponse, error) {
	cat, err := s.core.Gets(context.Background())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to get category.")
	}

	var data []*cgp.Category
	for _, value := range cat {
		data = append(data, &cgp.Category{
			ID:     value.ID,
			Name:   value.Name,
			Status: value.Status,
		})
	}

	return  &cgp.GetsCategoryResponse{
		Category: data,
	}, nil
}