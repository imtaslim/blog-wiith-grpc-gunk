package post

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog-gunk/blog/storage"
	pgp "blog-gunk/gunk/v1/post"
)

func (s *PostSvc) Create(ctx context.Context, req *pgp.CreatePostRequest ) (*pgp.CreatePostResponse, error) {
	post := storage.Post{
		ID: req.GetPost().ID,
		CatID: req.GetPost().CatID,
		Title: req.GetPost().Title,
		Description: req.GetPost().Description,
		Image: req.GetPost().Image,
	}
	id, err := s.core.CreatePost(context.Background(), post)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to create post.")
	}
	return &pgp.CreatePostResponse{
		ID: id,
	}, nil
}

func (s *PostSvc) Get(ctx context.Context, req *pgp.GetPostRequest) (*pgp.GetPostResponse, error) {
	var post storage.Post

	post, err := s.core.GetPost(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to get post.")
	}

	return  &pgp.GetPostResponse{
		Post : &pgp.Post{
			ID: post.ID,
			CatID: post.CatID,
			Title: post.Title,
			Description: post.Description,
			Image: post.Image,
			CatName: post.CatName,
		},
	}, nil
}

func (s *PostSvc) List(ctx context.Context, req *pgp.GetsPostRequest) (*pgp.GetsPostResponse, error) {
	post, err := s.core.ListPost(context.Background())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to get category.")
	}

	var data []*pgp.Post
	for _, value := range post {
		data = append(data, &pgp.Post{
			ID:     value.ID,
			CatID:     value.CatID,
			Title:   value.Title,
			Description: value.Description,
			Image: value.Image,
			CatName: value.CatName,
		})
	}

	return  &pgp.GetsPostResponse{
		Post: data,
	}, nil
}

func (s *PostSvc) PaginateSearch(ctx context.Context, req *pgp.GetsPSRequest) (*pgp.GetsPSResponse, error) {
	post, total, err := s.core.PaginateSearch(context.Background(), req.Offset, req.Limit, req.Search)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to get category.")
	}

	var data []*pgp.Post
	for _, value := range post {
		data = append(data, &pgp.Post{
			ID:     value.ID,
			CatID:     value.CatID,
			Title:   value.Title,
			Description: value.Description,
			Image: value.Image,
			CatName: value.CatName,
		})
	}

	return  &pgp.GetsPSResponse{
		Post: data,
		Total: total,
	}, nil
}