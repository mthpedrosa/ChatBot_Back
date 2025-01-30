package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigRepository struct {
	db *mongo.Database
}

func NewConfigRepository(client *mongo.Client) *ConfigRepository {
	return &ConfigRepository{db: client.Database(dbName)}
}

func (r *ConfigRepository) Insert(ctx context.Context, config models.Config) (string, error) {
	collection := r.db.Collection(configCollection)
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	result, err := collection.InsertOne(ctx, config)
	if err != nil {
		return "", err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get inserted ID")
	}
	return insertedID.Hex(), nil
}

func (r *ConfigRepository) Edit(ctx context.Context, ID string, newConfig models.Config) error {
	collection := r.db.Collection(configCollection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	newConfig.UpdatedAt = time.Now()
	update := bson.M{"$set": newConfig}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *ConfigRepository) Find(ctx context.Context, query string) ([]models.Config, error) {
	collection := r.db.Collection(configCollection)
	filter := bson.M{"name": bson.M{"$regex": query, "$options": "i"}}
	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	var configs []models.Config
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ConfigRepository) FindByID(ctx context.Context, ID string) (*models.Config, error) {
	collection := r.db.Collection(configCollection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	var config models.Config
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigRepository) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(configCollection)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
