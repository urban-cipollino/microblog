package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLikePost(t *testing.T) {
	svc := setup(t)
	ctx := context.Background()

	t.Run("likes post", func(t *testing.T) {
		author, err := svc.createUser.CreateUser(ctx, "alice", "password123")
		require.NoError(t, err)

		liker, err := svc.createUser.CreateUser(ctx, "bob", "password123")
		require.NoError(t, err)

		post, err := svc.createPost.CreatePost(ctx, author.ID, "hello world")
		require.NoError(t, err)

		dto, err := svc.likePost.LikePost(ctx, post.Post.ID, liker.ID)
		require.NoError(t, err)

		require.Equal(t, post.Post.ID, dto.Post.ID)

		// Wait for async worker to process the like
		require.Eventually(t, func() bool {
			var count int
			err := svc.pool.QueryRow(ctx,
				"SELECT COUNT(*) FROM post_likes WHERE post_id = $1 AND user_id = $2",
				post.Post.ID, liker.ID,
			).Scan(&count)
			return err == nil && count == 1
		}, 2*time.Second, 50*time.Millisecond)
	})

	t.Run("idempotent like", func(t *testing.T) {
		author, err := svc.createUser.CreateUser(ctx, "charlie", "password123")
		require.NoError(t, err)

		liker, err := svc.createUser.CreateUser(ctx, "dave", "password123")
		require.NoError(t, err)

		post, err := svc.createPost.CreatePost(ctx, author.ID, "test post")
		require.NoError(t, err)

		_, err = svc.likePost.LikePost(ctx, post.Post.ID, liker.ID)
		require.NoError(t, err)

		_, err = svc.likePost.LikePost(ctx, post.Post.ID, liker.ID)
		require.NoError(t, err)

		// Wait for both jobs to process, then verify only 1 like exists
		require.Eventually(t, func() bool {
			var count int
			err := svc.pool.QueryRow(ctx,
				"SELECT COUNT(*) FROM post_likes WHERE post_id = $1 AND user_id = $2",
				post.Post.ID, liker.ID,
			).Scan(&count)
			return err == nil && count == 1
		}, 2*time.Second, 50*time.Millisecond)
	})
}
