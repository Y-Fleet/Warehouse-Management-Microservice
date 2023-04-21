package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	service "warehouse/service"

	pb "github.com/Y-Fleet/Grpc-Api/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedWarehouseManagementServiceServer
}

func (s *server) RenderWarehouse(ctx context.Context, req *pb.RenderWarehouseRequest) (*pb.RenderWarehouseResponse, error) {
	client, cancel := ConToDb()
	DataWarehouse, err := service.ReadWarehouse(client)
	if err != nil {
		return nil, err
	}
	protoWarehouses := service.ConvertToProtoWarehouses(DataWarehouse, nil)
	cancel()
	return &pb.RenderWarehouseResponse{Warehouse: protoWarehouses}, nil
}

func (s *server) InfoWarehouse(ctx context.Context, req *pb.InfoWarehouseRequest) (*pb.InfoWarehouseResponse, error) {
	client, cancel := ConToDb()
	cancel()
	response, err := service.GetInfoWarehouse(client, req)
	return &pb.InfoWarehouseResponse{Warehouse: response.Warehouse}, err
}

func (s *server) AddStock(ctx context.Context, req *pb.AddStockRequest) (*pb.AddStockResponse, error) {
	client, cancel := ConToDb()
	response, err := service.AddStock(client, req)
	cancel()
	if err != nil {
		return nil, err
	}
	return &pb.AddStockResponse{Item: response}, nil
}

func ConToDb() (*mongo.Client, context.CancelFunc) {
	const connectionString = "mongodb://Stock:Stock@localhost:27021"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	fmt.Println("Connected successfully to MongoDB")
	return client, cancel
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWarehouseManagementServiceServer(s, &server{})
	reflection.Register(s)
	log.Println("Starting microservice on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
