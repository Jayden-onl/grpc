//
// link: https://grpc.io/docs/quickstart/go.html
// protoc -I proto/ proto/quickstart.proto --go_out=plugins=grpc:proto
//

package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "hongkang.name/grpc/quickstart/proto"
)

const (
	port1 = ":50051"
	port2 = ":50052"
	port3 = ":50053"
)

type server struct {
	Port string
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("%v Received: %v, %v", s.Port, p.Addr.String(), in.Name)
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("%v Received: %v, %v again", s.Port, p.Addr.String(), in.Name)
	}
	return &pb.HelloReply{Message: "Hello " + in.Name + " again"}, nil
}

func (s *server) Say(ctx context.Context, in *pb.World) (*pb.Reply, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("%v Received world: %v, %v", s.Port, p.Addr.String(), in.World)
	}
	return &pb.Reply{Reply: "are you said " + in.World}, nil
}

func main() {
	lis1, err := net.Listen("tcp", port1)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lis2, err := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lis3, err := net.Listen("tcp", port3)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{Port: port1})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve lis1: %v", err)
		} else {

		}
	}(lis1)

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{Port: port2})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve lis2: %v", err)
		} else {

		}
	}(lis2)

	wg.Add(1)
	go func(listener net.Listener) {
		defer wg.Done()
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{Port: port3})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve lis3: %v", err)
		} else {

		}
	}(lis3)
	wg.Wait()
}
