package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/google/uuid"
)

type Item struct {
	ItemID    		uuid.UUID		`json:"item_id"`
	UserID    		uuid.UUID		`json:"user_id"`
	Owner	 		string			`json:"owner"`
	Name      		string			`json:"name"`
	Quantity  		int				`json:"quantity"`
	Price     		int				`json:"price"`
	Description		string			`json:"description"`
	CreatedAt 		time.Time		`json:"created_at"`
	UpdatedAt 		time.Time		`json:"updated_at"`
}
type CreateItemInput struct {
	ItemID			uuid.UUID		`json:"item_id"`
	Name      		string			`json:"name" validate:"required"`
	Quantity  		int				`json:"quantity" validate:"required,gt=0"`
	Price     		int				`json:"price" validate:"required,gt=0"`
	Description		string			`json:"description"`
}
type GetItemInput struct {
	ItemID			uuid.UUID		`json:"item_id"`
	Owner			string			`json:"owner"`
}
type UpdateItemInput struct {
	Name      string    `json:"name" validate:"required,min=3"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	Price     int       `json:"price" validate:"required,gt=0"`
	UpdatedAt time.Time `json:"updated_at"`
	Description	string	`json:"description"`
}
type ItemResp struct {
	ItemID    uuid.UUID `json:"item_id"`
	Owner	string	`json:"owner"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	Description	string		`json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ItemsPageReq struct {
	Username	string	`json:"username"`
	Limit		int		`json:"limit"`
	Offset		int		`json:"offset"`
}
type ItemsPageRes struct {
	Items		[]ItemResp			`json:"items"`
	TotalItems	int					`json:"total_items"`
	TotalPages 	int					`json:"total_pages"`
	Current		int					`json:"current"`
	PageSize	int					`json:"page_size"`
}
func NewItem(ctx context.Context, input *CreateItemInput)(*Item, error){
	name := strings.TrimSpace(input.Name)
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
		Owner: ctxKey.UsernameKey,
		Name: name,
		Quantity: input.Quantity,
		Price: input.Price,
		Description: input.Description,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}