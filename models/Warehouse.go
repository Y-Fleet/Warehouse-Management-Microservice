package models

type Warehouse struct {
	ID        string          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string          `json:"name,omitempty" bson:"name,omitempty"`
	Address   string          `json:"address,omitempty" bson:"address,omitempty"`
	Location  Location        `json:"location,omitempty" bson:"location,omitempty"`
	Inventory []InventoryItem `json:"inventory,omitempty" bson:"inventory,omitempty"`
}

type Location struct {
	Type        string    `js"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

type InventoryItem struct {
	ID   string `json:"_id,omitempty" bson:"_id,omitempty"`
	Quantity int32  `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type Item struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Desc string `bson:"desc"`
	Kg   int32  `bson:"Kg"`
}
