package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils/responses"
	"gorm.io/gorm"
)

func AuthMiddleware(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			responses.Error(w, http.StatusUnauthorized, "Authenticated token is not set")
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userID, err := models.GetUserIDByAPIToken(db, token)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, "Authenticated token is invalid")
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		next.ServeHTTP(w, r)
	})
}
