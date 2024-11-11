package main

import (
	"log"
	"net/http"

	"github.com/Neukod-Academy/neukod-backend/handlers/session"
	"github.com/Neukod-Academy/neukod-backend/handlers/user"
	"github.com/Neukod-Academy/neukod-backend/middleware"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
	"github.com/Neukod-Academy/neukod-backend/utils"
)

func main() {

	http.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
		res := utils.HttpResponseBody{
			Status:  http.StatusOK,
			Message: "welcome to the homepage",
			Data:    nil,
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		res.UpdateHttpResponse(w)
	})
	http.HandleFunc("/v1/login", session.CreateSession)
	http.HandleFunc("/v1/trialclass", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			user.NewTrial(w, r)
		case http.MethodDelete:
			user.DeleteTrial(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			session.CreateAccount(w, r)
		case http.MethodGet:
			middleware.AuthMiddleware(session.ShowAccount)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("The server is running at http://localhost:%s", env.LOCAL_PORT)
	if err := http.ListenAndServe(":"+env.LOCAL_PORT, nil); err != nil {
		panic("unable to listen at this ports")
	}
}
