package orchestrator

import (
	"context"
	"log"

	pb "github.com/1minepowminx/distributed_calculator/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToAgent() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Could not connect to grpc server:", err)
		return nil, err
	}

	log.Println("Connection status:", conn.GetState())
	return conn, nil
}

// Соединение с агентом
func Calculate(ctx context.Context, expr string) (float64, error) {
	conn, err := ConnectToAgent()
	if err != nil {
		log.Fatalf("Failed connect to agent: %v", err)
		return 0, err
	}

	grpcClient := pb.NewCalculatorServiceClient(conn)

	ans1, err := grpcClient.Calculate(ctx, &pb.ExpressionRequest{Expression: expr})
	if err != nil {
		log.Printf("Failed get answer from agent: %v", err)
		return 0, err
	}

	return ans1.Result, nil
}
