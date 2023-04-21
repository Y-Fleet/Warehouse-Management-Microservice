package service

import (
	"context"
	"errors"
	"log"
	"time"

	models "warehouse/models"

	pb "github.com/Yfleet/shared_proto/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateWarehouse(client *mongo.Client, warehouse *models.Warehouse) error {
	collection := client.Database("Warehouse").Collection("Warehouse")
	_, err := collection.InsertOne(context.Background(), warehouse)
	if err != nil {
		return err
	}
	return nil
}

func ReadWarehouse(client *mongo.Client) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse

	collection := client.Database("Warehouse").Collection("Warehouse")
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var warehouse models.Warehouse
		err := cur.Decode(&warehouse)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return warehouses, nil
}

func UpdateWarehouse(client *mongo.Client, id string, warehouse *models.Warehouse) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid object id")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": warehouse}
	collection := client.Database("Warehouse").Collection("Warehouse")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWarehouse(client *mongo.Client, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid object id")
	}

	filter := bson.M{"_id": objectID}
	collection := client.Database("Warehouse").Collection("Warehouse")
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func ReadWarehouseByID(client *mongo.Client, id string) (*models.Warehouse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid object id")
	}

	filter := bson.M{"_id": objectID}
	collection := client.Database("Warehouse").Collection("Warehouse")
	var warehouse models.Warehouse
	err = collection.FindOne(context.Background(), filter).Decode(&warehouse)
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func AddStock(client *mongo.Client, req *pb.AddStockRequest) (*pb.InventoryItem, error) {
	collection := client.Database("Warehouse").Collection("Warehouse")
	warehouseObjectID, err := primitive.ObjectIDFromHex(req.IDWharehouse)
	if err != nil {
		log.Printf("Error converting warehouseID to ObjectID: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid warehouse ID")
	}

	itemIdObjectID, err := primitive.ObjectIDFromHex(req.IDItems)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		log.Printf("Error converting req.ItemId to ObjectID: %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid item ID")
	}

	filter := bson.M{"_id": warehouseObjectID, "inventory._id": itemIdObjectID}
	update := bson.M{"$inc": bson.M{"inventory.$.quantity": req.Stock}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating warehouse inventory: %v", err)
		return nil, status.Error(codes.Internal, "Error updating warehouse inventory")
	}

	// Find the updated item
	projection := bson.M{"inventory.$": 1}
	result := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection))

	if result.Err() != nil {
		log.Printf("Error finding the updated item: %v", result.Err())
		return nil, status.Error(codes.Internal, "Error finding the updated item")
	}

	var updatedWarehouse struct {
		Inventory []pb.InventoryItem `bson:"inventory"`
	}

	err = result.Decode(&updatedWarehouse)
	if err != nil {
		log.Printf("Error decoding the updated item: %v", err)
		return nil, status.Error(codes.Internal, "Error decoding the updated item")
	}

	if len(updatedWarehouse.Inventory) > 0 {
		return &updatedWarehouse.Inventory[0], nil
	} else {
		return nil, status.Error(codes.NotFound, "Updated item not found")
	}
}

func DelStock(client *mongo.Client, req *pb.AddStockRequest) (bool, error) {
	return false, nil
}
