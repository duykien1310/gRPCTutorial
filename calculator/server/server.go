package main

import (
	"context"
	"fmt"
	calculatorgb "grpcTutorial/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorgb.SumRequest) (*calculatorgb.SumResponse, error) {
	log.Println("sum called...")
	resp := &calculatorgb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	calculatorgb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running ...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}
