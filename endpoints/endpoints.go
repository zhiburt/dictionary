package endpoints

import (
	"context"
	"fmt"

	"github.com/dictionary/services"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints provide all Go kit endpoints for this one
type Endpoints struct {
	Words      endpoint.Endpoint
	AddNewWord endpoint.Endpoint
	GetByID    endpoint.Endpoint
	GetByW     endpoint.Endpoint
}

// NewEndpoints returns new object of Endpoints for service s
func NewEndpoints(s services.Dictionary) Endpoints {
	return Endpoints{
		Words:      makeWordsEndpoint(s),
		AddNewWord: makeAddNewWordEndpoint(s),
		GetByID:    makeGetByIDEndpoint(s),
		GetByW:     makeGetByWEndpoint(s),
	}
}

func makeWordsEndpoint(s services.Dictionary) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, ok := request.(WordsRequest)
		if ok {
			resp, err := s.Words(ctx)
			return WordsResponse{resp, err}, nil
		}
		return WordsResponse{nil, fmt.Errorf("doesn't have this request")}, nil
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

func makeGetByWEndpoint(s services.Dictionary) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetByWRequest)
		resp, err := s.GetByW(ctx, req.W)
		return GetByWResponse{resp, err}, nil
	}
}
