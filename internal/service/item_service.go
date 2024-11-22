package service

import (
	"context"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemServiceImpl interface {
	CreateItemService(ctx context.Context, new *model.CreateItemInput)(*model.Item, error)
}
type ItemService struct {
	repo repository.ItemRepoImpl
	db *pgxpool.Pool
}
func NewItemService(repo repository.ItemRepoImpl, db *pgxpool.Pool)ItemServiceImpl{
	return &ItemService{
		repo:repo,
		db:db,
	}
}
func(s *ItemService)CreateItemService(ctx context.Context, new *model.CreateItemInput)(*model.Item, error){
	item, err := model.NewItem(ctx, new)
	if err != nil {
		helper.ErrMsg(err, "failed to create item: ")
		return nil, err
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transaction: ")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)
	if err := s.repo.CreateItemRepo(ctx, tx, item); err != nil {
		helper.ErrMsg(err, "failed to create item(db err): ")
		return nil, err
	}
	return item, nil
}