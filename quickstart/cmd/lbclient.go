package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "hongkang.name/grpc/quickstart/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func main() {
	rso, cleanup := manual.GenerateAndRegisterManualResolver()
	rso.InitialState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:50051"},
			{Addr: "localhost:50052"},
			{Addr: "localhost:50053"},
		},
	})
	defer cleanup()

	conn, err := grpc.DialContext(
		context.Background(),
		rso.Scheme()+":///",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		return
	}

	c := pb.NewGreeterClient(conn)
	defer conn.Close()

	// Contact the server and print out its response.
	name := "Jacob"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	r, err = c.SayHelloAgain(ctx2, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet again: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	for {
		var input string
		fmt.Println("Speak something: ")
		fmt.Scan(&input)

		ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second)
		defer cancel3()
		r, err := c.Say(ctx3, &pb.World{World: input})
		if err != nil {
			log.Fatalf("could not greet again: %v", err)
		}
		log.Printf(r.Reply)
	}
}
