package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	svc := setup(t)
	ctx := context.Background()

	t.Run("creates post successfully", func(t *testing.T) {
		user, err := svc.createUser.CreateUser(ctx, "alice", "password123")
		require.NoError(t, err)

		post, err := svc.createPost.CreatePost(ctx, user.ID, "hello world")
		require.NoError(t, err)

		assert.NotEmpty(t, post.Post.ID)
		assert.Equal(t, "hello world", post.Post.Text)
		assert.Equal(t, user.ID, post.Author.ID)
		assert.Equal(t, "alice", post.Author.Name)
		assert.Equal(t, 0, post.Likes)
	})

	t.Run("fails with non-existent author", func(t *testing.T) {
		_, err := svc.createPost.CreatePost(ctx, [16]byte{1}, "hello")
		require.Error(t, err)
	})
}
