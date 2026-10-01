package user

import (
	"github.com/google/wire"

	"micro-blog/internal/handler"
)

var Set = wire.NewSet(
	NewCreateUserService,
	wire.Bind(new(handler.CreateUserService), new(*CreateUserService)),
)
