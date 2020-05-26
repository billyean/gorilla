package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet is invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes is invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		response := &greetpb.GreetManyTimesResponse {
			Result: result,
		}
		stream.Send(response)
		time.Sleep(2000 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet is invoked with %v\n", stream)
	result := ""
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			response := &greetpb.LongGreetResponse {
				Result: result,
			}
			return stream.SendAndClose(response)
		}
		if err != nil {
			fmt.Printf("Error happened in recving data: %v\n", err)
		}
		result += "\nHello , " + resp.GetGreeting().GetFirstName() + " " + resp.GetGreeting().GetLastName() + " !"
	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone is invoked with %v\n", stream)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error happened in recving data: %v\n", err)
		}
		resp := &greetpb.GreetEveryoneResponse {
			Result: "Hi, " + req.GetGreeting().FirstName + " " + req.GetGreeting().GetLastName() + " !",
		}
		stream.Send(resp)
	}

	return nil
}

func main() {
	fmt.Printf("Hello world\n")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen (%v)", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
