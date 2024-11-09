package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rluders/httpsuite"
	"log"
	"net/http"
)

type SampleRequest struct {
	Name string `json:"name" validate:"required,min=3"`
	Age  int    `json:"age" validate:"required,min=1"`
}

func (r *SampleRequest) SetParam(fieldName, value string) error {
	switch fieldName {
	case "name":
		r.Name = value
	}
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/submit/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Parse the request and validate it
		req, err := httpsuite.ParseRequest[*SampleRequest](w, r, "name")
		if err != nil {
			log.Printf("Error parsing or validating request: %v", err)
			return
		}

		// Step 2: Send a success response
		httpsuite.SendResponse(w, "Request received successfully", http.StatusOK, &req)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
