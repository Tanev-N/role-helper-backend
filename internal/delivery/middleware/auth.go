package middleware

import (
	"context"
	"log"
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
				log.Println(cookie.Value)
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
		log.Println("AAAAAAA")
		return nil
	}
	u, _ := v.(*models.User)
	log.Println(u.Username, u.ID)
	return u
}
