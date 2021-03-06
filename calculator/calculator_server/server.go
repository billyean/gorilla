package main

import (
	"context"
	"fmt"
	"github.com/billyean/tornado/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}


func (*server) Calculate(ctx context.Context, req *calculatorpb.CalRequest) (*calculatorpb.CalResponse, error)  {
	operand1 := req.GetOperand1()
	operand2 := req.GetOperand2()
	operator := req.GetOperator()

	fmt.Printf("Calling %d %v %d\n", operand1, operator, operand2)
	var result int32
	switch operator {
		case calculatorpb.CalRequest_ADD: result = operand1 + operand2
		case calculatorpb.CalRequest_MINUS: result = operand1 - operand2
		case calculatorpb.CalRequest_TIMES: result = operand1 * operand2
		case calculatorpb.CalRequest_DIVIDE: result = operand1 / operand2
		case calculatorpb.CalRequest_MOD: result = operand1 % operand2
	}
	res := &calculatorpb.CalResponse{
		Result: result,
	}
	return res, nil
}

func (*server)  FibonacciNumber(req *calculatorpb.FibonacciRequest, stream calculatorpb.CalculateService_FibonacciNumberServer) error {
	var firstNumber = req.GetFirstNumber()
	var secondNumber = req.GetSecondNumber()
	var n int32 = 0

	var result int64

	for {
		n += 1
		result = firstNumber + secondNumber
		fmt.Printf("Calling %d + %d = %d\n", firstNumber, secondNumber, result)
		response := &calculatorpb.FibonacciResponse{
			Result: result,
			N: n,
		}
		err := stream.Send(response)
		if err != nil {
			log.Fatalf("Error happened when sending response: %v\n", err)
			break
		}
		time.Sleep(200 * time.Millisecond)
		firstNumber = secondNumber
		secondNumber = result
	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculateService_AverageServer) error {
	n := 0
	sum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error when receiving request: %v\n", err)
		}
		n += 1
		sum += req.GetNumber()
 	}

 	response := &calculatorpb.AverageResponse{
 		Average: float64(sum) / float64(n),
	}
	stream.SendAndClose(response)

	return nil
}


func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen (%v)", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
