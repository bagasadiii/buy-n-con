package model

import (
	"context"
	"errors"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/google/uuid"
)

type Post struct {
	PostID		uuid.UUID		`json:"post_id"`
	Content		string			`json:"content"`
	Owner		string			`json:"owner"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	UserID		uuid.UUID		`json:"user_id"`
}
type PostInput struct {
	Content		string			`json:"content" validate:"required"`
}
type GetPostInput struct {
	PostID		uuid.UUID		`json:"post_id"`
	Owner		string			`json:"owner"`
}
type UpdatePostInput struct {
	PostID		uuid.UUID		`json:"post_id"`
	Owner		string			`json:"owner"`
	Content		string			`json:"content" validate:"required"`
	UpdatedAt	time.Time		`json:"updated_at"`
}
type PostsPageReq struct {
	Username	string			`json:"username"`
	Limit		int				`json:"limit"`
	Offset		int				`json:"offset"`
}
type PostsPageRes struct {
	Posts		[]Post			`json:"posts"`
	TotalPosts	int				`json:"total_posts"`
	TotalPages 	int				`json:"total_pages"`
	Current		int				`json:"current"`
	PageSize	int				`json:"page_size"`
}
func NewPost (ctx context.Context, input *PostInput)(*Post, error){
	ctxKey, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
	if !ok {
		helper.ErrMsg(nil, "failed to get context key: ")
		return nil, errors.New("failed to get context key")
	}
	if ctxKey.UserIDKey == uuid.Nil || ctxKey.UsernameKey == ""{
		helper.ErrMsg(nil, "no data in context")
		return nil, errors.New("no data in context")
	}
	if input.Content == "" {
		return nil, errors.New("content cannot be empty")
	}
	return &Post{
		PostID: uuid.New(),
		Content: input.Content,
		Owner: ctxKey.UsernameKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: ctxKey.UserIDKey,
	}, nil
}