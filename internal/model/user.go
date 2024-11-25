package model

import (
	"errors"
	"time"
	"unicode"

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
	Username string		`json:"username" validate:"required,min=3"`
	Email    string		`json:"email" validate:"required,email"`
	Password string		`json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Username string		`json:"username" validate:"required,min=3"`
	Password string		`json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UserID		uuid.UUID	`json:"user_id"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}
func NewUser(input *RegisterInput)(*User, error){
	for _, char := range input.Username {
		if unicode.IsUpper(char) {
			return nil, errors.New("username cannot contain uppercase")
		}
		
		if unicode.IsSpace(char) {
			return nil, errors.New("username cannot contain spaces")
		}
		
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_') {
			return nil, errors.New("username can only contain letters, digits, and underscores")
		}
	}

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