package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Conversations struct {
	db *mongo.Database
}

func NewConversationsRepository(db *mongo.Client) *Conversations {
	return &Conversations{
		db: db.Database(dbName),
	}
}

// Create inserts a new conversation into the database
func (c *Conversations) Insert(ctx context.Context, conversation models.Conversation) (string, error) {

	collection := c.db.Collection(conversationsCollection)

	conversation.CreatedAt = time.Now().Add(-3 * time.Hour)

	// insert
	result, err := collection.InsertOne(ctx, conversation)
	if err != nil {
		return "", err
	}

	// id generated
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("falha ao obter o ID da conversa")
	}

	return insertedID.Hex(), nil
}

// Edit - Edits Conversations values
func (c *Conversations) Edit(ctx context.Context, ID string, newConversation models.Conversation) error {
	collection := c.db.Collection(conversationsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	// Converts the novaConversa struct to a map to use in $set
	updateFields := bson.M{}
	bsonBytes, err := bson.Marshal(newConversation)
	if err != nil {
		return err
	}
	bson.Unmarshal(bsonBytes, &updateFields)

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Search conversations by workflow id or customer id
func (c *Conversations) Find(ctx context.Context, query string) ([]models.Conversation, error) {
	collection := c.db.Collection(conversationsCollection)

	queryFields, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	// Construct a filter based on the fields of the query
	filter := bson.M{}
	for key, values := range queryFields {
		for _, value := range values {
			filter[key] = bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}
		}
	}

	// Perform the query in the database
	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	for cursor.Next(ctx) {
		var conversation models.Conversation
		if err := cursor.Decode(&conversation); err != nil {
			return nil, err
		}

		conversations = append(conversations, conversation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// Search id brings up a specific conversation by id
func (c *Conversations) FindId(ctx context.Context, ID string) (*models.Conversation, error) {
	collection := c.db.Collection(conversationsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	//Id filter
	filter := bson.M{"_id": objectID}

	// Perform the query in the database
	var conversations models.Conversation
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&conversations)
	if err != nil {
		return nil, err
	}

	return &conversations, nil
}

// Delete - Delete the conversation with the id sent
func (c *Conversations) Delete(ctx context.Context, ID string) error {
	collection := c.db.Collection(conversationsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// Id filter
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// Insert messages into the 'messages' field
func (c *Conversations) InsertMessage(ctx context.Context, idConversa string, newMessage models.Message) error {
	collection := c.db.Collection(conversationsCollection)

	//set timestampo now
	newMessage.Timestamp = time.Now().Unix()

	// Converte a string para o object id do mongo
	objectID, erro := primitive.ObjectIDFromHex(idConversa)
	if erro != nil {
		return erro
	}

	// Define the update operation to add the new message to the array of messages.
	operacao := bson.M{
		"$push": bson.M{"mensagens": newMessage},
	}

	// update conversation document
	_, err := collection.UpdateByID(ctx, objectID, operacao)
	if err != nil {
		return err
	}

	return nil
}
