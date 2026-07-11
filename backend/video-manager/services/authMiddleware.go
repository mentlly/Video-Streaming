package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Fetch the secret key directly from your Go environment (.env)
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		// 2. Grab the cookie from the request
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims := &CustomClaims{}

		debugClaims := jwt.MapClaims{}
		_, _, debugErr := new(jwt.Parser).ParseUnverified(tokenString, &debugClaims)
		if debugErr == nil {
			fmt.Printf("\n[DEBUG] Raw Node Payload Map: %+v\n", debugClaims)
		} else {
			fmt.Printf("\n[DEBUG] Could not even read unverified token: %v\n", debugErr)
		}

		// 3. Parse and Validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			// Always validate the signing method is what you expect (HMAC-SHA256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		fmt.Printf("Token: %v\n", claims.UserID)

		// 4. Handle invalid signatures, expired tokens, or parsing errors
		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// 2. Add the UserID to the request context
		ctx := context.WithValue(r.Context(), "userId", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
