package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"micro-blog/internal/exception"
)

func TestCreateUser(t *testing.T) {
	svc := setup(t)
	ctx := context.Background()

	t.Run("creates user successfully", func(t *testing.T) {
		user, err := svc.createUser.CreateUser(ctx, "alice", "password123")
		require.NoError(t, err)

		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "alice", user.Name)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, "password123", user.Password)
	})

	t.Run("rejects duplicate name", func(t *testing.T) {
		_, err := svc.createUser.CreateUser(ctx, "bob", "password123")
		require.NoError(t, err)

		_, err = svc.createUser.CreateUser(ctx, "bob", "other_password")
		require.Error(t, err)

		var conflict *exception.ConflictError
		assert.ErrorAs(t, err, &conflict)
	})
}
