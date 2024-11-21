package model

import (
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID		uuid.UUID	`json:"user_id"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type RegisterInput struct {
	Username string		`json:"username"`
	Email    string		`json:"email"`
	Password string		`json:"password"`
}

type LoginInput struct {
	Username string		`json:"username"`
	Password string		`json:"password"`
}

func NewUser(input *RegisterInput)(*User, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.ErrMsg(err, "failed to hash password: ")
		return nil, err
	}
	return &User{
		UserID: uuid.New(),
		Username: input.Username,
		Email: input.Email,
		Password: string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}