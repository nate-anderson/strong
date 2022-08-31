package strong

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Request[T any] struct {
	raw  *http.Request
	body *T
}

type Response[T any] struct {
	w    http.ResponseWriter
	body *T
}

type ResponseError interface {
	error
	Code() int
}

// JSONRoute takes a strongly-typed handler function and returns a standard http.HandlerFunc
// The request body is JSON decoded as ReqBodyT, and the response is JSON encoded and returned to the client
func JSONRoute[ReqBodyT any, ResT any](handler func(req *Request[ReqBodyT]) (*Response[ResT], ResponseError)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bodyDecoder := json.NewDecoder(req.Body)
		reqBody := new(ReqBodyT)
		err := bodyDecoder.Decode(&reqBody)
		if err != nil {
			handleRequestError(http.StatusBadRequest, err, w)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w)
			return
		}

		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(res.body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w)
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
			handleRequestError(http.StatusBadRequest, err, w)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w)
			return
		}

		bodyEncoder := xml.NewEncoder(w)
		err = bodyEncoder.Encode(res.body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w)
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
			handleRequestError(http.StatusBadRequest, err, w)
			return
		}

		reqObj := Request[ReqBodyT]{
			raw:  req,
			body: reqBody,
		}

		res, resErr := handler(&reqObj)
		if err != nil {
			handleRequestError(resErr.Code(), err, w)
			return
		}

		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(res.body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			handleRequestError(http.StatusInternalServerError, err, w)
			return
		}
	}
}
