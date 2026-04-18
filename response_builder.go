package httpsuite

import "net/http"

// ReplyBuilder configures metadata and headers before writing a response.
type ReplyBuilder struct {
	meta    any
	headers http.Header
}

// ResponseBuilder builds and writes HTTP responses declaratively.
type ResponseBuilder[T any] struct {
	code    int
	data    T
	meta    any
	problem *ProblemDetails
	headers http.Header
}

// Respond starts a success response builder.
func Respond[T any](data T) *ResponseBuilder[T] {
	return &ResponseBuilder[T]{
		code: http.StatusOK,
		data: data,
	}
}

// RespondProblem starts a problem response builder.
func RespondProblem(problem *ProblemDetails) *ResponseBuilder[any] {
	code := http.StatusInternalServerError
	if problem != nil && problem.Status >= 400 && problem.Status <= 599 {
		code = problem.Status
	}
	return &ResponseBuilder[any]{
		code:    code,
		problem: problem,
	}
}

// Meta sets response metadata for a fluent helper chain.
func (b *ReplyBuilder) Meta(meta any) *ReplyBuilder {
	b.meta = meta
	return b
}

// Header sets a single response header for a fluent helper chain.
func (b *ReplyBuilder) Header(key, value string) *ReplyBuilder {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	b.headers.Set(key, value)
	return b
}

// Headers merges multiple response headers for a fluent helper chain.
func (b *ReplyBuilder) Headers(headers http.Header) *ReplyBuilder {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	for key, values := range headers {
		for _, value := range values {
			b.headers.Add(key, value)
		}
	}
	return b
}

// OK writes a 200 JSON response using the fluent helper configuration.
func (b *ReplyBuilder) OK(w http.ResponseWriter, data any) {
	Respond(data).Meta(b.meta).Headers(b.headers).Write(w)
}

// Created writes a 201 JSON response using the fluent helper configuration.
func (b *ReplyBuilder) Created(w http.ResponseWriter, data any, location string) {
	builder := Respond(data).Status(http.StatusCreated).Meta(b.meta).Headers(b.headers)
	if location != "" {
		builder.Header("Location", location)
	}
	builder.Write(w)
}

// Problem writes a problem response using the fluent helper configuration.
func (b *ReplyBuilder) Problem(w http.ResponseWriter, problem *ProblemDetails) {
	RespondProblem(problem).Headers(b.headers).Write(w)
}

// Status overrides the response status code.
func (b *ResponseBuilder[T]) Status(code int) *ResponseBuilder[T] {
	b.code = code
	return b
}

// Meta sets response metadata.
func (b *ResponseBuilder[T]) Meta(meta any) *ResponseBuilder[T] {
	b.meta = meta
	return b
}

// Header sets a single response header.
func (b *ResponseBuilder[T]) Header(key, value string) *ResponseBuilder[T] {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	b.headers.Set(key, value)
	return b
}

// Headers merges multiple response headers.
func (b *ResponseBuilder[T]) Headers(headers http.Header) *ResponseBuilder[T] {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	for key, values := range headers {
		for _, value := range values {
			b.headers.Add(key, value)
		}
	}
	return b
}

// Write writes the configured response.
func (b *ResponseBuilder[T]) Write(w http.ResponseWriter) {
	writeResponse(w, b.code, b.data, b.problem, b.meta, b.headers)
}
