package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OpenaiMongo struct {
	db *mongo.Database
}

func NewOpenAiMongoRepository(db *mongo.Client) *OpenaiMongo {
	return &OpenaiMongo{db: db.Database(dbName)}
}

// Insert the Meta Account in the database
func (o *OpenaiMongo) Insert(ctx context.Context, assistant models.CreateAssistant) (string, error) {
	collection := o.db.Collection(openaiCollection)

	assistant.CreatedAt = time.Now().Add(-3 * time.Hour)

	// Insert the account into the collection.
	result, err := collection.InsertOne(ctx, assistant)
	if err != nil {
		return "", err
	}

	// The ID generated for the created account.
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("falha ao obter o ID da conta criada inserido")
	}

	return insertedID.Hex(), nil
}

func (o *OpenaiMongo) Edit(ctx context.Context, ID string, assistant models.CreateAssistant) error {
	collection := o.db.Collection(openaiCollection)

	assistant.UpdateAt = time.Now().Add(-3 * time.Hour)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	// Converts the novaConversa struct to a map to use in $set
	updateFields := bson.M{}
	bsonBytes, err := bson.Marshal(assistant)
	if err != nil {
		return err
	}
	bson.Unmarshal(bsonBytes, &updateFields)

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (o *OpenaiMongo) FindAllUser(ctx context.Context, ID string) ([]models.CreateAssistant, error) {
	collection := o.db.Collection(openaiCollection)

	filter := bson.M{"user_id": ID}

	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println("error closing cursor: ", err)
		}
	}(cursor, ctx)

	var assistants []models.CreateAssistant
	for cursor.Next(ctx) {
		var assistant models.CreateAssistant
		if err := cursor.Decode(&assistant); err != nil {
			return nil, err
		}

		assistants = append(assistants, assistant)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return assistants, nil
}

func (o *OpenaiMongo) GetAssistant(ctx context.Context, assistantID string) (*models.CreateAssistant, error) {
	collection := o.db.Collection(openaiCollection)

	objectID, err := primitive.ObjectIDFromHex(assistantID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var assistant models.CreateAssistant
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&assistant)
	if err != nil {
		return nil, err
	}

	return &assistant, nil
}

func (o *OpenaiMongo) Delete(ctx context.Context, ID string) error {
	collection := o.db.Collection(openaiCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
