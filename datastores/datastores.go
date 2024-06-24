package datastores

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Record represents a generic record in our data store
type Record struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RecordType string             `bson:"record_type" json:"record_type"`
	Data       bson.M             `bson:"data" json:"data"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// Database defines the interface for our data store operations
type Database interface {
	Create(request map[string]interface{}) (*Record, error)
	Read(request map[string]interface{}) (*Record, error)
	Update(request map[string]interface{}) (*Record, error)
	Delete(request map[string]interface{}) (map[string]interface{}, error)
	Search(request map[string]interface{}) ([]Record, error)
	Inventory() (map[string]interface{}, error)
	Close() error
}
