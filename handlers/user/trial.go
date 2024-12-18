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
	if r.Method != http.MethodPost {
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
	}
	if err := newTrial.CheckIfEmpty(); len(err) != 0 {
		res.Status = http.StatusBadRequest
		res.Message = err
		res.UpdateHttpResponse(w)
	}

	isTheCourseAllowed := false
	for _, course := range models.TrialList {
		if newTrial.Course == course {
			isTheCourseAllowed = true
			continue
		}
	}

	if !isTheCourseAllowed {
		res.Status = http.StatusBadRequest
		res.Message = "This course is still unavailable"
		res.UpdateHttpResponse(w)
	}

	db := new(utils.Mongo)
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Failed to create a client to the database"
		res.UpdateHttpResponse(w)
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
	}

	res.Data = newTrial
	res.UpdateHttpResponse(w)
}

func ShowTrial(w http.ResponseWriter, r *http.Request) {
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to show the trial list",
		Data:    nil,
	}
	if r.Method != http.MethodGet {
		res.Status = http.StatusMethodNotAllowed
		res.Message = "Method is not allowed"
		res.Data = nil
		res.UpdateHttpResponse(w)
		return
	}

	db := utils.Mongo{}
	if err := db.CreateClient(env.MONGO_URI); err != nil {
		res.Message = "Failed to create a client to the database"
		res.Data = nil
		res.UpdateHttpResponse(w)
		return
	}
	defer db.CloseClientDB()

	col := db.Client.Database("Neukod").Collection("Trial")
	cursor, err := col.Find(db.Context, bson.M{}, options.Find())
	if err != nil {
		res.Message = "Failed to find the data"
		res.Data = nil
		res.UpdateHttpResponse(w)
		return
	}
	stored := models.Trial{}
	if err := cursor.All(db.Context, &stored); err != nil {
		res.Message = "Failed to find the data"
		res.Data = nil
		res.UpdateHttpResponse(w)
		return
	}
	res.Status = http.StatusOK
	res.Message = "Successful to get all of the trial data"
	res.Data = stored
	res.UpdateHttpResponse(w)
}

func ConfirmTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func EditTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func DeleteTrial(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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
	defer db.CloseClientDB()

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
		res.Message = "The trial class is not found to be deleted"
		res.Data = TrialId
		res.UpdateHttpResponse(w)
		return
	}
	res.Data = TrialId
	res.UpdateHttpResponse(w)
}
