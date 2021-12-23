package post

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog-gunk/blog/storage"
	pgp "blog-gunk/gunk/v1/post"
)

func (s *PostSvc) Update(ctx context.Context, req *pgp.UpdatePostRequest ) (*pgp.UpdatePostResponse, error) {
	post := storage.Post{
		ID: req.GetPost().ID,
		CatID: req.GetPost().CatID,
		Title: req.GetPost().Title,
		Description: req.GetPost().Description,
		Image: req.GetPost().Image,
	}
	err := s.core.UpdatePost(context.Background(), post)
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to Update post.")
	}
	return  &pgp.UpdatePostResponse{}, nil
}

func (s *PostSvc) Delete(ctx context.Context, req *pgp.DeletePostRequest ) (*pgp.DeletePostResponse, error) {
	err := s.core.DeletePost(context.Background(), req.GetID())
	if err != nil{
		return nil, status.Error(codes.Internal, "failed to Delete post.")
	}
	return  &pgp.DeletePostResponse{}, nil
}