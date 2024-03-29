package main

import (
	"context"
	"fmt"
	calculatorgb "grpcTutorial/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

// SumWithDeadLine implements calculatorgb.CalculatorServiceServer
func (*server) SumWithDeadLine(ctx context.Context, req *calculatorgb.SumRequest) (*calculatorgb.SumResponse, error) {
	log.Println("sum with deadline called...")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("context Canceled ...")
			return nil, status.Errorf(codes.Canceled, "client canceled request")
		}
		time.Sleep(1 * time.Second)
	}

	resp := &calculatorgb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

// Square implements calculatorgb.CalculatorServiceServer
func (*server) Square(ctx context.Context, req *calculatorgb.SquareRequest) (*calculatorgb.SquareResponse, error) {
	log.Println("Square called ...")
	num := req.GetNum()
	if num < 0 {
		log.Printf("req num < 0, return InvalidAgrument")
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, req num was %v", num)
	}

	return &calculatorgb.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}, nil
}

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
	n := req.GetNumber()
	for n > 1 {
		if n%k == 0 {
			n = n / k
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
