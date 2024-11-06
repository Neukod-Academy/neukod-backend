package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
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
	res.UpdateHttpResponse(w)
}

func DropSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to drop this session",
		Data:    nil,
	}
	res.UpdateHttpResponse(w)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "Failed to create database client"
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
	res.Status = http.StatusOK
	res.Message = "Success to add a new user"
	res.Data = newUser
	res.UpdateHttpResponse(w)
}

func ShowAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
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
