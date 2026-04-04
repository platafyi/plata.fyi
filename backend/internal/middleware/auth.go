package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/platafyi/plata.fyi/internal/database"
)

type contextKey string

const ownerIDKey contextKey = "owner_id"

func Auth(store database.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				http.Error(w, `{"error":"Потребна е автентикација"}`, http.StatusUnauthorized)
				return
			}

			ownerID, err := store.GetOwnerByToken(r.Context(), token)
			if err != nil {
				http.Error(w, `{"error":"Грешка при автентикација"}`, http.StatusInternalServerError)
				return
			}
			if ownerID == nil {
				http.Error(w, `{"error":"Неважечка или истечена сесија"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ownerIDKey, *ownerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetOwnerID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ownerIDKey).(string)
	return id, ok
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
