package repository

import (
	"context"
	"errors"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/jackc/pgx/v5"
)

type ItemRepoImpl interface{
	CreateItemRepo(ctx context.Context, tx pgx.Tx, item *model.Item)error
	GetItemRepo(ctx context.Context, tx pgx.Tx, itemID string)(*model.Item, error)
}
type ItemRepo struct{}

func NewItemRepository() ItemRepoImpl {
	return &ItemRepo{}
}

func(r *ItemRepo)CreateItemRepo(ctx context.Context, tx pgx.Tx, item *model.Item)error{
	query := `
		INSERT INTO items (item_id, user_id, belongs_to, name, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := tx.Exec(
		ctx, query, item.ItemID, item.UserID, item.BelongsTo,
		item.Name,item.Quantity,item.Price, item.CreatedAt,item.UpdatedAt,
	)
	if err != nil {
		helper.ErrMsg(err, "failed to create item: ")
		return err
	}
	return nil
}
func(r *ItemRepo)GetItemRepo(ctx context.Context, tx pgx.Tx, itemID string)(*model.Item, error){
	query := `
		SELECT item_id, user_id, belongs_to, name, quantity, price, created_at, updated_at
		FROM items
		WHERE item_id = $1
	`
	var item model.Item
	row := tx.QueryRow(ctx, query, itemID)
	err := row.Scan(
		&item.ItemID, &item.UserID, &item.BelongsTo, &item.Name, 
		&item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to fetch item (db error): ")
		return nil, err
	}
	return &item, nil
}
// func(r *ItemRepo)GetAllItemsRepo(ctx context.Context, tx pgx.Tx, userID string)([]model.Item, error){
// 	query := `
// 		SELECT item_id, name, quantity, price, created_at, updated_at
		
// 	`
// 	var items []model.Item
// }