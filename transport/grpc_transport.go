package transport

import (
	"context"

	"github.com/dictionary/models/word"

	"github.com/dictionary/endpoints"
	pb "github.com/dictionary/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	addNewWord grpctransport.Handler
	getByID    grpctransport.Handler
}

func (g *grpcServer) AddNewWord(ctx context.Context, req *pb.AddNewWordRequest) (*pb.AddNewWordResponce, error) {
	_, resp, err := g.addNewWord.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddNewWordResponce), nil
}

func (g *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponce, error) {
	_, resp, err := g.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponce), nil
}

// NewGRPCService create and return gprc service
func NewGRPCService(ctx context.Context, endpoint endpoints.Endpoints) pb.DictionaryServer {
	return &grpcServer{
		grpctransport.NewServer(
			endpoint.AddNewWord,
			decodeGRPCAddNewWordRequest,
			encodeGRPCAddNewWordResponce,
		),
		grpctransport.NewServer(
			endpoint.GetByID,
			decodeGRPCGetByIDRequest,
			encodeGRPCGetByIDResponce,
		)}
}

// encodeGRPCAddNewWordRequest endode one
func encodeGRPCAddNewWordResponce(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.AddNewWordResponse)
	return &pb.AddNewWordResponce{
		Id: req.ID,
	}, req.Err
}

// decodeGRPCAddNewWordRequest decode one
func decodeGRPCAddNewWordRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.AddNewWordRequest)
	return endpoints.AddNewWordRequest{
		Word: word.Word{
			W:             req.Word.Word,
			Examples:      req.Word.Examples,
			Transcription: req.Word.Transcription,
		},
	}, nil
}

// encodeGRPCAddNewWordRequest endode one
func encodeGRPCGetByIDResponce(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.GetByIDResponse)
	return &pb.GetByIDResponce{
		Word: &pb.Word{
			Word:          req.Word.W,
			Examples:      req.Word.Examples,
			Transcription: req.Word.Transcription,
		},
	}, req.Err
}

// decodeGRPCAddNewWordRequest decode one
func decodeGRPCGetByIDRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetByIDRequest)
	return endpoints.GetByIDRequest{ID: req.Id}, nil
}
