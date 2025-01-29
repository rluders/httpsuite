package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rluders/httpsuite/v2"
	"log"
	"net/http"
	"strconv"
)

type SampleRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3"`
	Age  int    `json:"age" validate:"required,min=1"`
}

type SampleResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (r *SampleRequest) SetParam(fieldName, value string) error {
	switch fieldName {
	case "id":
		id, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		r.ID = id
	}
	return nil
}

func ChiParamExtractor(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// You can test it using:
//
//	curl -X POST http://localhost:8080/submit/123 \
//		-H "Content-Type: application/json" \
//		-d '{"name": "John Doe", "age": 30}'
//
// And you should get:
//
// {"data":{"id":123,"name":"John Doe","age":30}}
func main() {
	// Creating the router with Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define the endpoint POST
	r.Post("/submit/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Using the function for parameter extraction to the ParseRequest
		req, err := httpsuite.ParseRequest[*SampleRequest](w, r, ChiParamExtractor, "id")
		if err != nil {
			log.Printf("Error parsing or validating request: %v", err)
			return
		}

		resp := &SampleResponse{
			ID:   req.ID,
			Name: req.Name,
			Age:  req.Age,
		}

		// Sending success response
		httpsuite.SendResponse[SampleResponse](w, http.StatusOK, *resp, nil, nil)
	})

	// Starting the server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
