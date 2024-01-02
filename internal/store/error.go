package store

import "errors"

var (
	EntityNotFound      = errors.New("entity not found")
	EntityAlreadyExists = errors.New("entity already exists")
)
