package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func compute(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.CalRequest{
		Operand1: 20,
		Operand2: 10,
		Operator: calculatorpb.CalRequest_MOD,
	}

	response, err :=  c.Calculate(context.Background(), request)
	if err != nil {
		log.Fatal("response with error: %v", err)
	}

	fmt.Printf("get response message : '%d'", response.Result)
}

func fibonacci(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.FibonacciRequest{
		FirstNumber: 1,
		SecondNumber: 1,
	}
	client, err := c.FibonacciNumber(context.Background(), request)
	if err != nil {
		log.Fatal("response with error: %v", err)
	}
	for {
		response, err := client.Recv()
		if err != nil {
			log.Fatal("error {} happened when Reciving messages: %v\n", err)
		}
		fmt.Printf("get %vth fibonnaci number : %v\n", response.N, response.Result)
	}
}



func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	if err != nil {
		log.Fatal("could not connect: %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculateServiceClient(conn)
	fibonacci(c)
}
