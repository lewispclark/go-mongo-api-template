package engine

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type Engine struct {
	Config *Config

	UsersCollection *mongo.Collection
}

var (
	ErrNotFound              = mongo.ErrNoDocuments
	ErrInternalFault         = errors.New("internal fault")
	ErrDocumentAlreadyExists = errors.New("document already exists")
	ErrTimeout               = errors.New("timeout")
)

func IsDuplicateKey(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

type Config struct {
	Database *mongo.Database
}

func New(cfg *Config) (*Engine, error) {
	return &Engine{
		UsersCollection: cfg.Database.Collection("Users"),
		Config:          cfg,
	}, nil
}
