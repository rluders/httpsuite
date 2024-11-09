# httpsuite

**httpsuite** is a Go library designed to simplify the handling of HTTP requests, validations, and responses 
in microservices. By providing a clear structure and modular approach, it helps developers write 
cleaner, more maintainable code with reduced boilerplate.

## Features

- **Request Parsing**: Streamline the parsing of incoming HTTP requests, including URL parameters.
- **Validation:** Centralize validation logic for easy reuse and consistency.
- **Response Handling:** Standardize responses across your microservices for a unified client experience.
- **Modular Design:** Each component (Request, Validation, Response) can be used independently, 
enhancing testability and flexibility.

> **Note:** Currently it only supports Chi.

## Installation

To install **httpsuite**, run:

```
go get github.com/rluders/httpsuite
```

## Usage

### Request Parsing with URL Parameters

Easily parse incoming requests and set URL parameters:

```go
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
```

Check out the [example folder for a complete project](./examples) demonstrating how to integrate **httpsuite** into 
your Go microservices.

## Contributing

Contributions are welcome! Feel free to open issues, submit pull requests, and help improve **httpsuite**.

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.