package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"fmt"

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
func (o *OpenaiMongo) Insert(ctx context.Context, dto models.Assistant) (string, error) {
	collection := o.db.Collection(openaiCollection)

	// Check if the user's email already exists in the database.
	existingAssistant := &models.Assistant{}
	filter := bson.M{"id": dto.ID}
	err := collection.FindOne(ctx, filter).Decode(existingAssistant)

	if err == nil {
		return "", errors.New("já exsite um assistante cadastrado com esse id")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	// Insert the account into the collection.
	result, err := collection.InsertOne(ctx, dto)
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

func (o *OpenaiMongo) FindAllUser(ctx context.Context, ID string) ([]models.Assistant, error) {
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

	var assistants []models.Assistant
	for cursor.Next(ctx) {
		var assistant models.Assistant
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
