package api

import (
	"errors"
	"net/http"

	"github.com/go-api-template/pkg/engine"
	"github.com/monzo/typhon"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Router) GetUsers(req typhon.Request) typhon.Response {
	users, err := r.Engine.GetUsers(req.Context)
	if err != nil {
		return req.ResponseWithCode(err, http.StatusInternalServerError)
	}

	return req.Response(users)
}

func (r *Router) GetUser(req typhon.Request) typhon.Response {
	// Get uuid from address params
	uuid := r.Params(req)["uuid"]

	user, err := r.Engine.GetUser(req.Context, uuid)
	if err != nil {
		switch err {
		case engine.ErrNotFound:
			return req.ResponseWithCode(err, http.StatusNotFound)
		default:
			return req.ResponseWithCode(err, http.StatusInternalServerError)
		}
	}

	return req.Response(user)
}

func (r *Router) CreateUser(req typhon.Request) typhon.Response {
	// Get user from req body
	user := &engine.User{}
	if err := req.Decode(user); err != nil {
		return req.ResponseWithCode(err, http.StatusBadRequest)
	}

	if err := r.Engine.CreateUser(req.Context, user); err != nil {
		switch err {
		case engine.ErrDocumentAlreadyExists:
			return req.ResponseWithCode(req, http.StatusConflict)
		case engine.ErrTimeout:
			return req.ResponseWithCode(req, http.StatusFailedDependency)
		default:
			return req.ResponseWithCode(req, http.StatusInternalServerError)
		}
	}

	return req.ResponseWithCode(req, http.StatusCreated)
}

func (r *Router) DeleteUser(req typhon.Request) typhon.Response {
	// Get uuid from address params
	uuid := r.Params(req)["uuid"]

	if err := r.Engine.DeleteUser(req.Context, uuid); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return req.ResponseWithCode(err.Error(), http.StatusNotFound)
		}
		return req.ResponseWithCode(err.Error(), http.StatusInternalServerError)
	}

	return req.ResponseWithCode(nil, http.StatusNoContent)
}
