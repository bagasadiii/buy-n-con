package repository

import (
	"context"
	"errors"
	"math"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/jackc/pgx/v5"
)

type PostRepoImpl interface {
	CreatePostRepo(ctx context.Context, tx pgx.Tx, new *model.Post) error
    GetPostByIDRepo(ctx context.Context, tx pgx.Tx, data *model.GetPostInput) (*model.Post, error)
    GetAllPostRepo(ctx context.Context, tx pgx.Tx, page *model.PostsPageReq)(*model.PostsPageRes, error)
    UpdatePostRepo(ctx context.Context, tx pgx.Tx, post *model.UpdatePostInput) (*model.Post, error)
    DeletePostRepo(ctx context.Context, tx pgx.Tx, post *model.GetPostInput) error
}
type PostRepo struct{}

func NewPostRepository()PostRepoImpl {
	return &PostRepo{}
}
func (r *PostRepo) CreatePostRepo(ctx context.Context, tx pgx.Tx, new *model.Post)error{
	query := `
		INSERT INTO posts (post_id, content, owner, created_at, updated_at, user_id)
		VALUES ($1, $2, $3, $4, $5,$6)
	`
	_, err := tx.Exec(ctx, query,
		new.PostID,
		new.Content,
		new.Owner,
		new.CreatedAt,
		new.UpdatedAt,
		new.UserID,
	)
	if err != nil {
		helper.ErrMsg(err, "failed to exec command")
		return err
	}
	return nil
}
func(r *PostRepo) GetPostByIDRepo(ctx context.Context, tx pgx.Tx, data *model.GetPostInput)(*model.Post, error){
	query := `
		SELECT post_id, content, owner, created_at, updated_at, user_id
		FROM posts
		WHERE post_id = $1 AND owner = $2
	`
	var post model.Post
	row := tx.QueryRow(ctx, query, data.PostID, data.Owner)
	err := row.Scan(
		&post.PostID,
		&post.Content,
		&post.Owner,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
	)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to fetch post (db error): ")
		return nil, err
	}
	return &post, nil
}
func(r *PostRepo)GetAllPostRepo(ctx context.Context, tx pgx.Tx, page *model.PostsPageReq)(*model.PostsPageRes, error){
	count := `
		SELECT COUNT (*)
		FROM posts
		WHERE owner = $1
	`
	var totalPosts int 
	err := tx.QueryRow(ctx, count, page.Username).Scan(&totalPosts)
	if err != nil {
		helper.ErrMsg(err, "failed to count posts (db err)")
		return nil, err
	}
	query := `
		SELECT post_id, owner, content, created_at, updated_at, user_id
		FROM posts
		WHERE owner = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := tx.Query(ctx, query, page.Username, page.Limit, page.Offset)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to fetch post (db error): ")
		return nil, err
	}
	defer rows.Close()

	var res model.PostsPageRes
	res.Posts = []model.Post{}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.PostID,
			&post.Owner,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UserID,
		)
		if err != nil {
			helper.ErrMsg(err, "scan posts err: ")
			return nil, err
		}
		res.Posts = append(res.Posts, post)
	}
	if rows.Err() != nil {
		helper.ErrMsg(rows.Err(), "iteration rows err: ")
		return nil, err
	}
	res.TotalPosts = totalPosts
	res.TotalPages = int(math.Ceil(float64(res.TotalPosts) / float64(page.Limit)))
	res.Current = (page.Offset/page.Limit) + 1
	res.PageSize = len(res.Posts)
	return &res, nil
}
func(r *PostRepo)UpdatePostRepo(ctx context.Context, tx pgx.Tx, post *model.UpdatePostInput)(*model.Post, error){
	query := `
		UPDATE posts
		SET content = $1,
			updated_at = $2
		WHERE post_id = $3 AND owner = $4
		RETURNING post_id, owner, content, created_at, updated_at
	`
	var updatedPost model.Post
	err := tx.QueryRow(ctx, query,
		post.Content,
		post.UpdatedAt,
		post.PostID,
		post.Owner,
	).Scan(
		&updatedPost.PostID,
		&updatedPost.Owner,
		&updatedPost.Content,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		helper.ErrMsg(err, "failed to update post (db err): ")
		return nil, err
	}
	return &updatedPost, nil
}
func(r *PostRepo)DeletePostRepo(ctx context.Context, tx pgx.Tx, post *model.GetPostInput)error{
	query := `
		DELETE FROM posts
        WHERE post_id = $1 AND owner = $2
	`
	_, err := tx.Exec(ctx, query, post.PostID, post.Owner)
    if err != nil {
        if err == pgx.ErrNoRows {
            return errors.New("item not found")
        }
        helper.ErrMsg(err, "failed to delete post (db err): ")
        return err
    }
    return nil
}