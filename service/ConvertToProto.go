package service

import (
	"warehouse/models"

	
	pb "github.com/Yfleet/shared_proto/api"
)

func ConvertToProtoWarehouses(modelWarehouses []*models.Warehouse, itemMap map[string]*models.Item) []*pb.Warehouse {
	protoWarehouses := make([]*pb.Warehouse, len(modelWarehouses))
	for i, modelWarehouse := range modelWarehouses {
		protoInventory := make([]*pb.InventoryItem, len(modelWarehouse.Inventory))
		for j, modelInventoryItem := range modelWarehouse.Inventory {
			item := itemMap[modelInventoryItem.ID]
			if item != nil {
				protoItem := &pb.ItemWaherouse{
					ID:   item.ID,
					Name: item.Name,
					Desc: item.Desc,
					Kg:   item.Kg,
				}
				protoInventory[j] = &pb.InventoryItem{
					ItemId:   modelInventoryItem.ID,
					Quantity: modelInventoryItem.Quantity,
					InfoItem: protoItem,
				}
			} else {
				protoInventory[j] = &pb.InventoryItem{
					ItemId:   modelInventoryItem.ID,
					Quantity: modelInventoryItem.Quantity,
				}
			}
		}
		protoWarehouse := &pb.Warehouse{
			Id:        modelWarehouse.ID,
			Name:      modelWarehouse.Name,
			Address:   modelWarehouse.Address,
			Latitude:  modelWarehouse.Location.Coordinates[1],
			Longitude: modelWarehouse.Location.Coordinates[0],
			Inventory: protoInventory,
		}
		protoWarehouses[i] = protoWarehouse
	}
	return protoWarehouses
}
