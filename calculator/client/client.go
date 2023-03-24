package main

import (
	"context"
	calculatorgb "grpcTutorial/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()

	client := calculatorgb.NewCalculatorServiceClient(cc)

	// log.Printf("service client %f", client)
	callSum(client)
}

func callSum(c calculatorgb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorgb.SumRequest{
		Num1: 7,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api err %v", err)
	}

	log.Printf("sum api response %v", resp.GetResult())
}
