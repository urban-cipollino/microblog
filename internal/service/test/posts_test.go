package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosts(t *testing.T) {
	svc := setup(t)
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		posts, err := svc.posts.Posts(ctx)
		require.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("returns posts in desc order", func(t *testing.T) {
		user, err := svc.createUser.CreateUser(ctx, "alice", "password123")
		require.NoError(t, err)

		first, err := svc.createPost.CreatePost(ctx, user.ID, "first post")
		require.NoError(t, err)

		second, err := svc.createPost.CreatePost(ctx, user.ID, "second post")
		require.NoError(t, err)

		posts, err := svc.posts.Posts(ctx)
		require.NoError(t, err)

		require.Len(t, posts, 2)
		assert.Equal(t, second.Post.ID, posts[0].Post.ID)
		assert.Equal(t, first.Post.ID, posts[1].Post.ID)
		assert.Equal(t, "alice", posts[0].Author.Name)
	})
}
