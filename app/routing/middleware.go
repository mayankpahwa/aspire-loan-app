package routing

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mayankpahwa/aspire-loan-app/internal/repo/mysql"
	"github.com/mayankpahwa/aspire-loan-app/internal/resources/types"
	"github.com/pkg/errors"
)

// ContextKey is string alias for context type.
type ContextKey string

const ContextUserIDKey ContextKey = "user_id"

func AuthMiddleware(repo mysql.Repo) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, password, ok := r.BasicAuth()
			if !ok {
				WriteError(r, w, errors.Wrap(types.ErrUnauthorized, "credentials not found in request header"))
				return
			}
			userFromDB, err := repo.GetUserByID(r.Context(), user)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					WriteError(r, w, errors.Wrap(types.ErrUnauthorized, "user not found"))
					return
				}
				WriteError(r, w, err)
				return
			}
			if userFromDB.Password != password {
				WriteError(r, w, errors.Wrap(types.ErrUnauthorized, "incorrect user_id / password combination"))
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserIDKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidationUserIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDPath := chi.URLParam(r, "userID")
		userIDAuth := r.Context().Value(ContextUserIDKey)
		if userIDPath != userIDAuth {
			errors.Wrap(types.ErrUnauthorized, "cannot access resources of a different user")
			return
		}
		next.ServeHTTP(w, r)
	})
}
