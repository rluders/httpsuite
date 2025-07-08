# httpsuite

httpsuite is a lightweight, idiomatic Go library that simplifies HTTP request parsing, validation, 
and response handling in microservices. It’s designed to reduce boilerplate and promote clean, 
maintainable, and testable code — all while staying framework-agnostic.

## ✨ Features

- 🧾 **Request Parsing**: Automatically extract and map JSON payloads and URL path parameters to Go structs.
- ✅ **Validation:** Centralized validation using struct tags, integrated with standard libraries like `go-playground/validator`.
- 📦 **Unified Responses:** Standardize your success and error responses (e.g., [RFC 7807 Problem Details](https://datatracker.ietf.org/doc/html/rfc7807)) for a consistent API experience.
- 🔌 **Modular Design:** Use each component independently — ideal for custom setups, unit testing, or advanced use cases.
- 🧪 **Test-Friendly:** Decouple parsing and validation logic for simpler, more focused test cases.

### 🔌 Supported routers

- [Chi](https://github.com/go-chi/chi)
- [Gorilla MUX](https://github.com/gorilla/mux)
- Go standard `http.ServeMux`
- ...and potentially more — [Submit a PR with an example!](https://github.com/rluders/httpsuite)

## 🛠 Installation

To install **httpsuite**, run:

```
go get github.com/rluders/httpsuite/v2
```

## 🚀 Usage

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/rluders/httpsuite/v2"
    "net/http"
)

type SampleRequest struct {
    ID   int    `json:"id" validate:"required"`
    Name string `json:"name" validate:"required,min=3"`
}

func (r *SampleRequest) SetParam(fieldName, value string) error {
    if fieldName == "id" {
        id, err := strconv.Atoi(value)
        if err != nil {
            return err
        }
        r.ID = id
    }
    return nil
}

func main() {
    r := chi.NewRouter()

    r.Post("/submit/{id}", func(w http.ResponseWriter, r *http.Request) {
        req, err := httpsuite.ParseRequest[*SampleRequest](w, r, chi.URLParam, "id")
        if err != nil {
            return // ProblemDetails already sent
        }

        httpsuite.SendResponse(w, http.StatusOK, req, nil, nil)
    })

    http.ListenAndServe(":8080", r)
}
```

💡 Try it:

```
curl -X POST http://localhost:8080/submit/123 \
    -H "Content-Type: application/json" \
    -d '{"name":"John"}'
```

## 📂 Examples

Check out the `examples/` folder for a complete working project demonstrating:

- Full request lifecycle
- Param parsing
- Validation
- ProblemDetails usage
- JSON response formatting

## 📖 Tutorial & Article

- [Improving Request Validation and Response Handling in Go Microservices](https://medium.com/@rluders/improving-request-validation-and-response-handling-in-go-microservices-cc54208123f2)

## 🤝 Contributing

All contributions are welcome! Whether it's a bug fix, feature proposal, or router integration example:

- Open an issue
- Submit a PR
- Join the discussion!

## 🪪 License

The MIT License (MIT). Please see [License File](LICENSE) for more information.