package auth

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewJWTService,
	NewPasswordService,
	NewLoginService,
)
