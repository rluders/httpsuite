module gorillamux_example

go 1.23

require (
	github.com/gorilla/mux v1.8.1
	github.com/rluders/httpsuite/v3 v3.0.0
)

replace github.com/rluders/httpsuite/v3 => ../..
