package dict

import (
	"log"
	"time"

	pb "github.com/dictionary/tgbot/dict/pb"

	"google.golang.org/grpc"
)

// NewDict create gRPC connection with service and return one
func NewDict(addr string) pb.DictionaryClient {
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	_panicIfErorr(err)
	log.Printf("dict creasdasdated")

	client := pb.NewDictionaryClient(conn)

	return client
}

func _panicIfErorr(err error) {
	if err != nil {
		panic(err)
	}
}
