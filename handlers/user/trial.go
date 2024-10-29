package user

import (
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	res := utils.HttpResponseBody{
		Status:  http.StatusCreated,
		Message: "Trial has been booked, let us reach you directly for the update",
		Data:    nil,
	}
	var newTrial models.Trial

	newTrial, err := utils.HttpReqReader[models.Trial](r)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to add a new trial session"
		res.UpdateHttpResponse(w)
		return
	}

	isAllowed := false
	for _, course := range models.TrialList {
		if newTrial.Course == course {
			isAllowed = true
			continue
		}
	}

	if !isAllowed {
		res.Status = http.StatusBadRequest
		res.Message = "This course is still unavailable"
		res.UpdateHttpResponse(w)
		return
	}

	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to create a client to the database"
		res.UpdateHttpResponse(w)
		return
	}

	defer db.CloseClientDB()
	col := db.Client.Database("Neukod").Collection("Trial")
	newTrial.TrialId = uuid.New().String()
	newTrial.CreatedAt = time.Now()
	newTrial.UpdatedAt = time.Now()
	if _, err := col.InsertOne(db.Context, newTrial, options.InsertOne()); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to add a new trial session"
		res.UpdateHttpResponse(w)
		return
	}

	res.Data = newTrial
	res.UpdateHttpResponse(w)
}

func ConfirmTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func EditTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func DeleteTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

	res := utils.HttpResponseBody{
		Status:  http.StatusOK,
		Message: "Successful delete the trial session",
		Data:    nil,
	}

	TrialId := r.URL.Query().Get("trial_id")

	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to connect to the database"
		res.Data = TrialId
		res.UpdateHttpResponse(w)
		return
	}

	col := db.Client.Database("Neukod").Collection("Trial")

	filter := bson.M{"trial_id": TrialId}
	if delRes, err := col.DeleteOne(db.Context, filter); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to delete the data"
		res.Data = TrialId
		res.UpdateHttpResponse(w)
		return
	} else if delRes.DeletedCount < 1 {
		res.Status = http.StatusNotFound
		res.Message = "No document found to be deleted"
		res.Data = TrialId
		res.UpdateHttpResponse(w)
		return
	}
	res.Data = TrialId
	res.UpdateHttpResponse(w)
}
