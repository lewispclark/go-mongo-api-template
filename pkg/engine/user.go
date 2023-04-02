package engine

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UUID      string `json:"uuid" bson:"_id" validate:"required,email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
}

func (e *Engine) GetUsers(ctx context.Context) ([]User, error) {
	// Filter for Find query
	filter := bson.D{}

	cur, err := e.UsersCollection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute mongo aggregation on users collection")
		return nil, err
	}

	// Users array to store the decoded documents
	users := []User{}

	// Convert results to User structs
	if err = cur.All(ctx, &users); err != nil {
		log.Error().Err(err).Msg("failed to decode cursor into 'User' struct")
		return nil, err
	}

	return users, nil
}

func (e *Engine) GetUser(ctx context.Context, uuid string) (*User, error) {
	// Filter for Find query
	filter := bson.D{{Key: "_id", Value: uuid}}

	// User struct to store decoded document
	user := &User{}
	if err := e.UsersCollection.FindOne(ctx, filter).Decode(user); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			log.Error().Str("uuid", uuid).Err(err).Msg("user not found")
			return nil, ErrNotFound
		default:
			log.Error().Str("uuid", uuid).Err(err).Msg("failed to execute mongo find query on users collection")
			return nil, fmt.Errorf("%w: %s", ErrInternalFault, err.Error())
		}
	}

	return user, nil
}

func (e *Engine) CreateUser(ctx context.Context, user *User) error {
	if _, err := e.UsersCollection.InsertOne(ctx, *user); err != nil {
		switch {
		case mongo.IsDuplicateKeyError(err):
			log.Error().Str("uuid", user.UUID).Err(err).Msg("user already exists")
			return ErrDocumentAlreadyExists
		case mongo.IsTimeout(err):
			log.Error().Err(err).Msg("mongo connection timed out")
			return ErrTimeout
		default:
			log.Error().Err(err).Msg("failed to insert document into users collection")
			return err
		}
	}

	return nil
}

func (e *Engine) DeleteUser(ctx context.Context, uuid string) error {
	// Filter for Delete query
	filter := bson.D{{Key: "_id", Value: uuid}}

	if _, err := e.UsersCollection.DeleteOne(ctx, filter); err != nil {
		log.Error().Str("uuid", uuid).Err(err).Msg("failed to execute mongo delete on users collection")
		return err
	}

	return nil
}
