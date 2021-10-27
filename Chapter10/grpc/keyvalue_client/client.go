package main

import (
	"context"
	pb "github.com/renanvicente/CloudNativeGo/Chapter10/grpc/keyvalue"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// Use context to set a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We can still set other options as desired.
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

	conn, err := grpc.DialContext(ctx, "127.0.0.1:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		//grpclog.Fatalf(err)
	}
	defer conn.Close()

	//conn, err := grpc.Dial("127.0.0.1:50051",
	//	grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()

	// Get a new instance of our client
	client := pb.NewKeyValueClient(conn)

	var action, key, value string

	// Expect something like "set foo bar"
	if len(os.Args) > 2 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}

	//// Use context to establish a 1-second timeout.
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// Call client.Get() or client.Put() as appropriate.
	switch action {
	case "get":
		r, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}
		log.Printf("Get %s returns: %s", key, r.Value)
	case "put":
		_, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("could not put key %s: %v\n", key, err)
		}
		log.Printf("Put %s", key)
	case "delete":
		_, err := client.Delete(ctx, &pb.DeleteRequest{Key: key})
		if err != nil {
			log.Fatalf("could not delete key %s: %v\n", key, err)
		}
		log.Printf("Delete %s", key)
	default:
		log.Fatalf("Syntax: go run [get|put] KEY VALUE...")
	}
}
