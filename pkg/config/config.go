package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Config struct {
	Database *DatabaseConfig `json:"database" yaml:"database"`

	SecretKey   string `yaml:"secret_key" json:"secret_key"`
	UI          string `yaml:"ui" json:"ui"`
	RoutePrefix string `json:"route_prefix" yaml:"route_prefix"`
}

type DatabaseConfig struct {
	URI  string `json:"uri" yaml:"uri"`
	Name string `json:"name" yaml:"name"`
}

func (d *DatabaseConfig) Connect(ctx context.Context) (*mongo.Database, error) {
	wc := writeconcern.New(writeconcern.J(true))
	opts := options.Client().ApplyURI(d.URI).SetWriteConcern(wc)

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return client.Database(d.Name), nil
}
