package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(bs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	userAddHandler := kithttp.NewServer(
		makeAddUserEndpoint(bs),
		decodeAddUserRequest,
		encodeResponse,
		opts...,
	)

	authListHandler := kithttp.NewServer(
		makeUserListEndpoint(bs),
		decodeUserListRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/users/v1/add", userAddHandler).Methods("POST")
	r.Handle("/users/v1/list", authListHandler).Methods("GET")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeAddUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request userRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUserListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request userRequest
	//value := r.FormValue("token")
	return request, nil
}


func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
