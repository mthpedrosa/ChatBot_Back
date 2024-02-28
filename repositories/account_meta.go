package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Metas struct {
	db *mongo.Database
}

func NewMetaRepository(db *mongo.Client) *Metas {
	return &Metas{db: db.Database(dbName)}
}

// Insert the Meta Account in the database
func (r *Metas) Insert(ctx context.Context, meta models.Meta) (string, error) {
	collection := r.db.Collection(metaAccountsCollection)

	meta.CreatedAt = time.Now().Add(-3 * time.Hour)
	meta.Webhook = WebhookGenerate(10)

	// Check if the user's email already exists in the database.
	existingContaMeta := &models.Meta{}
	filter := bson.M{"token": meta.Token}
	err := collection.FindOne(ctx, filter).Decode(existingContaMeta)

	if err == nil {
		return "", errors.New("já exsite uma conta cadastrado com esse token de conexão")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	// Insert the account into the collection.
	result, err := collection.InsertOne(ctx, meta)
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

// Edit the Meta  Account values
func (r *Metas) Edit(ctx context.Context, ID string, newMeta models.Meta) error {
	collection := r.db.Collection(metaAccountsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// filter id
	filter := bson.M{"_id": objectID}

	// Convert the newMeta struct to a map for use in the $set operation.
	updateFields := bson.M{}
	bsonBytes, err := bson.Marshal(newMeta)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bsonBytes, &updateFields)
	if err != nil {
		return err
	}

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Find search for the account with the query sent
func (r *Metas) Find(ctx context.Context, query string) ([]models.Meta, error) {
	collection := r.db.Collection(metaAccountsCollection)

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

	var metas []models.Meta
	for cursor.Next(ctx) {
		var meta models.Meta
		if err := cursor.Decode(&meta); err != nil {
			return nil, err
		}

		metas = append(metas, meta)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return metas, nil
}

// FindId search for the account with the id sent
func (r *Metas) FindId(ctx context.Context, ID string) (*models.Meta, error) {
	collection := r.db.Collection(metaAccountsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var meta models.Meta
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

// Delete the meta account with the id sent
func (r *Metas) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(metaAccountsCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// Search for the meta account based on the provided phone ID.
func (r *Metas) FindPhoneID(ctx context.Context, ID string) (*models.Meta, error) {
	collection := r.db.Collection(metaAccountsCollection)

	filter := bson.M{
		"phones_meta": bson.M{
			"$elemMatch": bson.M{"id": ID},
		},
	}

	var meta models.Meta
	err := collection.FindOne(ctx, filter).Decode(&meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func WebhookGenerate(tamanho int) string {
	var letras = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, tamanho)
	for i := range b {
		b[i] = letras[rand.Intn(len(letras))]
	}
	return string(b)
}
