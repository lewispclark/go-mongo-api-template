package api

import (
	"github.com/go-api-template/pkg/engine"
	"github.com/monzo/typhon"
)

type Config struct {
	Engine *engine.Engine
}

type Router struct {
	*typhon.Router
	*Config
}

func New(cfg *Config) *Router {
	return &Router{
		Router: &typhon.Router{},
		Config: cfg,
	}
}

func (r *Router) Serve() typhon.Service {
	r.GET("/users", typhon.Service(r.GetUsers))
	r.PUT("/users", typhon.Service(r.CreateUser))
	r.GET("/users/:uuid", typhon.Service(r.GetUser))
	r.DELETE("/users/:uuid", typhon.Service(r.DeleteUser))

	return r.Router.Serve()
}
