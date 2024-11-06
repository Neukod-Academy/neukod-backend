package main

import (
	"log"
	"net/http"

	"github.com/Neukod-Academy/neukod-backend/handlers/admin"
	"github.com/Neukod-Academy/neukod-backend/handlers/session"
	"github.com/Neukod-Academy/neukod-backend/handlers/user"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the home page"))
	})
	http.HandleFunc("/admin/contents", admin.CreateContent)
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
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			session.CreateAccount(w, r)
		case http.MethodGet:
			session.ShowAccount(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Printf("The server is running at http://localhost:%s", env.LOCAL_PORT)
	if err := http.ListenAndServe(":"+env.LOCAL_PORT, nil); err != nil {
		panic("unable to listen at this ports")
	}
}
