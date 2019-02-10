package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/dictionary/proto"
	"google.golang.org/grpc"

	"github.com/dictionary/models/word"

	"github.com/dictionary/endpoints"
	"github.com/dictionary/services"
	"github.com/dictionary/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Addr this is address our server
var Addr = flag.String("address", ":8080", "address http server")

// AddrGRPC this is address of our grpc server
var AddrGRPC = flag.String("gRPCaddr", ":8081", "handled address by gRPC serv")

// StorageDir it's path to dir for storage
var StorageDir = flag.String("storage", "./storage", "path to storage dir")

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger,
			"service", "dictionary",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	logger.Log("info", "server started")
	defer logger.Log("info", "server shoutdowned")

	var service services.Dictionary
	{
		repository := word.NewBadgerRepository(*StorageDir, logger)
		service = services.NewDictionary(repository, logger)
	}

	endpoints := endpoints.NewEndpoints(service)
	var h http.Handler
	{
		h = transport.NewService(endpoints, logger)
	}

	errsChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		errsChan <- fmt.Errorf("%c", <-c)
	}()

	server := http.Server{
		Addr:    *Addr,
		Handler: h,
	}

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *Addr)
		errsChan <- server.ListenAndServe()
	}()

	go grpcListen(errsChan, endpoints, logger)

	level.Error(logger).Log("exit", <-errsChan)
}

func grpcListen(errChan chan<- error, endpoints endpoints.Endpoints, logger log.Logger) {
	listener, err := net.Listen("tcp", *AddrGRPC)
	if err != nil {
		errChan <- err
		return
	}

	srv := transport.NewGRPCService(context.Background(), endpoints)
	grpcServ := grpc.NewServer()
	pb.RegisterDictionaryServer(grpcServ, srv)

	logger.Log("stars grpc listening")

	errChan <- grpcServ.Serve(listener)

	logger.Log("end grpc listening")
}
