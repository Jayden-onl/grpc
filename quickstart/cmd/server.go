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

	"github.com/spf13/pflag"

	pb "hongkang.name/grpc/quickstart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	ports = pflag.StringSlice("ports", []string{":50051"}, "serve ports")
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
	pflag.Parse()

	var wg sync.WaitGroup

	for _, port := range *ports {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		} else {
			log.Printf("listening: %v", port)
		}

		wg.Add(1)
		go func(listener net.Listener, port string) {
			defer wg.Done()
			s := grpc.NewServer()
			pb.RegisterGreeterServer(s, &server{Port: port})
			if err := s.Serve(listener); err != nil {
				log.Fatalf("failed to serve %v: %v", port, err)
			} else {
				log.Printf("serving: %v", port)
			}
		}(lis, port)
	}

	wg.Wait()
}
