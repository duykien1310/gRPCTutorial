package main

import (
	"context"
	"fmt"
	calculatorgb "grpcTutorial/calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// FindMax implements calculatorgb.CalculatorServiceServer
func (*server) FindMax(stream calculatorgb.CalculatorService_FindMaxServer) error {
	log.Println("Find max called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF ...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while Recv Find Max %v", err)
			return err
		}

		num := req.GetNum()
		log.Printf("recv num %v", num)
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorgb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
	}
}

// Average implements calculatorgb.CalculatorServiceServer
func (*server) Average(stream calculatorgb.CalculatorService_AverageServer) error {
	log.Println("Average called")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//tinh trung binh va return cho client
			resp := calculatorgb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(&resp)
		}
		if err != nil {
			log.Fatalf("err while Recv Average %v", err)
		}

		log.Println("receive num")
		total += req.GetNum()
		count++
	}
}

// PrimeNumberDecomposition implements calculatorgb.CalculatorServiceServer
func (*server) PrimeNumberDecomposition(req *calculatorgb.PNDRequest, stream calculatorgb.CalculatorService_PrimeNumberDecompositionServer) error {

	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			// Send to client
			stream.Send(&calculatorgb.PNDResponse{
				Result: k,
			})
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

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
