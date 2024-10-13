package main

import (
	"log"
	"net/http"

	"github.com/Neukod-Academy/neukod-backend/handlers/admin"
	"github.com/Neukod-Academy/neukod-backend/pkg/env"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the home page"))
	})
	http.HandleFunc("/admin/contents", admin.CreateContent)

	log.Printf("The server is running at http://localhost:%s", env.LOCAL_PORT)
	if err := http.ListenAndServe(":"+env.LOCAL_PORT, nil); err != nil {
		panic("unable to listen at this ports")
	}
}
