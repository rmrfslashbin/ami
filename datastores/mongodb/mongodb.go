package mongodb

import (
	"context"
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

	now := time.Now().UTC() // Use UTC time
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

	// Handle date-based queries
	if createdAt, ok := query["created_at"].(map[string]interface{}); ok {
		for op, value := range createdAt {
			if timeStr, ok := value.(string); ok {
				t, err := time.Parse(time.RFC3339, timeStr)
				if err != nil {
					return nil, &ErrInvalidDateFormat{Field: "created_at", Err: err}
				}
				createdAt[op] = t
			}
		}
		query["created_at"] = createdAt
	}

	// Construct the final query
	finalQuery := bson.M{}
	for key, value := range query {
		if key == "record_type" {
			finalQuery[key] = value
		} else {
			finalQuery["data."+key] = value
		}
	}

	cursor, err := ds.db.Collection("records").Find(context.Background(), finalQuery)
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

		// Convert primitive.DateTime to time.Time
		oldestRecord := result["oldest_record"].(primitive.DateTime).Time()
		newestRecord := result["newest_record"].(primitive.DateTime).Time()

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
