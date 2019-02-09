package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dictionary/endpoints"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

var (
	//ErrBadRouting for happend bad routing
	ErrBadRouting = errors.New("bad routing")
)

// NewService return http.Handler for service
func NewService(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	router.Methods("GET").Path("/words/{id}").Handler(
		kithttp.NewServer(
			endpoints.GetByID,
			decodeGetByIDRequest,
			encodeResponse,
			options...),
	)

	router.Methods("POST").Path("/words").Handler(
		kithttp.NewServer(
			endpoints.AddNewWord,
			decodeAddNewWordRequest,
			encodeResponse,
			options...),
	)

	return router
}

func decodeAddNewWordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.AddNewWordRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoints.GetByIDRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encodeError(ctx, err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(b)
}
