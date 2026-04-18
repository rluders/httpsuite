module chi_example

go 1.25.0

require (
	github.com/go-chi/chi/v5 v5.2.0
	github.com/rluders/httpsuite/v3 v3.0.0
	github.com/rluders/httpsuite/validation/playground v0.0.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)

replace github.com/rluders/httpsuite/v3 => ../..

replace github.com/rluders/httpsuite/validation/playground => ../../validation/playground
