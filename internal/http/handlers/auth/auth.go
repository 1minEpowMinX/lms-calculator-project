package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TTL токена
var tokenExpire = time.Now().Add(5 * time.Minute)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInteractor interface {
	RegisterUser(ctx context.Context, uname, pswrd string) error
	LoginUser(ctx context.Context, uname, pswrd string) (string, error)
}

// Обработчик регистрации
func RegisterUserHandler(ctx context.Context, userInteractor UserInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error: %v", err)
			return
		}

		if err := userInteractor.RegisterUser(ctx, req.Username, req.Password); err != nil {
			http.Error(w, "This username is already registered", http.StatusBadRequest)
			log.Printf("Error: %v", err)
			return
		}

		log.Print("Success RegisterUserHandler")
		w.WriteHeader(http.StatusCreated)
	}
}

// Обработчик аутентефикации
func LoginUserHandler(ctx context.Context, userInteractor UserInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error: %v", err)
			return
		}

		token, err := userInteractor.LoginUser(ctx, req.Username, req.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Printf("Error: %v", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			Expires:  tokenExpire,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
		log.Print("Success LoginUserHandler")
	}
}
