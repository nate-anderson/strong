package strong

import "net/http"

// handle different error codes differently
var errorHandlers map[int]ErrorFunc

func init() {
	errorHandlers = make(map[int]ErrorFunc)
}

type ErrorFunc func(w http.ResponseWriter, req *http.Request, err error)

var defaultErrorFunc ErrorFunc = func(w http.ResponseWriter, req *http.Request, err error) {
	http.Error(
		w,
		err.Error(),
		http.StatusBadRequest,
	)
}

func HandleErrorCode(statusCode int, handler ErrorFunc) {
	errorHandlers[statusCode] = handler
}

// use the default error handler unless an override is set by the caller
func handleRequestError(statusCode int, err error, w http.ResponseWriter, req *http.Request) {
	if handler, ok := errorHandlers[statusCode]; ok {
		handler(w, req, err)
	}
	defaultErrorFunc(w, req, err)
}
