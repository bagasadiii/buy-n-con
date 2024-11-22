package model

import (
	"context"
	"errors"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/google/uuid"
)

type Item struct {
	ItemID    uuid.UUID `json:"item_id"`
	UserID    uuid.UUID `json:"user_id"`
	BelongsTo string	`json:"belongs_to"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateItemInput struct {
	ItemID   	uuid.UUID	`json:"item_id"`
	Name      	string		`json:"name"`
	Quantity  	int			`json:"quantity"`
	Price     	int			`json:"price"`
}
func NewItem(ctx context.Context, input *CreateItemInput)(*Item, error){
	ctxKey, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
	if !ok {
		helper.ErrMsg(nil, "failed to get context key: ")
		return nil, errors.New("failed to get context key")
	}
	if ctxKey.UserIDKey == uuid.Nil || ctxKey.UsernameKey == ""{
		helper.ErrMsg(nil, "no data in context")
		return nil, errors.New("no data in context")
	}
	return &Item{
		ItemID: uuid.New(),
		UserID: ctxKey.UserIDKey,
		BelongsTo: ctxKey.UsernameKey,
		Name: input.Name,
		Quantity: input.Quantity,
		Price: input.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}