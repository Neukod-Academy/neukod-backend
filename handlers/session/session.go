package session

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/middleware"
	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
	"github.com/golang-jwt/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to create a session",
		Data:    nil,
	}
	loginCred, err := utils.HttpReqReader[models.UserLogin](r)
	if err != nil {
		res.Message = "Unable while reading the login request credential"
		res.UpdateHttpResponse(w)
		return
	}
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "failed to create database client"
		res.UpdateHttpResponse(w)
		return
	}
	coll := db.Client.Database("Neukod").Collection("Users")
	var stored models.User
	if err := coll.FindOne(db.Context, bson.M{"username": loginCred.Username}, options.FindOne()).Decode(&stored); err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = http.StatusNotFound
			res.Message = "Unable to find this credential or still not registered"
			res.UpdateHttpResponse(w)
			return
		}
		res.Message = "failed to create database client"
		res.UpdateHttpResponse(w)
		return
	}
	log.Println(loginCred.Password)
	log.Println(stored.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(loginCred.Password)); err != nil {
		res.Message = "failed while validating the user password"
		res.UpdateHttpResponse(w)
		return
	}
	cookie := &http.Cookie{}

	if cookie.Value == "" {
		newToken, err := middleware.CreateToken(loginCred.Username)
		if err != nil {
			res.Message = "failed to create session cookie"
			res.UpdateHttpResponse(w)
			return
		}
		cookie = &http.Cookie{
			Name:     "session_id",
			Value:    newToken,
			Path:     "/",
			Expires:  time.Now().Add(1 * time.Minute),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	res.Status = http.StatusCreated
	res.Message = "Succesful to create a session"
	res.UpdateHttpResponse(w)
}

func DropSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to create a session",
		Data:    nil,
	}
	_, err := r.Cookie("session_id")
	if err != nil {
		res.Message = "No active session found"
		res.UpdateHttpResponse(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	res.Status = http.StatusOK
	res.Message = "Succesful to drop a session"
	res.UpdateHttpResponse(w)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to process this request",
		Data:    nil,
	}
	newUser, err := utils.HttpReqReader[models.User](r)
	if err != nil {
		res.UpdateHttpResponse(w)
		return
	}
	AvailRole := map[string]struct{}{
		"admin":       {},
		"super_admin": {},
	}
	if _, exist := AvailRole[newUser.Role]; !exist {
		res.UpdateHttpResponse(w)
		return
	}
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "failed to create database client"
		res.UpdateHttpResponse(w)
		return
	}
	coll := db.Client.Database("Neukod").Collection("Users")
	var stored models.User
	availableUsername := false
	if err := coll.FindOne(db.Context, bson.M{
		"username": newUser.Username,
	}, options.FindOne()).Decode(&stored); err == nil {
		res.Message = "This username is not available"
		res.UpdateHttpResponse(w)
		return
	} else if err == mongo.ErrNoDocuments {
		availableUsername = true
	}
	if availableUsername {
		if newNanoId, err := gonanoid.New(10); err != nil {
			res.Message = "Failed while generating the new user id"
			res.UpdateHttpResponse(w)
			return
		} else {
			newUser.Id = newNanoId
		}
		hashingPassed := false
		timedOut := 0
		for !hashingPassed {
			if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost); err != nil {
				res.Message = "Failed while hashing the credential"
				res.UpdateHttpResponse(w)
				return
			} else {
				if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(newUser.Password)); err != nil {
					fmt.Println("hashed result and the password is not matching, trying again..")
					timedOut++
					if timedOut >= 5 {
						res.Message = "Failed while hashing the credential"
						res.UpdateHttpResponse(w)
						return
					}
					continue
				}
				newUser.Password = string(hashedPassword)
				hashingPassed = true
			}
		}
		newUser.CreatedAt = time.Now()
		newUser.UpdatedAt = time.Now()
		if _, err := coll.InsertOne(db.Context, newUser); err != nil {
			res.Message = "Failed while adding the new data"
			res.UpdateHttpResponse(w)
			return
		}
	}
	res.Status = http.StatusCreated
	res.Message = "Success to add a new user"
	res.Data = newUser
	res.UpdateHttpResponse(w)
}

func ShowAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	AllowedRole := map[string]struct{}{
		"admin":       {},
		"super_admin": {},
	}
	userData, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed conversion the claims", http.StatusInternalServerError)
		return
	}
	role, ok := userData["role"].(string)
	fmt.Println(role)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else if _, allowed := AllowedRole[role]; !allowed {
		http.Error(w, "Unauthorized: this action is not belong to the authorized role", http.StatusUnauthorized)
		return
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "error to show all of the account",
		Data:    nil,
	}
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "Failed to create the database client"
		res.UpdateHttpResponse(w)
		return
	}
	coll := db.Client.Database("Neukod").Collection("Users")
	if cursor, err := coll.Find(db.Context, bson.M{}, options.Find()); err != nil {
		res.Message = "Failed to find the data in the database"
		res.UpdateHttpResponse(w)
		return
	} else {
		stored := []models.User{}
		if err := cursor.All(db.Context, &stored); err != nil {
			res.Message = "Failed to serve the data"
			res.UpdateHttpResponse(w)
			return
		}
		res.Status = http.StatusOK
		res.Message = "Success to get all of the users data"
		res.Data = stored
		res.UpdateHttpResponse(w)
	}
}

func RemoveAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	userData, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed conversion the claims", http.StatusInternalServerError)
		return
	}
	role, ok := userData["role"].(string)
	AllowedRole := map[string]struct{}{
		"admin":       {},
		"super_admin": {},
	}
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else if _, allowed := AllowedRole[role]; !allowed {
		http.Error(w, "Unauthorized: this action is not belong to the authorized role", http.StatusUnauthorized)
		return
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "failed to delete an account",
		Data:    nil,
	}
	userId := r.PathValue("id")
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "Failed to create the database client"
		res.UpdateHttpResponse(w)
		return
	}
	defer db.CloseClientDB()

	coll := db.Client.Database("Neukod").Collection("Users")
	tempColl := db.Client.Database("Neukod").Collection("Users_temp")

	var stored models.User
	filter := bson.M{
		"user_id": userId,
	}
	err := coll.FindOne(db.Context, filter).Decode(&stored)
	if err != nil {
		res.Message = "ERR: Failed to find the user_id for deletion"
		res.UpdateHttpResponse(w)
		return
	}
	stored.RemovedAt = time.Now()
	stored.UpdatedAt = time.Now()
	if _, err := tempColl.InsertOne(db.Context, stored); err != nil {
		res.Message = "ERR: Failed to move the selected user to temporary database for deletion"
		res.UpdateHttpResponse(w)
		return
	}
	if _, err := coll.DeleteOne(db.Context, filter); err != nil {
		res.Message = "ERR: Failed while deleting the data"
		res.UpdateHttpResponse(w)
		return
	}
	if userData["sub"] == stored.Username {
		_, err := r.Cookie("session_id")
		if err != nil {
			res.Message = "No active session found"
			res.UpdateHttpResponse(w)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Session",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		})
	}
	res.Status = http.StatusOK
	res.Message = "Successful to delete a data"
	res.Data = stored.Username
	res.UpdateHttpResponse(w)
}
