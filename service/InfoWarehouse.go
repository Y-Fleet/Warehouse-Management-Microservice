package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"warehouse/models"

	pb "github.com/Yfleet/shared_proto/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetInfoWarehouse(client *mongo.Client, rq *pb.InfoWarehouseRequest) (*pb.InfoWarehouseResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	warehouseID := rq.Id
	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		return nil, errors.New("Invalid object id")
	}
	filter := bson.M{"_id": objectID}
	result := client.Database("Warehouse").Collection("Warehouse").FindOne(ctx, filter)
	if result.Err() != nil {
		log.Printf("Error finding warehouse: %v\n", result.Err())
		return nil, result.Err()
	}
	var warehouse models.Warehouse

	err = result.Decode(&warehouse)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var itemIDs []string
	for _, item := range warehouse.Inventory {
		itemIDs = append(itemIDs, item.ID)
		fmt.Println(itemIDs)
	}
	var itemsProtoIDs []*pb.ItemsId
	for _, itemID := range itemIDs {
		itemsProtoIDs = append(itemsProtoIDs, &pb.ItemsId{ID: itemID})
	}
	conn, err := ConnToInventory()
	if err != nil {
		return nil, err
	}
	inventoryClient := pb.NewInventoryServiceClient(conn)
	inventoryResponse, err := inventoryClient.GetInventory(ctx, &pb.GetInventoryRequest{ID: itemsProtoIDs})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("inventoryResponse.Items: %#v\n", inventoryResponse.Items)

	itemMap := make(map[string]*models.Item)
	for _, item := range inventoryResponse.Items {
		itemMap[item.ID] = &models.Item{
			ID:   item.ID,
			Name: item.Name,
			Desc: item.Desc,
			Kg:   item.Kg,
		}
	}

	fmt.Printf("itemMap: %#v\n", itemMap)

	protoWarehouses := ConvertToProtoWarehouses([]*models.Warehouse{&warehouse}, itemMap)
	response := &pb.InfoWarehouseResponse{
		Warehouse: protoWarehouses[0],
	}
	defer conn.Close()
	return response, nil
}
