package user

import (
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewTrial(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to execute this request",
		Data:    nil,
	}
	var newTrial models.Trial
	newTrial, err := utils.HttpReqReader[models.Trial](r)
	if err != nil {
		return err
	}
	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.UpdateHttpResponse(w, http.StatusInternalServerError, "Unable to connect to the database")
		return err
	}
	defer db.CloseClientDB()
	col := db.Client.Database("Neukod").Collection("Trial")
	newTrial.TrialId = uuid.New().String()
	newTrial.CreatedAt = time.Now()
	newTrial.UpdatedAt = time.Now()
	if _, err := col.InsertOne(db.Context, newTrial, options.InsertOne()); err != nil {
		res.UpdateHttpResponse(w, http.StatusInternalServerError, "Failed to add the data to the database")
		return err
	}
	res.UpdateHttpResponse(w, http.StatusCreated, newTrial)
	return nil
}

func ConfirmTrial(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	return nil
}

func EditTrial(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "PUT" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	return nil
}

func DeleteTrial(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "DELETE" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	return nil
}
