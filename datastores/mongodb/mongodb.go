package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rmrfslashbin/ami/datastores"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataStore represents our data store and contains the MongoDB client
type DataStore struct {
	mongoDbURI   string
	databaseName string
	client       *mongo.Client
	db           *mongo.Database
}

// Option is a configuration option.
type Option func(config *DataStore)

// New creates a new Messages configuration.
func New(opts ...func(*DataStore)) (*DataStore, error) {
	dataStore := &DataStore{}
	for _, opt := range opts {
		opt(dataStore)
	}

	if dataStore.mongoDbURI == "" {
		return nil, &ErrMissingMongoDbURI{}
	}

	if dataStore.databaseName == "" {
		return nil, &ErrMissingDatabaseName{}
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dataStore.mongoDbURI))
	if err != nil {
		return nil, &ErrFailedToConnectToMongoDB{Err: err}
	}

	dataStore.client = client
	dataStore.db = client.Database(dataStore.databaseName)

	return dataStore, nil
}

// WithMongoDbURI sets the MongoDB URI for the data store
func WithMongoDbURI(uri string) Option {
	return func(config *DataStore) {
		config.mongoDbURI = uri
	}
}

// WithDatabaseName sets the database name for the data store
func WithDatabaseName(name string) Option {
	return func(config *DataStore) {
		config.databaseName = name
	}
}

// GetClaudeToolSpec returns the Claude Tool specification as a blockquote string
func (ds *DataStore) GetClaudeToolSpec() string {
	spec := `{
  "type": "function",
  "function": {
    "name": "data_store",
    "description": "Interact with a persistent data store to create, read, update, and delete records, as well as search and retrieve inventory information.",
    "parameters": {
      "type": "object",
      "properties": {
        "action": {
          "type": "string",
          "enum": ["create", "read", "update", "delete", "search", "inventory"],
          "description": "The action to perform on the data store"
        },
        "record_type": {
          "type": "string",
          "description": "The type of record (e.g., 'person', 'company', 'product')"
        },
        "id": {
          "type": "string",
          "description": "The unique identifier of the record (required for read, update, and delete actions)"
        },
        "data": {
          "type": "object",
          "description": "The data to be stored or updated (required for create and update actions)"
        },
        "query": {
          "type": "object",
          "description": "The search query parameters (required for search action)"
        },
        "fields": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Specific fields to update (optional for update action)"
        }
      },
      "required": ["action"]
    }
  }
}`
	return fmt.Sprintf("> %s", spec)
}

// HandleRequest processes incoming requests from Claude
func HandleRequest(db datastores.Database, input []byte) ([]byte, error) {
	var request map[string]interface{}
	err := json.Unmarshal(input, &request)
	if err != nil {
		return nil, &ErrFailedToUnmarshalRequest{Err: err}
	}

	action, ok := request["action"].(string)
	if !ok {
		return nil, &ErrActionMissingOrInvalid{}
	}

	var result interface{}

	switch action {
	case "create":
		result, err = db.Create(request)
	case "read":
		result, err = db.Read(request)
	case "update":
		result, err = db.Update(request)
	case "delete":
		result, err = db.Delete(request)
	case "search":
		result, err = db.Search(request)
	case "inventory":
		result, err = db.Inventory()
	default:
		err = &ErrActionUnknown{Action: action}
	}

	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

// Create creates a new record in the data store
func (ds *DataStore) Create(request map[string]interface{}) (*datastores.Record, error) {
	recordType, ok := request["record_type"].(string)
	if !ok {
		return nil, &ErrRecordTypeMissingOrInvalid{}
	}

	data, ok := request["data"].(map[string]interface{})
	if !ok {
		return nil, &ErrDataMissingOrInvalid{}
	}

	now := time.Now()
	record := &datastores.Record{
		RecordType: recordType,
		Data:       data,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result, err := ds.db.Collection("records").InsertOne(context.Background(), record)
	if err != nil {
		return nil, &ErrFailedToInsertRecord{Err: err}
	}

	record.ID = result.InsertedID.(primitive.ObjectID)
	return record, nil
}

// Close closes the MongoDB connection
func (ds *DataStore) Close() error {
	return ds.client.Disconnect(context.Background())
}

// Read retrieves a record from the data store
func (ds *DataStore) Read(request map[string]interface{}) (*datastores.Record, error) {
	idStr, ok := request["id"].(string)
	if !ok {
		return nil, &ErrIdMissingOrInvalid{}
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, &ErrIdMissingOrInvalid{Err: err}
	}

	var record datastores.Record
	err = ds.db.Collection("records").FindOne(context.Background(), bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &ErrRecordNotFound{}
		}
		return nil, &ErrFailedToRetrieveRecord{Err: err}
	}

	return &record, nil
}

// Update updates an existing record in the data store
func (ds *DataStore) Update(request map[string]interface{}) (*datastores.Record, error) {
	idStr, ok := request["id"].(string)
	if !ok {
		return nil, &ErrIdMissingOrInvalid{}
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, &ErrIdMissingOrInvalid{Err: err}
	}

	data, ok := request["data"].(map[string]interface{})
	if !ok {
		return nil, &ErrDataMissingOrInvalid{}
	}

	update := bson.M{
		"$set": bson.M{
			"data":       data,
			"updated_at": time.Now(),
		},
	}

	var record datastores.Record
	err = ds.db.Collection("records").FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&record)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &ErrRecordNotFound{}
		}
		return nil, &ErrFailedToUpdateRecord{Err: err}
	}

	return &record, nil
}

// Delete removes a record from the data store
func (ds *DataStore) Delete(request map[string]interface{}) (map[string]interface{}, error) {
	idStr, ok := request["id"].(string)
	if !ok {
		return nil, &ErrIdMissingOrInvalid{}
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, &ErrIdMissingOrInvalid{Err: err}
	}

	result, err := ds.db.Collection("records").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, &ErrFailedToDeleteRecord{Err: err}
	}

	if result.DeletedCount == 0 {
		return nil, &ErrRecordNotFound{}
	}

	return map[string]interface{}{
		"message": "Record deleted successfully",
		"id":      idStr,
	}, nil
}

// Search searches for records in the data store
func (ds *DataStore) Search(request map[string]interface{}) ([]datastores.Record, error) {
	query, ok := request["query"].(map[string]interface{})
	if !ok {
		return nil, &ErrQueryMissingOrInvalid{}
	}

	cursor, err := ds.db.Collection("records").Find(context.Background(), query)
	if err != nil {
		return nil, &ErrFailedToSearchRecords{Err: err}
	}
	defer cursor.Close(context.Background())

	var records []datastores.Record
	if err = cursor.All(context.Background(), &records); err != nil {
		return nil, &ErrFailedToDecodeSearchResults{Err: err}
	}

	return records, nil
}

// Inventory retrieves a summary of the data store contents
func (ds *DataStore) Inventory() (map[string]interface{}, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":           "$record_type",
			"count":         bson.M{"$sum": 1},
			"fields":        bson.M{"$addToSet": "$data"},
			"oldest_record": bson.M{"$min": "$created_at"},
			"newest_record": bson.M{"$max": "$updated_at"},
		}},
	}

	cursor, err := ds.db.Collection("records").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, &ErrFailedToExecuteInventoryAggregation{Err: err}
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, &ErrFailedToDecodeInventoryResults{Err: err}
	}

	inventory := make(map[string]interface{})
	for _, result := range results {
		recordType := result["_id"].(string)
		count := result["count"].(int32)
		fields := result["fields"].(primitive.A)
		oldestRecord := result["oldest_record"].(time.Time)
		newestRecord := result["newest_record"].(time.Time)

		// Get a sample of field names
		sampleFields := make(map[string]bool)
		for _, f := range fields {
			for k := range f.(bson.M) {
				sampleFields[k] = true
			}
		}

		inventory[recordType] = map[string]interface{}{
			"count":         count,
			"sample_fields": keys(sampleFields),
			"oldest_record": oldestRecord,
			"newest_record": newestRecord,
		}
	}

	return inventory, nil
}

// keys is a helper function to get keys from a map
func keys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
