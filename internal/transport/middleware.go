package transport

import (
	"GoAssignment/internal/authentication"
	"GoAssignment/internal/contextkey"
	"context"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).
			Info("handled request")
		next.ServeHTTP(w, r)
	})
}

// Key for context key
const contextKey = "createdby"

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("an unauthorized request has been made")
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("authorization header could not be parsed")
			return
		}
		claims, err := authentication.ValidateJWT(authHeaderParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("could not validate incoming token")
			return
		} else {
			log.Info(claims)
			// Extract user ID from claims
			customClaims := claims.UserID // Type assert to your Claims struct

			// Add user ID to context
			log.Info(customClaims)
			ctx := context.WithValue(r.Context(), contextkey.UserIDKey, customClaims)
			r = r.WithContext(ctx)

			// Call the next handler with the updated context
			original(w, r)
		}

	}
}
