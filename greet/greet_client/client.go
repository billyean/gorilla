package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func doServerRequest(c greetpb.GreetServiceClient) {
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
func doServerStreaming(c greetpb.GreetServiceClient) {
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Haibo",
			LastName: "Yan",
		},
	}

	greetManyTimeClient, err :=  c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatal("response with error: %v\n", err)
	}

	for {
		response, err :=  greetManyTimeClient.Recv()
		if (err == io.EOF) {
			greetManyTimeClient.CloseSend()
			break
		}
		if err != nil {
			log.Fatal("error {} happened when Reciving messages: %v\n", err)
		}
		fmt.Printf("get response message : '%v'\n", response.Result)
	}

}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doServerStreaming(c)
}
