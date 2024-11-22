package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	ID       uuid.UUID
	Username string
	jwt.RegisteredClaims
}
type UserValidation struct {
	Token		string
	ID			uuid.UUID
	Username 	string
	Err 		error
}
func GenerateToken(id uuid.UUID, username string) (string, error) {
	exp := time.Now().Add(24 * time.Hour)
	newClaims := &Claims {
		ID: id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,newClaims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		helper.ErrMsg(err, "failed to generate token")
		return "", err
	}
	helper.SuccessMsg("create token success")
	return tokenString, nil
}

func ValidateToken(tokenString string)*UserValidation{
	newClaims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, newClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		helper.ErrMsg(err, "failed to validate token")
		return &UserValidation{
			Token: tokenString,
			Err: err,
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &UserValidation{
			Token: tokenString,
			Username: claims.Username,
			ID: claims.ID,
		}
	}
	return &UserValidation{
		Token: tokenString,
		Err:   errors.New("invalid token"),
	}
}