package main

import (
	"context"
	pb "github.com/renanvicente/grpc_sample/grpc/keyvalue"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedKeyValueServer
}

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	// The local Get function is implemented back in Chapter 5
	value, err := Get(r.Key)

	// Return expects a GetResponse pointer and an err
	return &pb.GetResponse{Value: value}, err
}

func (s *server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received PUT key=%v value=%v", r.Key, r.Value)

	return &pb.PutResponse{}, Put(r.Key, r.Value)

}

func (s *server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Received DELETE key=%v", r.Key)

	return &pb.DeleteResponse{}, Delete(r.Key)
}

func main() {
	// Create a gRPC server and register our KeyValueServer with it
	s := grpc.NewServer()
	pb.RegisterKeyValueServer(s, &server{})

	// Open a listening port on 50051
	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	// Start accepting connections on the listening port
	if err := s.Serve(list); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
