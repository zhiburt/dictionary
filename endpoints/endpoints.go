package endpoints

import (
	"context"

	"github.com/dictionary/services"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints provide all Go kit endpoints for this one
type Endpoints struct {
	AddNewWord endpoint.Endpoint
	GetByID    endpoint.Endpoint
}

// NewEndpoints returns new object of Endpoints for service s
func NewEndpoints(s services.Dictionary) Endpoints {
	return Endpoints{
		AddNewWord: makeAddNewWordEndpoint(s),
		GetByID:    makeGetByIDEndpoint(s),
	}
}

func makeAddNewWordEndpoint(s services.Dictionary) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddNewWordRequest)
		resp, err := s.AddNewWord(ctx, req.Word)
		return AddNewWordResponse{resp, err}, nil
	}
}

func makeGetByIDEndpoint(s services.Dictionary) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByIDRequest)
		resp, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{resp, err}, nil
	}
}
