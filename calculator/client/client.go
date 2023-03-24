package main

import (
	"context"
	calculatorgb "grpcTutorial/calculator/calculatorpb"
	"io"
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
	// callSum(client)
	// callPND(client)
	callAverage(client)
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

func callPND(c calculatorgb.CalculatorServiceClient) {
	log.Println("calling PND api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorgb.PNDRequest{
		Number: 120,
	})
	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, errErr := stream.Recv()

		if errErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("prime number %v", resp.GetResult())
	}
}

func callAverage(c calculatorgb.CalculatorServiceClient) {
	log.Println("calling Average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("call average err %v", err)
	}

	listReq := []calculatorgb.AverageRequest{
		calculatorgb.AverageRequest{
			Num: 5,
		},
		calculatorgb.AverageRequest{
			Num: 10,
		},
		calculatorgb.AverageRequest{
			Num: 12,
		},
		calculatorgb.AverageRequest{
			Num: 3,
		},
		calculatorgb.AverageRequest{
			Num: 4.2,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("send average request err %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response err %v", err)
	}
	log.Printf("average response %+v", resp)
}
