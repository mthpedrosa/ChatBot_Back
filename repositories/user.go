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

type Users struct {
	db *mongo.Database
}

func NewUsersRepository(db *mongo.Client) *Users {

	return &Users{
		db: db.Database(dbName),
	}
}

// Insert adds a new user to the "users" collection, setting the creation timestamp.
// Checks if the user's email already exists to prevent duplicates. Returns the inserted ID on success.
// Returns an error if the email already exists or if there's an issue during insertion.
func (r *Users) Insert(ctx context.Context, user models.User) (string, error) {
	collection := r.db.Collection(userCollection)

	// define date of now
	user.CreatedAt = time.Now()

	// check if the user's email already exists in the database
	existingUser := &models.User{}
	filter := bson.M{"email": user.Email}
	err := collection.FindOne(ctx, filter).Decode(existingUser)

	if err == nil {
		// The WhatsApp ID already exists in the database, return an error or an appropriate response
		return "", errors.New("a user with this Email is already registered")
	} else if err != mongo.ErrNoDocuments {
		// There was an error while querying the database, return the error
		return "", err
	}

	// Insert user in the collection
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// The ID generated for the entered user
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("falha ao obter o ID do usuário inserido")
	}

	return insertedID.Hex(), nil
}

// Edit updates a user's name in the "users" collection based on the provided ID.
// Performs an update operation to change the user's name field.
// Returns an error if there's an issue during the update process.
func (r *Users) Edit(ctx context.Context, ID string, newUser models.User) error {
	collection := r.db.Collection(userCollection)

	newUser.UpdateAt = time.Now().Add(-3 * time.Hour)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// filter the id
	filter := bson.M{"_id": objectID}

	update := bson.M{"$set": bson.M{"name": newUser.Name, "update_at": newUser.UpdateAt}} // editable fields

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Find searches for users in the "users" collection using the provided query string.
// It looks for users whose name or email matches the query's filter criteria.
// Receives a query string with filters for name or email fields.
// Returns a list of models.User for users matching the filter criteria.
// Returns an error if any issue arises during the search operation.
func (r *Users) Find(ctx context.Context, query string) ([]models.User, error) {
	collection := r.db.Collection(userCollection)

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

	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindId fetches a specific user by their ID.
// It retrieves the user details based on the provided ID from the database.
// Returns a pointer to models.user representing the found user or an error if not found.
func (r *Users) FindId(ctx context.Context, ID string) (*models.User, error) {
	collection := r.db.Collection(userCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	//filter the id
	filter := bson.M{"_id": objectID}

	// Perform the query in the database
	var user models.User
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete - Deletes the user account with the provided ID.
func (r *Users) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(userCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// Search using the email
func (r *Users) FindbyEmail(ctx context.Context, email string) (models.User, error) {
	collection := r.db.Collection(userCollection)

	filter := bson.M{"email": email}

	// Define the projections to select the desired fields
	projection := bson.M{
		"name":     1,
		"email":    1,
		"password": 1,
		"_id":      1,
		"profile":  1,
	}

	existingUser := &models.User{}
	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(existingUser)
	if err != nil {
		return models.User{}, err
	}

	return *existingUser, nil
}
