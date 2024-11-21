package repository

import (
	"context"
	"errors"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoImpl interface {
	RegisterRepo(ctx context.Context,user *model.User)error
	LoginRepo(ctx context.Context, input *model.LoginInput)(*model.User, error)
}
type UserRepository struct {
	db *pgxpool.Pool
}
func NewUserRepository(db *pgxpool.Pool)UserRepoImpl{
	return &UserRepository{db:db}
}
func(r *UserRepository)RegisterRepo(ctx context.Context,user *model.User)error{
	query := `
		INSERT INTO users (user_id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := r.db.Exec(ctx, query, user.UserID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return helper.ErrMsg(err, "failed to create user: ")
	}

	helper.SuccessMsg("user created")
	return nil
}
func(r *UserRepository)LoginRepo(ctx context.Context, input *model.LoginInput)(*model.User, error){
	query := `
		SELECT user_id, password, email
		FROM users
		WHERE username = $1
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, input.Username).Scan(&user.UserID, &user.Password, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, errors.New("no data found")
		}
		return nil, helper.ErrMsg(err, "failed to find data: ")
	}
	helper.SuccessMsg("login successful")
	return &user, nil
}