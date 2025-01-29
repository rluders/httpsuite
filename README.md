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

### Supported routers

- Gorilla MUX
- Chi
- Go Standard
- ...maybe more? Submit a PR with an example.

## Installation

To install **httpsuite**, run:

```
go get github.com/rluders/httpsuite/v2
```

## Usage

### Request Parsing with URL Parameters

Check out the [example folder for a complete project](./examples) demonstrating how to integrate **httpsuite** into 
your Go microservices.

## Contributing

Contributions are welcome! Feel free to open issues, submit pull requests, and help improve **httpsuite**.

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.