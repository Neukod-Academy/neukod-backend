package main

import (
	"log"
	"net/http"

	"github.com/Neukod-Academy/neukod-backend/handlers/index"
	"github.com/Neukod-Academy/neukod-backend/handlers/session"
	"github.com/Neukod-Academy/neukod-backend/handlers/user"
	"github.com/Neukod-Academy/neukod-backend/middleware"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
)

func main() {
	http.HandleFunc("/v1", index.RetreiveHomepage)
	http.HandleFunc("/v1/auth/signin", session.CreateSession)
	http.HandleFunc("/v1/auth/signout", middleware.AuthMiddleware(session.DropSession))
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
			middleware.AuthMiddleware(session.ShowAccounts)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/v1/users/{id}", middleware.AuthMiddleware(session.RemoveAccount))

	log.Printf("The server is running at http://localhost:%s", env.LOCAL_PORT)
	if err := http.ListenAndServe(":"+env.LOCAL_PORT, nil); err != nil {
		panic("unable to listen at this ports")
	}
}
