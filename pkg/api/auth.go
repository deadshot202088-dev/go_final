package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 8 * time.Hour

type signInRequest struct {
	Password string `json:"password"`
}

// passwordClaims хранит хэш пароля так при смене
// TODO_PASSWORD все выданные ранее токены перестают быть валидными
type passwordClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func jwtSecret(password string) []byte {
	return []byte("todo-scheduler:" + password)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" || req.Password != storedPassword {
		writeError(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	claims := passwordClaims{
		PasswordHash: passwordHash(storedPassword),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret(storedPassword))
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": signed})
}

//  middleware, требующий валидный JWT-токен в куке token, если задан
// пароль в переменной окружения TODO_PASSWORD. Если пароль не задан, доступ
// к API открыт без аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storedPassword := os.Getenv("TODO_PASSWORD")
		if storedPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims := &passwordClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(*jwt.Token) (any, error) {
			return jwtSecret(storedPassword), nil
		})
		if err != nil || !token.Valid || claims.PasswordHash != passwordHash(storedPassword) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
