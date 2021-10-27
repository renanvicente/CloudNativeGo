package frontend

import (
	"context"
	pb "github.com/renanvicente/CloudNativeGo/Chapter10/grpc/keyvalue"
	"github.com/renanvicente/CloudNativeGo/Chapter10/hexarch/core"
	"google.golang.org/grpc"
	"log"
	"net"
)

// grpcFrontEnd contains a reference to the core application logic,
// and complies with the contract defined by the FrontEnd interface.
type grpcFrontEnd struct {
	store *core.KeyValueStore
	pb.UnimplementedKeyValueServer
}

func (s *grpcFrontEnd) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	// The local Get function is implemented back in Chapter 5
	value, err := s.store.Get(r.Key)

	// Return expects a GetResponse pointer and an err
	return &pb.GetResponse{Value: value}, err
}

func (s *grpcFrontEnd) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received PUT key=%v value=%v", r.Key, r.Value)
	return &pb.PutResponse{}, s.store.Put(r.Key, r.Value)

}

func (s *grpcFrontEnd) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Received DELETE key=%v", r.Key)

	return &pb.DeleteResponse{}, s.store.Delete(r.Key)
}

func (f *grpcFrontEnd) Start(store *core.KeyValueStore) error {
	// Remember our core application reference.
	f.store = store
	// Create a gRPC server and register our KeyValueServer with it
	s := grpc.NewServer()
	pb.RegisterKeyValueServer(s, f)

	// Open a listening port on 50051
	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	// Start accepting connections on the listening port
	return s.Serve(list)
	//if err := s.Serve(list); err != nil {
	//	log.Fatal("failed to serve: %v", err)
	//}
}
