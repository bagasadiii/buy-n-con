package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/google/uuid"
	router "github.com/julienschmidt/httprouter"
)

type ContextKey struct {
	UserIDKey uuid.UUID
	UsernameKey string
}
type ctxKey string
const UserContextKey = ctxKey("context_key")

func Auth(next router.Handle)router.Handle{
	return func(w http.ResponseWriter, r *http.Request, p router.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			res := helper.UnauthorizedErr("Missing Header: ", nil)
			helper.JSONResponse(w, res.Status, res)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			res := helper.UnauthorizedErr("Missing token: ", nil)
			helper.JSONResponse(w, res.Status, res)
			return
		}
		validation := ValidateToken(token)
		if validation.Err != nil{
			res := helper.UnauthorizedErr("Expired or invalid token ", validation.Err)
			helper.JSONResponse(w, res.Status, res)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, &ContextKey{
			UserIDKey: validation.ID,
			UsernameKey: validation.Username,
		})
		next(w, r.WithContext(ctx), p)
	}
}