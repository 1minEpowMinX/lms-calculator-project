package agent

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/1minepowminx/distributed_calculator/internal/utils/agent/calculation"
	itp "github.com/1minepowminx/distributed_calculator/internal/utils/agent/infix_to_postfix"
	"github.com/1minepowminx/distributed_calculator/internal/utils/agent/validator"
	pb "github.com/1minepowminx/distributed_calculator/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.CalculatorServiceServer
}

func NewServer() *Server {
	return &Server{}
}

type CalculatorServiceServer interface {
	Calculate(context.Context, *pb.ExpressionRequest) (*pb.ExpressionResponse, error)
	mustEmbedUnimplementedCalculatorServiceServer()
}

func (s *Server) Calculate(ctx context.Context, in *pb.ExpressionRequest) (*pb.ExpressionResponse, error) {
	if !validator.IsValidExpression(in.Expression) {
		return nil, fmt.Errorf("invalid expression: %s", in.Expression)
	}
	postfixed := itp.ToPostfix(in.Expression)
	res, err := calculation.Evaluate(postfixed)
	if err != nil {
		return nil, err
	}
	log.Println("Successful operation!")
	return &pb.ExpressionResponse{Result: res}, nil
}

func RunAgentServer() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%v:%v", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error starting TCP listener: ", err)
		os.Exit(1)
	}

	log.Printf("TCP listener started at %s", addr)

	grpcServer := grpc.NewServer()

	expressionServiceServer := NewServer()

	pb.RegisterCalculatorServiceServer(grpcServer, expressionServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("Error serving grpc: ", err)
		os.Exit(1)
	}
}
