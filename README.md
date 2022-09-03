# strong

Strongly-typed HTTP handler functions in Go

This is an experimental library for writing RPC-style, strongly-typed HTTP handler functions in Go. It essentially wraps a function of signature `func(request) (response, error)` and abstracts away the JSON, XML or form data marshaling that would normally happen inside the handler.

## Examples:

```go
func (h Handler) CreateUser(req *strong.Request[CreateUserInput]) (*strong.Response[User], error) {
    u, err := h.userRepo.Store(input.Ctx())
    if err != nil {
        return nil, strong.Error(http.StatusInternalServerError, err)
    }

    return &strong.Response[User]{u}, nil
}

h := Handler{}

mux.Post("/user", strong.JSONRoute(h.CreateUser))
```

The `JSONRoute`, `XMLRoute`, and `FormRoute` functions all work similarly. The library is currently limited in that request-type response-type pairs are hardcoded.