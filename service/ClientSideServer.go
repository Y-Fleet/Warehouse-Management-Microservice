package service

import (
	"log"

	"google.golang.org/grpc"
)

func ConnToInventory() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conn, nil
}
