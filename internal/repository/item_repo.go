package repository

import (
	"context"
	"errors"
	"math"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ItemRepoImpl interface{
	CreateItemRepo(ctx context.Context, tx pgx.Tx, item *model.Item)error
	GetItemByIDRepo(ctx context.Context, tx pgx.Tx, input *model.GetItemInput)(*model.ItemResp, error)
	GetAllItemsRepo(ctx context.Context, tx pgx.Tx, page *model.ItemsPageReq)(*model.ItemsPageRes, error)
	ItemUpdateRepo(ctx context.Context, tx pgx.Tx, input *model.UpdateItemInput, id uuid.UUID)(*model.ItemResp, error)
	ItemDeleteRepo(ctx context.Context, tx pgx.Tx, id *uuid.UUID)error
}
type ItemRepo struct{}

func NewItemRepository() ItemRepoImpl {
	return &ItemRepo{}
}

func(r *ItemRepo)CreateItemRepo(ctx context.Context, tx pgx.Tx, item *model.Item)error{
	query := `
		INSERT INTO items (item_id, user_id, owner, name, quantity, price, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := tx.Exec(ctx, query, 
		item.ItemID, 
		item.UserID, 
		item.Owner,
		item.Name,
		item.Quantity,
		item.Price,
		item.Description,
		item.CreatedAt,
		item.UpdatedAt,
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
		SELECT item_id, owner, name, quantity, price, description, created_at, updated_at
		FROM items
		WHERE item_id = $1 AND owner = $2
	`
	var item model.ItemResp
	row := tx.QueryRow(ctx, query, input.ItemID, input.Owner)
	err := row.Scan(
		&item.ItemID,
		&item.Owner,
		&item.Name, 
		&item.Quantity,
		&item.Price,
		&item.Description,
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
func(r *ItemRepo)GetAllItemsRepo(ctx context.Context, tx pgx.Tx, page *model.ItemsPageReq)(*model.ItemsPageRes, error){
	count := `
		SELECT COUNT (*)
		FROM items
		WHERE owner = $1
	`
	var totalItems int 
	err := tx.QueryRow(ctx, count, page.Username).Scan(&totalItems)
	if err != nil {
		helper.ErrMsg(err, "failed to count(db err)")
		return nil, err
	}
	query := `
		SELECT item_id, owner, name, quantity, price, description, created_at, updated_at
		FROM items
		WHERE owner = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := tx.Query(ctx, query, page.Username, page.Limit, page.Offset)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to fetch items (db err): ")
		return nil, err
	}
	defer rows.Close()

	var res model.ItemsPageRes
	res.Items = []model.ItemResp{}

	for rows.Next() {
		var item model.ItemResp
		err := rows.Scan(
			&item.ItemID,
			&item.Owner,
			&item.Name, 
			&item.Quantity,
			&item.Price,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			helper.ErrMsg(err, "scan items err: ")
			return nil, err
		}
		res.Items = append(res.Items, item)
	}
	if rows.Err() != nil {
		helper.ErrMsg(rows.Err(), "iteration rows err: ")
		return nil, rows.Err()
	}
	res.TotalItems = totalItems
	res.TotalPages = int(math.Ceil(float64(res.TotalItems) / float64(page.Limit)))
	res.Current = (page.Offset/page.Limit) + 1
	res.PageSize = len(res.Items)
	return &res, nil
}
func(r *ItemRepo)ItemUpdateRepo(ctx context.Context, tx pgx.Tx, input *model.UpdateItemInput, id uuid.UUID)(*model.ItemResp, error){
	query := `
		UPDATE items
		SET name = COALESCE($1, name),
			quantity = COALESCE($2, quantity),
			price = COALESCE($3, price), 
			description = COALESCE($4, description),
			updated_at = $5
		WHERE item_id = $6
		RETURNING item_id, owner, name, quantity, price, description, created_at, updated_at
	`
	var updatedItem model.ItemResp
	err := tx.QueryRow(ctx, query,
		input.Name,
		input.Quantity,
		input.Price,
		input.Description,
		input.UpdatedAt,
		id,
	).Scan(
		&updatedItem.ItemID,
		&updatedItem.Owner,
		&updatedItem.Name, 
		&updatedItem.Quantity,
		&updatedItem.Price,
		&updatedItem.Description,
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
