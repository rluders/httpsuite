package main

import (
	"github.com/rluders/httpsuite/v3"
	"log"
	"net/http"
	"strconv"
)

type SampleRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
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

func StdMuxParamExtractor(r *http.Request, key string) string {
	// Remove "/submit/" (7 characters) from the URL path to get just the "id"
	// Example: /submit/123 -> 123
	return r.URL.Path[len("/submit/"):] // Skip the "/submit/" part
}

// You can test it using:
//
//	curl -X POST http://localhost:8080/submit/123 \
//		-H "Content-Type: application/json" \
//		-d '{"name": "John Doe", "age": 30}'
func main() {
	// Creating the router using the Go standard mux
	mux := http.NewServeMux()

	// Define the endpoint POST
	mux.HandleFunc("/submit/", func(w http.ResponseWriter, r *http.Request) {
		// Using the function for parameter extraction to the ParseRequest
		req, err := httpsuite.ParseRequest[*SampleRequest](w, r, StdMuxParamExtractor, nil, "id")
		if err != nil {
			log.Printf("Error parsing request: %v", err)
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
	log.Fatal(http.ListenAndServe(":8080", mux))
}
