package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Neukod-Academy/neukod-backend/models"
	"github.com/google/uuid"
)

func CreateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := new(models.Content)
	var ByteBody []byte
	var err error

	if ByteBody, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "Failed while decoding into json format", http.StatusInternalServerError)
		return
	} else if err = json.Unmarshal(ByteBody, &data); err != nil {
		http.Error(w, "Failed while decoding into json format", http.StatusInternalServerError)
		return
	} else {
		data.ID = uuid.New().String()
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if ByteBody, err = json.Marshal(data); err != nil {
			http.Error(w, "Failed while decoding into json format", http.StatusInternalServerError)
			return
		} else {
			w.Write(ByteBody)
		}

	}

}

func DeleteContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}

func EditContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT " {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}

}

func GetContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}
