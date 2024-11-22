package service

import (
	"context"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl interface {
	RegisterService(ctx context.Context, new *model.RegisterInput)(*model.User, error)
	LoginService(ctx context.Context, new *model.LoginInput)(string, error)
}
type UserService struct {
	repo repository.UserRepoImpl
}
func NewUserService(repo repository.UserRepoImpl)UserServiceImpl{
	return &UserService{repo:repo}
}
func(s *UserService)RegisterService(ctx context.Context, new *model.RegisterInput)(*model.User, error){
	user, err := model.NewUser(new)
	if err != nil {
		helper.ErrMsg(err, "failed to create user: ")
		return nil, err
	}
	if err := s.repo.RegisterRepo(ctx, user); err != nil {
		helper.ErrMsg(err, "failed to create user (database error): ")
		return nil, err
	}
	helper.SuccessMsg("user created")
	return user, nil
}
func(s *UserService)LoginService(ctx context.Context, new *model.LoginInput)(string, error){
	user, err := s.repo.LoginRepo(ctx, new)
	if err != nil {
		helper.ErrMsg(err, "failed: ")
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(new.Password)); err != nil {
		helper.ErrMsg(err, "invalid password: ")
		return "", err
	}
	token, err := middleware.GenerateToken(user.UserID, new.Username)
	if err != nil {
		helper.ErrMsg(err, "failed to create token: ")
		return "", err
	}
	return token, nil
}