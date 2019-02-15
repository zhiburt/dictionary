package transport

import (
	"context"
	"fmt"

	"github.com/dictionary/dictionary_service/models/word"

	"github.com/dictionary/dictionary_service/endpoints"
	pb "github.com/dictionary/dictionary_service/proto"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	words      grpctransport.Handler
	addNewWord grpctransport.Handler
	getByID    grpctransport.Handler
	getByW     grpctransport.Handler
}

func (g *grpcServer) Words(ctx context.Context, req *pb.WordsRequest) (*pb.WordsResponce, error) {
	_, resp, err := g.words.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.WordsResponce), nil
}

func (g *grpcServer) AddNewWord(ctx context.Context, req *pb.AddNewWordRequest) (*pb.AddNewWordResponce, error) {
	fmt.Println("asdasdasdasd")
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

func (g *grpcServer) GetByW(ctx context.Context, req *pb.GetByWRequest) (*pb.GetByWResponce, error) {
	_, resp, err := g.getByW.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByWResponce), nil
}

// NewGRPCService create and return gprc service
func NewGRPCService(ctx context.Context, endpoint endpoints.Endpoints) pb.DictionaryServer {
	return &grpcServer{
		grpctransport.NewServer(
			endpoint.Words,
			decodeGRPCWordsRequest,
			encodeGRPCWordsResponce,
		),
		grpctransport.NewServer(
			endpoint.AddNewWord,
			decodeGRPCAddNewWordRequest,
			encodeGRPCAddNewWordResponce,
		),
		grpctransport.NewServer(
			endpoint.GetByID,
			decodeGRPCGetByIDRequest,
			encodeGRPCGetByIDResponce,
		),
		grpctransport.NewServer(
			endpoint.GetByW,
			decodeGRPCGetByWRequest,
			encodeGRPCGetByWResponce,
		),
	}
}

// encodeGRPCAddNewWordRequest endode one
func encodeGRPCWordsResponce(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.WordsResponse)
	var responceWords []*pb.Word
	for _, w := range req.Words {
		responceWords = append(responceWords, &pb.Word{
			Word:          w.W,
			Transcription: w.Transcription,
			Examples:      w.Examples,
		})
	}
	return &pb.WordsResponce{Words: responceWords}, req.Err
}

// decodeGRPCAddNewWordRequest decode one
func decodeGRPCWordsRequest(ctx context.Context, r interface{}) (interface{}, error) {
	_, ok := r.(*pb.WordsRequest)
	if ok {
		return endpoints.WordsRequest{}, nil
	}
	return nil, fmt.Errorf("cannt decode")
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

// encodeGRPCAddNewWordRequest endode one
func encodeGRPCGetByWResponce(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.GetByWResponse)
	return &pb.GetByWResponce{
		Word: &pb.Word{
			Word:          req.Word.W,
			Examples:      req.Word.Examples,
			Transcription: req.Word.Transcription,
		},
	}, req.Err
}

// decodeGRPCAddNewWordRequest decode one
func decodeGRPCGetByWRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetByWRequest)
	return endpoints.GetByWRequest{W: req.W}, nil
}
