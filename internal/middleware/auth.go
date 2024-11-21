package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bagasadiii/buy-n-con/helper"
)

type contextKey string

var (
	UserIDKey = contextKey("user_id")
	UsernameKey = contextKey("username")
)

func Auth(next http.HandlerFunc)http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			res := helper.Response{
				Status: http.StatusUnauthorized,
				Message: "Missing auth header",
				Data: nil,
				Err: nil,
			}
			helper.JSONResponse(w, res.Status, res)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			res := helper.Response{
				Status:  http.StatusUnauthorized,
				Message: "Invalid Authorization header format",
				Data:    nil,
				Err:     nil,
			}
			helper.JSONResponse(w, res.Status, res)
			return
		}
		validation := ValidateToken(token)
		if validation.Err != nil{
			res := helper.Response{
				Status: http.StatusUnauthorized,
				Message: "Invalid or expired token",
				Data: validation,
				Err: validation.Err,
			}
			helper.JSONResponse(w, res.Status, res)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, validation.ID)
		ctx = context.WithValue(ctx, UsernameKey, validation.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}