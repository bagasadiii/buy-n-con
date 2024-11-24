package repository

import (
	"context"
	"errors"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ItemRepoImpl interface{
	CreateItemRepo(ctx context.Context, tx pgx.Tx, item *model.Item)error
	GetItemByIDRepo(ctx context.Context, tx pgx.Tx, input *model.GetItemInput)(*model.ItemResp, error)
	GetAllItemsRepo(ctx context.Context, tx pgx.Tx, username string)([]model.ItemResp, error)
	ItemUpdateRepo(ctx context.Context, tx pgx.Tx, input *model.UpdateItemInput, id uuid.UUID)(*model.ItemResp, error)
	ItemDeleteRepo(ctx context.Context, tx pgx.Tx, id *uuid.UUID)error
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
	helper.SuccessMsg("item created")
	return nil
}
func(r *ItemRepo)GetItemByIDRepo(ctx context.Context, tx pgx.Tx, input *model.GetItemInput)(*model.ItemResp, error){
	query := `
		SELECT item_id, belongs_to, name, quantity, price, created_at, updated_at
		FROM items
		WHERE item_id = $1 AND belongs_to = $2
	`
	var item model.ItemResp
	row := tx.QueryRow(ctx, query, input.ItemID, input.BelongsTo)
	err := row.Scan(
		&item.ItemID,
		&item.BelongsTo,
		&item.Name, 
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
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
func(r *ItemRepo)GetAllItemsRepo(ctx context.Context, tx pgx.Tx, username string)([]model.ItemResp, error){
	query := `
		SELECT item_id, belongs_to, name, quantity, price, created_at, updated_at
		FROM items
		WHERE belongs_to = $1
	`
	rows, err := tx.Query(ctx, query, username)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to fetch items (db err): ")
		return nil, err
	}
	defer rows.Close()

	var items []model.ItemResp
	for rows.Next() {
		var item model.ItemResp
		err := rows.Scan(
			&item.ItemID,
			&item.BelongsTo,
			&item.Name, 
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			helper.ErrMsg(err, "scan items err: ")
			return nil, err
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		helper.ErrMsg(err, "iteration rows err: ")
		return nil, err
	}
	return items, nil
}
func(r *ItemRepo)ItemUpdateRepo(ctx context.Context, tx pgx.Tx, input *model.UpdateItemInput, id uuid.UUID)(*model.ItemResp, error){
	query := `
		UPDATE items
		SET name = COALESCE($1, name),
			quantity = COALESCE($2, quantity),
			price = COALESCE($3, price), 
			updated_at = $4
		WHERE item_id = $5
		RETURNING item_id, belongs_to, name, quantity, price, created_at, updated_at
	`
	var updatedItem model.ItemResp
	err := tx.QueryRow(ctx, query,
		input.Name,
		input.Quantity,
		input.Price,
		input.UpdatedAt,
		id,
	).Scan(
		&updatedItem.ItemID,
		&updatedItem.BelongsTo,
		&updatedItem.Name, 
		&updatedItem.Quantity,
		&updatedItem.Price,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to update item (db err): ")
		return nil, err
	}
	return &updatedItem, nil
}
func(r *ItemRepo)ItemDeleteRepo(ctx context.Context, tx pgx.Tx, id *uuid.UUID) error {
    query := `
        DELETE FROM items
        WHERE item_id = $1
    `

    _, err := tx.Exec(ctx, query, id)
    if err != nil {
        if err == pgx.ErrNoRows {
            return errors.New("item not found")
        }
        helper.ErrMsg(err, "failed to delete item (db err): ")
        return err
    }

    return nil
}
