package middleware

import (
	"context"
	"net/http"
	"role-helper/internal/models"
)

type ctxKey string

const userCtxKey ctxKey = "user"

func Auth(us models.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token := cookie.Value
				user, err := us.ValidateToken(token)
				if err == nil && user != nil {
					ctx := context.WithValue(r.Context(), userCtxKey, user)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	v := r.Context().Value(userCtxKey)
	if v == nil {
		return nil
	}
	u, _ := v.(*models.User)
	return u
}
