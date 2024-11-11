package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := utils.HttpResponseBody{
			Status:  http.StatusUnauthorized,
			Message: "Unable to verify the session cookie",
			Data:    nil,
		}
		cookie, err := r.Cookie("Session")
		if err != nil {
			if err == http.ErrNoCookie {
				res.UpdateHttpResponse(w)
				return
			} else {
				res.Message = "Unable to retrieving the cookie"
				res.UpdateHttpResponse(w)
				return
			}
		}
		log.Println(cookie.Value)
		tokenString := cookie.Value
		claims, err := ValidateToken(tokenString)
		if err != nil {
			res.UpdateHttpResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next(w, r.WithContext(ctx))
	})
}

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
		"iss":  "neukod-backend",
		"sub":  username,
		"role": role,
		"exp":  time.Now().Add(48 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	SECRET := []byte(env.SECRET)
	if generatedToken, err := claims.SignedString(SECRET); err != nil {
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
		SECRET := []byte(env.SECRET)
		return SECRET, nil
	}); err != nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, fmt.Errorf("invalid token")
	}
}
