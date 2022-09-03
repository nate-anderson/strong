package strong

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Request[T any] struct {
	raw  *http.Request
	body *T
}

func (r Request[T]) Ctx() context.Context {
	return r.raw.Context()
}

func (r Request[T]) Body() *T {
	return r.body
}

func (r Request[T]) Header() *http.Header {
	return &r.raw.Header
}

func (r Request[T]) Cookies() *http.CookieJar {
	return r.Cookies()
}

type Response[T any] struct {
	w    http.ResponseWriter
	Body *T
}

type ResponseError interface {
	error
	Code() int
}

type responseError struct {
	err  error
	code int
}

func (e responseError) Error() string {
	return e.err.Error()
}

func (e responseError) Code() int {
	return e.code
}

func Error(code int, err error) responseError {
	return responseError{
		err:  err,
		code: code,
	}
}

// JSONRoute takes a strongly-typed handler function and returns a standard http.HandlerFunc
// The request body is JSON decoded as ReqBodyT, and the response is JSON encoded and returned to the client
func JSONRoute[ReqBodyT any, ResT any](handler func(req *Request[ReqBodyT]) (*Response[ResT], ResponseError)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bodyDecoder := json.NewDecoder(req.Body)
		reqBody := new(ReqBodyT)
		err := bodyDecoder.Decode(&reqBody)
		if err != nil {
			handleRequestError(http.StatusBadRequest, err, w, req)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w, req)
			return
		}

		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(res.Body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w, req)
			return
		}
	}
}

// XMLRoute takes a strongly-typed handler function and returns a standard http.HandlerFunc
// The request body is XML decoded as ReqBodyT, and the response is XML encoded and returned to the client
func XMLRoute[ReqBodyT any, ResT any](handler func(req *Request[ReqBodyT]) (*Response[ResT], ResponseError)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bodyDecoder := xml.NewDecoder(req.Body)
		reqBody := new(ReqBodyT)
		err := bodyDecoder.Decode(&reqBody)
		if err != nil {
			handleRequestError(http.StatusBadRequest, err, w, req)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w, req)
			return
		}

		bodyEncoder := xml.NewEncoder(w)
		err = bodyEncoder.Encode(res.Body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w, req)
			return
		}
	}
}

// JSONRoute takes a strongly-typed handler function and returns a standard http.HandlerFunc
// The request body is decoded as  ReqBodyT, and the response is JSON encoded and returned to the client
func FormRoute[ReqBodyT any, ResT any](handler func(req *Request[ReqBodyT]) (*Response[ResT], ResponseError)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		reqBody := new(ReqBodyT)
		err = unmarshalForm(req.Form, &reqBody)
		if err != nil {
			handleRequestError(http.StatusBadRequest, err, w, req)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w, req)
			return
		}

		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(res.Body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w, req)
			return
		}
	}
}
