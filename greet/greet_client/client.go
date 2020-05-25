package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Haibo",
				LastName:  "Yan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tina",
				LastName:  "Luo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Yan",
				LastName:  "Li",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Zachary",
				LastName:  "Stephen",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tristan",
				LastName:  "Timerlake",
			},
		},
	}

	stream, err :=  c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf(" with error: %v", err)
	}

	for _, request := range requests {
		stream.Send(request)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		if err != nil {
			log.Fatalf("Recv response with error: %v\n", err)
		}
	}
	fmt.Printf("get response message : '%v'\n", res.GetResult())
}



func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doClientStreaming(c)
}
