package service

import (
	"context"
	"errors"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostServiceImpl interface {
	CreatePostService(ctx context.Context, input *model.PostInput)(*model.Post, error)
	GetPostByIDService(ctx context.Context, input *model.GetPostInput)(*model.Post, error)
	GetAllPostService(ctx context.Context, page *model.PostsPageReq)(*model.PostsPageRes, error)
	UpdatePostService(ctx context.Context, new *model.UpdatePostInput, getPost *model.GetPostInput)(*model.Post, error)
	DeletePostService(ctx context.Context, getPost *model.GetPostInput)error
}

type PostService struct {
	repo repository.PostRepoImpl
	db *pgxpool.Pool
}

func NewServiceImpl(repo repository.PostRepoImpl, db *pgxpool.Pool)PostServiceImpl{
	return &PostService{
		repo:repo,
		db:db,
	}
}

func(s *PostService)CreatePostService(ctx context.Context, input *model.PostInput)(*model.Post, error){
	post, err := model.NewPost(ctx, input)
	if err != nil {
		helper.ErrMsg(err, "failed to create posts")
		return nil, err
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transaction")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)
	if err := s.repo.CreatePostRepo(ctx, tx, post); err != nil {
		helper.ErrMsg(err, "failed to create post (db err)")
		return nil, err
	}
	return post, nil
}
func(s *PostService)GetPostByIDService(ctx context.Context, input *model.GetPostInput)(*model.Post, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transactions")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	post, err := s.repo.GetPostByIDRepo(ctx,tx, input)
	if err != nil {
		helper.ErrMsg(err, "failed to get post")
		return nil, err
	}
	return post, nil
}
func(s *PostService)GetAllPostService(ctx context.Context, page *model.PostsPageReq)(*model.PostsPageRes, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transactions")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)

	if page.Username == "" {
		return nil, errors.New("invalid username")
	}
	if page.Limit <= 0 {
		page.Limit = 10
	}
	if page.Offset < 0 {
		page.Offset = 0
	}
	res, err := s.repo.GetAllPostRepo(ctx, tx, page)
	if err != nil {
		helper.ErrMsg(err, "failed to get posts page")
		return nil, err
	}
	return res, nil
}
func(s *PostService)UpdatePostService(ctx context.Context, new *model.UpdatePostInput, getPost *model.GetPostInput)(*model.Post, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transaction: ")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)
	existingPost, err := s.repo.GetPostByIDRepo(ctx, tx, getPost)
	if err != nil {
		helper.ErrMsg(err, "failed to get post: ")
		return nil, err
	}
	if new.Content == "" {
		new.Content = existingPost.Content
	}
	new.UpdatedAt = time.Now()
	res, err := s.repo.UpdatePostRepo(ctx, tx, new)
	if err != nil {
		helper.ErrMsg(err, "failed to while update post: ")
		return nil, err
	}
	return res, nil
}
func(s *PostService)DeletePostService(ctx context.Context, getPost *model.GetPostInput)error{
    tx, err := s.db.Begin(ctx)
    if err != nil {
        helper.ErrMsg(err, "failed to begin transaction: ")
        return err
    }
    defer helper.CommitOrRollback(ctx, tx)
    err = s.repo.DeletePostRepo(ctx, tx ,getPost)
    if err != nil {
        helper.ErrMsg(err, "failed to delete post: ")
        return err
    }
    return nil
}