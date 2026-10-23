# httpsuite

[![Go Reference](https://pkg.go.dev/badge/github.com/rluders/httpsuite/v3.svg)](https://pkg.go.dev/github.com/rluders/httpsuite/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/rluders/httpsuite/v3)](https://goreportcard.com/report/github.com/rluders/httpsuite/v3)
[![GitHub release](https://img.shields.io/github/v/release/rluders/httpsuite?sort=semver)](https://github.com/rluders/httpsuite/releases)
[![codecov](https://codecov.io/gh/rluders/httpsuite/graph/badge.svg?token=E9mGGRHicy)](https://codecov.io/gh/rluders/httpsuite)
[![License](https://img.shields.io/github/license/rluders/httpsuite)](LICENSE)
![stdlib-only](https://img.shields.io/badge/dependencies-stdlib--only-blue)
![RFC 9457](https://img.shields.io/badge/errors-RFC%209457-orange)

`httpsuite` is a Go library for request parsing, response writing, and RFC 9457 problem responses.

`v3` keeps the root module stdlib-only and moves validation to an optional submodule.

## Features

* Parse JSON request bodies with a default `1 MiB` limit
* Return `413 Payload Too Large` when the configured body limit is exceeded
* Bind path params explicitly through a router-specific extractor
* Validate automatically during `ParseRequest` when a global validator is configured
* Keep `ParseRequest` panic-safe for invalid inputs and return regular errors instead
* Return consistent RFC 9457 Problem Details
* Write success responses with optional generic metadata
* Support both direct helpers and optional builders

## Why httpsuite?

Writing HTTP handlers in Go often leads to repetitive and error-prone code:

* Manual JSON decoding
* Manual validation wiring
* Ad-hoc path parameter extraction
* Inconsistent error handling
* Repeated boilerplate across handlers

`httpsuite` provides a consistent and safe flow for parsing requests and writing responses, while staying lightweight and idiomatic.

## What you gain

### Less boilerplate

A single entry point replaces multiple manual steps:

```go
req, err := httpsuite.ParseRequest[*MyRequest](w, r, chi.URLParam, nil, "id")
if err != nil {
    return
}
```

* JSON decoding
* Path param binding
* Optional validation

All handled in one place, reducing handler complexity.

### Safe by default

* `ParseRequest` never panics on invalid inputs
* Handles nil request and nil body safely
* Enforces body size limits with proper HTTP responses

This reduces the risk of runtime crashes and undefined behavior.

### Consistent error responses

Errors follow **RFC 9457 Problem Details** out of the box.

Instead of ad-hoc responses like:

```json
{ "error": "invalid request" }
```

You get structured, standardized responses:

```json
{
  "type": "...",
  "title": "...",
  "status": 400,
  "detail": "..."
}
```

This improves API consistency and client integration.

### Centralized validation

* Configure once with `SetValidator(...)`
* Automatically applied during `ParseRequest`
* Can be overridden per request

This eliminates duplicated validation logic across handlers.

### Lightweight and modular

* Core module is **stdlib-only**
* No framework lock-in
* Validation is optional

You keep full control over your stack.

### Flexible response API

Choose the level of abstraction you need:

**Simple**

```go
httpsuite.OK(w, data)
```

**Fluent**

```go
httpsuite.Reply().Meta(meta).OK(w, data)
```

**Advanced**

```go
httpsuite.Problem(...).Title(...).Build()
```

### Router-agnostic

Works with:

* Chi
* Gorilla Mux
* net/http

No need to change your router or architecture.

## Comparison

### ❌ Without httpsuite

```go
var req MyRequest

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "invalid body", 400)
    return
}

idStr := chi.URLParam(r, "id")
id, err := strconv.Atoi(idStr)
if err != nil {
    http.Error(w, "invalid id", 400)
    return
}

// validation...
// more error handling...
```

**Problems:**

* Repeated logic across handlers
* Inconsistent error formats
* Easy to miss edge cases

### ✅ With httpsuite

```go
req, err := httpsuite.ParseRequest[*MyRequest](w, r, chi.URLParam, nil, "id")
if err != nil {
    return
}

httpsuite.OK(w, req)
```

**Benefits:**

* Single, consistent flow
* Less code
* Safer defaults

## Supported routers

* Chi
* Gorilla Mux
* Go standard `http.ServeMux`

## Installation

### Core

```bash
go get github.com/rluders/httpsuite/v3
```

### Optional validation adapter

```bash
go get github.com/rluders/httpsuite/validation/playground
```

## Mental model

* request in: `ParseRequest(...)`
* success out: `OK(...)`, `Created(...)`, `Reply().Meta(...).OK(...)`
* problem out: `ProblemResponse(...)`, `NewBadRequestProblem(...)`, `Problem(...).Build()`
* validation: configure once with `SetValidator(...)`, override locally with `ParseOptions.Validator`

For simple handlers, prefer direct helpers.

When a handler needs custom headers, metadata, or problem composition, use the optional builders.

`ParseRequest` never panics on invalid inputs such as a nil request, nil body, or nil path extractor. These cases return regular Go errors so callers can fail safely.

## Quick start

### Core only

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rluders/httpsuite/v3"
)

type CreateUserRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (r *CreateUserRequest) SetParam(fieldName, value string) error {
	if fieldName != "id" {
		return nil
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	r.ID = id
	return nil
}

func main() {
	router := chi.NewRouter()

	router.Post("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		req, err := httpsuite.ParseRequest[*CreateUserRequest](w, r, chi.URLParam, nil, "id")
		if err != nil {
			return
		}

		httpsuite.OK(w, req)
	})

	_ = http.ListenAndServe(":8080", router)
}
```

Try it:

```bash
curl -X POST http://localhost:8080/users/123 \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada"}'
```

### Core + validation

```go
validator := playground.NewWithValidator(nil, &httpsuite.ProblemConfig{
	BaseURL: "https://api.example.com",
})

httpsuite.SetValidator(validator)

req, err := httpsuite.ParseRequest[*CreateUserRequest](
	w,
	r,
	chi.URLParam,
	&httpsuite.ParseOptions{
		MaxBodyBytes: 1 << 20,
	},
	"id",
)
```

### Problem type URLs

Set the base URL and extra error type paths once at startup; `GetProblemTypeURL` and the problem helpers pick them up. `ParseOptions.Problems` still overrides per call.

```go
cfg := httpsuite.NewProblemConfig()
cfg.BaseURL = "https://api.example.com"
cfg.ErrorTypePaths["payment_required_error"] = "/errors/payment-required"
httpsuite.SetProblemConfig(cfg)

// restore the built-in default
httpsuite.SetProblemConfig(httpsuite.NewProblemConfig())
```

## Direct helpers

```go
httpsuite.OK(w, user)
httpsuite.OKWithMeta(w, users, httpsuite.NewPageMeta(page, pageSize, totalItems))
httpsuite.Created(w, user, "/users/42")
httpsuite.Accepted(w, job)
httpsuite.NoContent(w)
httpsuite.ProblemResponse(w, httpsuite.NewNotFoundProblem("user not found"))
httpsuite.ProblemResponse(w, httpsuite.NewUnauthorizedProblem("missing token"))
httpsuite.ProblemResponse(w, httpsuite.NewForbiddenProblem("admin only"))
httpsuite.ProblemResponse(w, httpsuite.NewConflictProblem("email already taken"))
httpsuite.ProblemResponse(w, httpsuite.NewUnprocessableProblem("order already shipped"))
httpsuite.ProblemResponse(w, httpsuite.NewInternalProblem("unexpected error"))
```

## Fluent helpers

```go
httpsuite.Reply().
	Meta(httpsuite.NewPageMeta(page, pageSize, totalItems)).
	OK(w, users)

httpsuite.Reply().
	Header("X-Request-ID", requestID).
	Created(w, user, "/users/42")

httpsuite.Reply().
	Header("Location", "/jobs/7").
	Accepted(w, job)

httpsuite.Reply().
	Header("X-Request-ID", requestID).
	NoContent(w)
```

## Builders

```go
problem := httpsuite.Problem(http.StatusNotFound).
	Type(httpsuite.GetProblemTypeURL("not_found_error")).
	Title("User Not Found").
	Detail("user 42 does not exist").
	Instance("/users/42").
	Build()

httpsuite.RespondProblem(problem).
	Header("X-Trace-ID", traceID).
	Write(w)
```

## Precognition

Let clients validate a request against the endpoint's real decoding, binding and validation rules without running business logic. Opt in per endpoint:

```go
req, err := httpsuite.ParseRequest[*CreateUserRequest](w, r, chi.URLParam, &httpsuite.ParseOptions{
	Precognition: true,
}, "id")
if err != nil {
	// Covers parse errors, validation errors and handled precognitive requests.
	return
}

// Never reached for a precognitive request.
user := store.Create(req)
```

A client sends the same request with `Precognition: true`:

| Outcome | Response | `err` |
| --- | --- | --- |
| Valid | `204 No Content`, `Precognition-Success: true` | `ErrPrecognitionHandled` |
| Validation failed | `422` with the validator's `ProblemDetails` | validation error |
| Malformed body / too large / bad path param | usual `400` / `413` problem | parse error |

Notes:

* Without `ParseOptions.Precognition`, the header is ignored and the request is processed normally.
* Opted-in endpoints add `Vary: Precognition` to responses.
* Use `httpsuite.IsPrecognitionHandled(err)` to tell a handled precognitive request apart from a failure, and `httpsuite.IsPrecognitionRequest(r)` to inspect the request.
* With `SkipValidation`, a precognitive request only checks decoding and path binding.
* A successful pre-validation does not guarantee the real request will succeed; state may change in between.

## Architecture

* root module: `github.com/rluders/httpsuite/v3`
* optional validation adapter: `github.com/rluders/httpsuite/validation/playground`
* root stays stdlib-only
* validation is opt-in at bootstrap, automatic at parse time when configured
* response metadata is generic and can use `PageMeta` or `CursorMeta`

## Migration from v2 to v3

* update imports from `v2` to `v3`
* update `ParseRequest` calls to pass `opts` before `pathParams`
* configure validation globally with `httpsuite.SetValidator(...)`
* `ParseRequest` now validates automatically when a validator is configured
* validator-provided status is respected
* use `ParseOptions.SkipValidation` to opt out per call
* use `ParseOptions.Validator` to override per call
* use `ProblemConfig` for custom problem type URLs

## Examples

Examples live in `examples/`.

* `examples/stdmux`: core-only with `http.ServeMux`
* `examples/gorillamux`: path params with Gorilla Mux
* `examples/chi`: global validation with Chi
* `examples/restapi`: full REST example with pagination and custom problems

## Notes for contributors

* request helpers: `request*.go`
* response helpers/builders: `response*.go`
* problem handling: `problem*.go`

## Release workflow

Supports:

* pushing an existing `v*` tag
* manual release with semantic version bump

## Tutorial

* [Improving Request Validation and Response Handling in Go Microservices](https://medium.com/@rluders/improving-request-validation-and-response-handling-in-go-microservices-cc54208123f2), about the first released version.

## Contributing

* open an issue
* submit a PR
* add a router example

## License

MIT
