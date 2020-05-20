package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Haibo",
			LastName: "Yan",
		},
	}
	response, err :=  c.Greet(context.Background(), request)
	if err != nil {
		log.Fatal("response with error: %v", err)
	}

	fmt.Printf("get response message : '%v'", response.Result)
}
