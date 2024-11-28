package service

import (
	"context"
	"errors"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemServiceImpl interface {
	CreateItemService(ctx context.Context, new *model.CreateItemInput)(*model.Item, error)
	GetItemByIDService(ctx context.Context, input *model.GetItemInput)(*model.ItemResp, error)
	GetAllItemsService(ctx context.Context, page *model.ItemsPageReq)(*model.ItemsPageRes, error)
	UpdateItemService(ctx context.Context, new *model.UpdateItemInput, getItem *model.GetItemInput)(*model.ItemResp, error)
	DeleteItemService(ctx context.Context, id *uuid.UUID)error
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
func(s *ItemService)GetItemByIDService(ctx context.Context, input *model.GetItemInput)(*model.ItemResp, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to get item(tx error): ")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)
	item, err := s.repo.GetItemByIDRepo(ctx, tx, input)
	if err != nil {
		helper.ErrMsg(err, "failed to get item(db error): ")
		return nil, err
	}
	return item, nil
}
func(s *ItemService)GetAllItemsService(ctx context.Context, page *model.ItemsPageReq)(*model.ItemsPageRes, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transaction: ")
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
	res, err := s.repo.GetAllItemsRepo(ctx, tx, page)
	if err != nil {
		helper.ErrMsg(err, "failed to get itemtransaction: ")
		return nil, err
	}
	return res, nil
}
func(s *ItemService)UpdateItemService(ctx context.Context, new *model.UpdateItemInput, getItem *model.GetItemInput)(*model.ItemResp, error){
	tx, err := s.db.Begin(ctx)
	if err != nil {
		helper.ErrMsg(err, "failed to begin transaction: ")
		return nil, err
	}
	defer helper.CommitOrRollback(ctx, tx)
	existingItem, err := s.repo.GetItemByIDRepo(ctx, tx, getItem)
	if err != nil {
		helper.ErrMsg(err, "failed to get item: ")
		return nil, err
	}
	if new.Name == "" {
		new.Name = existingItem.Name
	}
	if new.Price == 0 {
		new.Price = existingItem.Price
	}
	if new.Quantity == 0 {
		new.Quantity = existingItem.Quantity
	}
	updateItem := model.UpdateItemInput{
		Name: new.Name,
		Quantity: new.Quantity,
		Price: new.Price,
		Description: new.Description,
		UpdatedAt: time.Now(),
	}
	id := getItem.ItemID
	res, err := s.repo.ItemUpdateRepo(ctx, tx, &updateItem, id)
	if err != nil {
		helper.ErrMsg(err, "failed to while update item: ")
		return nil, err
	}
	return res, nil
}
func(s *ItemService)DeleteItemService(ctx context.Context, id *uuid.UUID)error{
    tx, err := s.db.Begin(ctx)
    if err != nil {
        helper.ErrMsg(err, "failed to begin transaction: ")
        return err
    }
    defer helper.CommitOrRollback(ctx, tx)
    err = s.repo.ItemDeleteRepo(ctx, tx, id)
    if err != nil {
        helper.ErrMsg(err, "failed to delete item: ")
        return err
    }
    return nil
}