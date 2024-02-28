package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Customers struct {
	db *mongo.Database
}

func NewCustomersRepository(db *mongo.Client) *Customers {
	//return &Customers{db}

	return &Customers{
		db: db.Database(dbName),
	}
}

// Insert creates a new customer record in the "customers" collection.
// It takes a models.Customer object as input and sets the creation timestamp.
// Verifies if the user's WhatsApp ID already exists in the database and prevents duplicate entries.
// If the WhatsApp ID already exists, it returns an error indicating that the user is already registered.
// If an error occurs during the database query, it returns the respective error.
// Inserts the new customer into the database and returns the inserted ID in hexadecimal format on success.
// Returns an error if there's a failure while inserting the customer into the database.
func (r *Customers) Insert(ctx context.Context, customer models.Customer) (string, error) {
	collection := r.db.Collection(customerCollection)

	customer.CreatedAt = time.Now().Add(-3 * time.Hour)
	customer.UpdateAt = time.Now().Add(-3 * time.Hour)
	customer.OtherFields = []models.Fields{}

	// Check if the user's WhatsApp ID already exists in the database
	existingUser := &models.Customer{}
	filter := bson.M{"whatsapp_id": customer.WhatsAppID}
	err := collection.FindOne(ctx, filter).Decode(existingUser)
	if err == nil {
		// The WhatsApp ID already exists in the database, return an error or an appropriate response
		return "", errors.New("A user with this WhatsApp ID is already registered")
	} else if err != mongo.ErrNoDocuments {
		// There was an error while querying the database, return the error
		return "", err
	}

	// Insert the customer into the database
	result, err := collection.InsertOne(context.Background(), customer)
	if err != nil {
		return "", fmt.Errorf("failed to insert customer into the database: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to retrieve the inserted Customer ID")
	}

	return insertedID.Hex(), nil
}

// Edit updates specific customer information in the "customers" collection based on the provided ID.
// It takes the customer's ID and a new customer model (models.Customer) containing fields to be updated.
// Executes an update operation in the MongoDB database, modifying the "name", "email", and "phone" fields as provided.
// Returns an error, if any, while attempting to perform the update.
func (r *Customers) Edit(ctx context.Context, ID string, newCustomer models.Customer) error {
	collection := r.db.Collection(customerCollection)

	newCustomer.UpdateAt = time.Now().Add(-3 * time.Hour)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"name":  newCustomer.Name,
			"email": newCustomer.Email,
			"phone": newCustomer.Phone,
		},
	} // Editable fields

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Find retrieves customer records from the database based on the provided query string.
// The query string can contain multiple fields and values in the format "field1=value1&field2=value2".
// It constructs a flexible filter for MongoDB to search across various fields.
// Returns a slice of models.Customer containing matching records or an error if the query fails.
func (r *Customers) Find(ctx context.Context, query string) ([]models.Customer, error) {
	collection := r.db.Collection(customerCollection)

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

	var customers []models.Customer
	for cursor.Next(ctx) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

// FindId fetches a specific customer by their ID.
// It retrieves the customer details based on the provided ID from the database.
// Returns a pointer to models.Customer representing the found customer or an error if not found.
func (r *Customers) FindId(ctx context.Context, identifier string, byWhatsAppID bool) (*models.Customer, error) {
	collection := r.db.Collection(customerCollection)

	var filter bson.M
	var err error

	if byWhatsAppID {
		filter = bson.M{"whatsapp_id": identifier}
	} else {
		objectID, err := primitive.ObjectIDFromHex(identifier)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": objectID}
	}

	var customer models.Customer
	err = collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

// Delete - Deletes the customer account with the provided ID.
// It performs a deletion operation in the "customers" collection based on the given ID.
// Receives the ID of the customer to be deleted and removes the corresponding record from the database.
// Returns an error if there's any issue encountered during the deletion process.
func (r *Customers) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(customerCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// UpdateCustomerField atualiza ou adiciona um field em other_fields[] na sessão especificada
func (r *Customers) UpdateCustomerField(ctx context.Context, customerID string, field models.Fields) error {
	if len(customerID) != 24 {
		return fmt.Errorf("customerID '%s' is not a valid ObjectID", customerID)
	}

	collection := r.db.Collection(customerCollection)

	// Convert the string to the MongoDB ObjectID
	objectID, erro := primitive.ObjectIDFromHex(customerID)
	if erro != nil {
		return erro
	}

	// Adicione o campo diretamente ao array, sem fazer uma verificação inicial
	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"other_fields": field}}
	_, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}
