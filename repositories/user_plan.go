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

// UserPlanRepository handles CRUD operations for UserPlan.
type UserPlanRepository struct {
	db *mongo.Database
}

// NewUserPlanRepository initializes a new repository for UserPlan.
func NewUserPlanRepository(db *mongo.Client) *UserPlanRepository {
	return &UserPlanRepository{
		db: db.Database(dbName),
	}
}

func (r *UserPlanRepository) Insert(ctx context.Context, userPlan models.UserPlan) (string, error) {
	collection := r.db.Collection(userPlanCollection)

	userPlan.CreatedAt = time.Now().Add(-3 * time.Hour)

	// Check if the user's email already exists in the database.
	existingContaMeta := &models.UserPlan{}
	filter := bson.M{"user_id": userPlan.UserID}
	err := collection.FindOne(ctx, filter).Decode(existingContaMeta)

	if err == nil {
		return "", errors.New("já exsite uma conta cadastrado com esse token de conexão")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	// Insert the account into the collection.
	result, err := collection.InsertOne(ctx, userPlan)
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

// Find retrieves all UserPlans that match a given filter.
func (r *UserPlanRepository) Find(ctx context.Context, query string) ([]models.UserPlan, error) {
	collection := r.db.Collection(userPlanCollection)

	queryFields, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	// Build a filter based on the query fields
	filter := bson.M{}
	for key, values := range queryFields {
		for _, value := range values {
			filter[key] = bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}
		}
	}

	// Execute the query on the database
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var usersPlans []models.UserPlan
	for cursor.Next(ctx) {
		var userPlan models.UserPlan
		if err := cursor.Decode(&userPlan); err != nil {
			return nil, err
		}

		usersPlans = append(usersPlans, userPlan)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return usersPlans, nil
}

// FindByID retrieves a UserPlan by its ID.
func (r *UserPlanRepository) FindId(ctx context.Context, ID string) (*models.UserPlan, error) {
	collection := r.db.Collection(userPlanCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	//filter the id
	filter := bson.M{"_id": objectID}

	// Perform the query in the database
	var userPlan models.UserPlan
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&userPlan)
	if err != nil {
		return nil, err
	}

	return &userPlan, nil
}

// Update updates an existing UserPlan with new values.
func (r *UserPlanRepository) Edit(ctx context.Context, ID string, newUserPlan models.UserPlan) error {
	collection := r.db.Collection(userPlanCollection)

	newUserPlan.UpdatedAt = time.Now().Add(-3 * time.Hour)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	// Convert the newUserPlan struct to a map for use in the $set operation.
	updateFields := bson.M{}
	bsonBytes, err := bson.Marshal(newUserPlan)
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

// Delete removes a UserPlan by its ID.
func (r *UserPlanRepository) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(userPlanCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
