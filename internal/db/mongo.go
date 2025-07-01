package db

import (
	"Jsonb/data"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// MongoDB connection setup
func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	MongoClient = client
	MongoDB = client.Database(os.Getenv("MONGO_DB"))

	// Create indexes equivalent to PostgreSQL
	collection := MongoDB.Collection("json_data")
	
	// Index equivalent to PostgreSQL: (data->>'id', id)
	// In MongoDB: compound index on (id, _id)
	idIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
			{Key: "_id", Value: 1},
		},
		Options: options.Index().SetName("idx_id_objectid"),
	}

	// Index equivalent to PostgreSQL: ((data->>'age')::int)  
	// In MongoDB: simple index on age field
	ageIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "age", Value: 1},
		},
		Options: options.Index().SetName("idx_age"),
	}

	indexModels := []mongo.IndexModel{idIndexModel, ageIndexModel}
	
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	
	_, err = collection.Indexes().CreateMany(ctx2, indexModels)
	if err != nil {
		return err
	}

	fmt.Println("✅ Index on (id, _id) created")
	fmt.Println("✅ Index on age created")
	fmt.Println("✅ Connected to MongoDB")
	return nil
}

// Create single document
func CreateMongo(record data.DummyData) error {
	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, record)
	return err
}

// Bulk insert multiple documents (optimized)
func InsertMongoBulk(records []data.DummyData) error {
	if len(records) == 0 {
		return nil
	}

	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert to interface slice
	docs := make([]interface{}, len(records))
	for i, record := range records {
		docs[i] = record
	}

	// Use ordered=false for better performance
	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(ctx, docs, opts)
	return err
}

// Insert multiple documents (kept for backward compatibility)
func InsertMongo(records []data.DummyData) error {
	return InsertMongoBulk(records)
}

// Find single document by 'id' - uses idx_id_objectid index
func FindMongoByID(id int) (data.DummyData, error) {
	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result data.DummyData
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	return result, err
}

// Query by age condition - uses idx_age index
func QueryMongo(minAge int) ([]data.DummyData, error) {
	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"age": bson.M{"$gt": minAge}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []data.DummyData
	err = cursor.All(ctx, &results)
	return results, err
}

// Update documents where age > minAge - uses idx_age index
func UpdateMongoCity(minAge int, newCity string) error {
	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"age": bson.M{"$gt": minAge}}
	update := bson.M{"$set": bson.M{"address.city": newCity}}

	_, err := collection.UpdateMany(ctx, filter, update)
	return err
}

// Delete documents where age < maxAge - uses idx_age index
func DeleteMongoByAge(maxAge int) error {
	collection := MongoDB.Collection("json_data")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"age": bson.M{"$lt": maxAge}}

	_, err := collection.DeleteMany(ctx, filter)
	return err
}
