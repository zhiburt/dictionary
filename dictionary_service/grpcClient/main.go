package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/dictionary/endpoints"

	pb "github.com/dictionary/proto"
	"google.golang.org/grpc"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// Addr addr to grpc server
var Addr = flag.String("Address", ":8083", "Address for connection")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*Addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	_panicIfErorr(err)

	client := pb.NewDictionaryClient(conn)
	ctx := context.Background()
	{
		req1 := &pb.AddNewWordRequest{
			Word: &pb.Word{Word: "c++"},
		}

		resp1, err := client.AddNewWord(ctx, req1)
		_panicIfErorr(err)

		req2 := &pb.GetByIDRequest{Id: resp1.Id}
		resp2, err := client.GetByID(ctx, req2)
		_panicIfErorr(err)
		fmt.Println(resp2.Word)
	}
}

// NewClient return clients endpoints
// But Why?
func NewClient(conn *grpc.ClientConn) endpoints.Endpoints {
	var dictAddNewWordMethod = grpctransport.NewClient(conn,
		"Dict", "AddNewWord",
		encodeGRPCRequest,
		decodeGRPCResponce,
		&pb.AddNewWordResponce{},
	).Endpoint()

	var dictGetByIDMethod = grpctransport.NewClient(conn,
		"Dict", "GetByID",
		encodeGRPCRequest,
		decodeGRPCResponce,
		&pb.GetByIDResponce{},
	).Endpoint()

	return endpoints.Endpoints{
		AddNewWord: dictAddNewWordMethod,
		GetByID:    dictGetByIDMethod,
	}
}

func encodeGRPCRequest(_ context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func decodeGRPCResponce(_ context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func _panicIfErorr(err error) {
	if err != nil {
		panic(err)
	}
}
