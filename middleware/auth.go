package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getRole(username string) (string, error) {
	db := new(utils.Mongo)
	var user models.User
	db.CreateClient(env.MONGO_URI)
	coll := db.Client.Database("Neukod").Collection("Users")
	if err := coll.FindOne(db.Context, bson.M{"username": username}, options.FindOne()).Decode(&user); err != nil {
		return "", err
	}
	return user.Username, nil
}

func CreateToken(username string) (string, error) {
	role, err := getRole(username)
	if err != nil {
		return "", err
	}
	mapClaims := jwt.MapClaims{
		"iss": "neukod-backend",
		"sub": username,
		"aud": role,
		"exp": time.Now().Add(48 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	if generatedToken, err := claims.SignedString(env.SECRET); err != nil {
		return "", err
	} else {
		return generatedToken, nil
	}
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return env.SECRET, nil
	}); err != nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, fmt.Errorf("invalid token")
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Cookie not found", http.StatusUnauthorized)
					return
				} else {
					http.Error(w, "Unable to retrieving the cookie", http.StatusUnauthorized)
					return
				}
			}
			tokenString := cookie.Value
			claims, err := ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}
