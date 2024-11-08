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

type Workflows struct {
	db *mongo.Database
}

func NewWorkflowsRepository(db *mongo.Client) *Workflows {
	return &Workflows{
		db: db.Database(dbName),
	}
}

// Inserts a new Workflow into the 'workflows' collection, setting the creation timestamp, and returns the generated ID.
func (r *Workflows) Insert(ctx context.Context, workflow models.Workflow) (string, error) {
	collection := r.db.Collection(workflowCollection)
	workflow.CreatedAt = time.Now().Add(-3 * time.Hour)
	workflow.UpdateAt = time.Now().Add(-3 * time.Hour)

	// Insert the workflow in colletion
	result, err := collection.InsertOne(context.Background(), workflow)
	if err != nil {
		return "", err
	}

	// Return id generate
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("falha ao obter o ID do Workflow inserido")
	}

	return insertedID.Hex(), nil
}

// Edit updates specific workflow information in the "workflows" collection based on the provided ID.
func (r *Workflows) Edit(ctx context.Context, ID string, newWorkflow models.Workflow) error {
	collection := r.db.Collection(workflowCollection)

	newWorkflow.UpdateAt = time.Now().Add(-3 * time.Hour)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// Id filter
	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": newWorkflow, // Update all fields
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Find retrieves workflows records from the database based on the provided query string.
func (r *Workflows) Find(ctx context.Context, query string) ([]models.Workflow, error) {
	collection := r.db.Collection(workflowCollection)

	// Parse the query string into individual fields
	queryFields, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

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

	var workflows []models.Workflow
	for cursor.Next(ctx) {
		var workflow models.Workflow
		if err := cursor.Decode(&workflow); err != nil {
			return nil, err
		}

		workflows = append(workflows, workflow)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return workflows, nil
}

// FindId fetches a specific workflows by their ID.
func (r *Workflows) FindId(ctx context.Context, ID string) (*models.Workflow, error) {
	collection := r.db.Collection(workflowCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	// Id filter
	filter := bson.M{"_id": objectID}

	// Perform the query in the database
	var workflow models.Workflow
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

// Delete - Deletes the workflows account with the provided ID.
func (r *Workflows) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(workflowCollection)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// Id Filter
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
