package repositories

import (
	"autflow_back/models"
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	db *mongo.Database
}

func NewSessionsRepository(db *mongo.Client) *Session {

	return &Session{
		db: db.Database(dbName),
	}
	//return &Session{db}
}

// Create inserts a new conversation into the database
func (r *Session) Insert(ctx context.Context, sessions models.Session) (models.Session, error) {

	collection := r.db.Collection(sessionsCollection)

	sessions.CreatedAt = time.Now().Add(-3 * time.Hour)
	sessions.UpdateAt = time.Now().Add(-3 * time.Hour)
	sessions.OtherFields = []models.Fields{}

	// insert
	result, err := collection.InsertOne(ctx, sessions)
	if err != nil {
		return models.Session{}, err
	}

	// id generated
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.Session{}, errors.New("Falha ao obter o ID da sessão")
	}
	sessions.ID = insertedID

	return sessions, nil
}

// Edit an existing session using the ID
func (r *Session) Edit(ctx context.Context, ID string, session models.Session) error {
	collection := r.db.Collection(sessionsCollection)

	session.UpdateAt = time.Now().Add(-3 * time.Hour)

	if session.Status == "finished" {
		session.FinishedAt = time.Now().Add(-3 * time.Hour)
	}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	// ID filter
	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"status":      session.Status,
			"finished_at": session.FinishedAt,
			"update_at":   session.UpdateAt,
		},
	} // Editable fields

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// This repository contains the implementation of the Session data access layer in Go. It's responsible for handling all interactions with the MongoDB database related to session data. The repository includes functionality to dynamically construct and execute queries based on various filter criteria, making it flexible for different use cases. It supports advanced filtering, including handling complex structures like `other_fields`.
func (r *Session) Find(ctx context.Context, queryString string) ([]models.Session, error) {
	collection := r.db.Collection(sessionsCollection)

	queryFields, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	filter := bson.M{}
	for key, values := range queryFields {
		for _, value := range values {
			if key == "other_fields" {
				otherFieldsFilters := parseOtherFields(value)
				for _, f := range otherFieldsFilters {
					filter[string(f.Key)] = f.Value
				}
			} else if key == "created_at" || key == "update_at" || key == "finished_at" {
				dateFilter, err := createDateFilter(key, value)
				if err != nil {
					fmt.Println("Erro ao criar filtro de data:", err)
					continue
				}
				filter[key] = dateFilter
			} else {
				// Para outros campos
				filter[key] = bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}
			}
		}
	}

	fmt.Println("MongoDB filter:", filter)

	// Execute the query on the database and handle the results.
	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	for cursor.Next(ctx) {
		var session models.Session
		if err := cursor.Decode(&session); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// FindId fetches a specific session by their ID.
func (r *Session) FindId(ctx context.Context, identifier string) (*models.Session, error) {
	collection := r.db.Collection(sessionsCollection)

	objectID, err := primitive.ObjectIDFromHex(identifier)
	if err != nil {
		return nil, err
	}

	//filter the id
	filter := bson.M{"_id": objectID}

	// Perform the query in the database
	var session models.Session
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Delete - Deletes the session account with the provided ID.
func (r *Session) Delete(ctx context.Context, ID string) error {
	collection := r.db.Collection(sessionsCollection)

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

// Insert messages into the 'messages' field
func (r *Session) InsertMessage(ctx context.Context, idSession string, newMessage models.Message) error {
	collection := r.db.Collection(sessionsCollection)

	//set timestampo now
	newMessage.Timestamp = time.Now().Unix()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string to the mongo object id
	objectID, erro := primitive.ObjectIDFromHex(idSession)
	if erro != nil {
		return erro
	}

	// Defines the update operation to add the new message to the message array
	operacao := bson.M{
		"$push": bson.M{"messages": newMessage},
	}

	// Update the conversation document in the database
	_, err := collection.UpdateByID(ctx, objectID, operacao)
	if err != nil {
		return err
	}

	return nil
}

// Update the last traversed node
func (r *Session) UpdateLastNode(ctx context.Context, idNode string, idSession string) error {
	fmt.Println("Alterando nó da session : " + idNode)
	fmt.Println("ID da session : " + idSession)
	collection := r.db.Collection(sessionsCollection)

	// Convert string for object id
	objectID, erro := primitive.ObjectIDFromHex(idSession)
	if erro != nil {
		return erro
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update to set the "last_node" field with the provided value
	update := bson.M{"$set": bson.M{"last_node": idNode}}

	// Perform the update of the conversation document in the database
	_, err := collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return err
	}

	return nil
}

// updateSessionField atualiza ou adiciona um field em other_fields[] na sessão especificada
func (r *Session) UpdateSessionField(ctx context.Context, sessionID string, field models.Fields) error {
	collection := r.db.Collection(sessionsCollection)

	// Convert the string to the MongoDB ObjectID
	objectID, erro := primitive.ObjectIDFromHex(sessionID)
	if erro != nil {
		return erro
	}

	// Checks if the field already exists and updates or adds as necessary
	filter := bson.M{"_id": objectID, "other_fields.name": field.Name}
	update := bson.M{"$set": bson.M{"other_fields.$": field}}
	result, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	// If no document has been updated, it means the field does not exist and needs to be added
	if result.MatchedCount == 0 {
		filter = bson.M{"_id": objectID}
		update = bson.M{"$push": bson.M{"other_fields": field}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}

// parseOtherFields analisa a string other_fields e retorna um filtro para o MongoDB.
func parseOtherFields(otherFieldsStr string) bson.D {
	var otherFieldsFilters bson.D
	fields := strings.Split(otherFieldsStr, ";")
	for _, field := range fields {
		keyValue := strings.SplitN(field, "=", 2)
		if len(keyValue) == 2 {
			// Cria um filtro que corresponde aos objetos dentro do array other_fields
			filter := bson.E{"other_fields", bson.M{"$elemMatch": bson.M{"name": keyValue[0], "value": keyValue[1]}}}
			otherFieldsFilters = append(otherFieldsFilters, filter)
		}
	}
	return otherFieldsFilters
}

func createDateFilter(field, value string) (bson.M, error) {
	// Verifica se o valor é um intervalo de datas
	if strings.Contains(value, ",") {
		dates := strings.Split(value, ",")
		if len(dates) != 2 {
			return nil, errors.New("formato de intervalo de datas inválido")
		}
		startDate, err := time.Parse(time.RFC3339, dates[0])
		if err != nil {
			return nil, err
		}
		endDate, err := time.Parse(time.RFC3339, dates[1])
		if err != nil {
			return nil, err
		}
		return bson.M{"$gte": startDate, "$lte": endDate}, nil
	} else {
		// Trata o valor como uma data única
		date, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, err
		}
		return bson.M{"$eq": date}, nil
	}
}
