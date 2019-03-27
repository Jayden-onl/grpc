package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "hongkang.name/grpc/quickstart/proto"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "Jacob"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
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
