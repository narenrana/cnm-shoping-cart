package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"shopping-cart/cnm-auth/models"
	"shopping-cart/cnm-auth/services"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the booking service.
func MakeHandler(bs services.Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	authLoginHandler := kithttp.NewServer(
		makeAuthLoginEndpoint(bs),
		decodeAuthLoginRequest,
		encodeResponse,
		opts...,
	)

	authLogoutHandler := kithttp.NewServer(
		makeAuthLogoutEndpoint(bs),
		decodeAuthLogoutRequest,
		encodeResponse,
		opts...,
	)

	authSignUpHandler := kithttp.NewServer(
		makeSignUpEndpoint(bs),
		decodeAuthSignUpRequest,
		encodeResponse,
		opts...,
	)

	authRecoverPasswordHandler := kithttp.NewServer(
		makeRecoverPasswordEndpoint(bs),
		decodeAuthRecoverPasswordRequest,
		encodeResponse,
		opts...,
	)

	authRefreshTokenHandler := kithttp.NewServer(
		makeRefreshTokenEndpoint(bs),
		decodeAuthRefreshTokenRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/auth/v1/login", authLoginHandler).Methods("POST")
	r.Handle("/auth/v1/logout", authLogoutHandler).Methods("GET")
	r.Handle("/auth/v1/signup", authSignUpHandler).Methods("POST")
	r.Handle("/auth/v1/recoverPassword", authRecoverPasswordHandler).Methods("POST")
	r.Handle("/auth/v1/refreshToken", authRefreshTokenHandler).Methods("POST")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeAuthLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AuthLoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAuthLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {

	value := r.FormValue("token")
	return models.AuthLogoutRequest{
		Token:       value,
	}, nil
}

func decodeAuthSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AuthSignUpRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAuthRecoverPasswordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AuthRecoverPasswordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}



func decodeAuthRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AuthRecoverPasswordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		return nil, err
	}
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
	case services.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
