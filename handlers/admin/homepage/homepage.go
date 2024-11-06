package admin

import (
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/Neukod-Academy/neukod-backend/utils"
)

func CreateAboutUs(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusInternalServerError,
		Message: "Failed to process this request",
		Data:    nil,
	}
	db := new(utils.Mongo)
	if stored, err := db.FindAllData("Neukod", "Home"); err != nil {
		res.Message = "Failed to find the data"
		res.UpdateHttpResponse(w)
		return err
	} else if len(stored.Data) > 0 {
		res.Message = "Unable to add a new data since the data is already exist"
		res.UpdateHttpResponse(w)
		return err
	}

	var newAboutUs models.AboutUs
	newAboutUs, err := utils.HttpReqReader[models.AboutUs](r)
	if err != nil {
		res.Message = "ERROR: failed while reading the request data"
		res.UpdateHttpResponse(w)
		return err
	}
	newAboutUs.CreatedAt = time.Now()
	newAboutUs.UpdatedAt = time.Now()

	if _, err := db.InsertNewData("Neukod", "Home", newAboutUs); err != nil {
		res.Message = "ERROR: FAILED while adding a new data"
		res.UpdateHttpResponse(w)
		return err
	}
	res.UpdateHttpResponse(w)
	return nil
}

// func UpdateAboutUse(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
// 	}
// 	res := utils.HttpResponseBody{
// 		Status:  http.StatusInternalServerError,
// 		Message: "Failed to process this request",
// 		Data:    nil,
// 	}
// 	return nil

// }
