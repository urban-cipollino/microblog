package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"micro-blog/internal/database"
	"micro-blog/internal/domain"
)

type PostLikeRepository struct {
	db database.DBTX
	qb sq.StatementBuilderType
}

func NewPostLikeRepository(db *pgxpool.Pool, qb sq.StatementBuilderType) *PostLikeRepository {
	return &PostLikeRepository{
		db: db,
		qb: qb,
	}
}

func (r *PostLikeRepository) WithTx(tx pgx.Tx) *PostLikeRepository {
	return &PostLikeRepository{db: tx, qb: r.qb}
}

func (r *PostLikeRepository) Create(ctx context.Context, like *domain.PostLike) (*domain.PostLike, error) {
	query, args, err := r.qb.
		Insert("post_likes").
		Columns("post_id", "post_created_at", "user_id").
		Values(like.PostID, like.PostCreatedAt, like.UserID).
		Suffix("ON CONFLICT DO NOTHING RETURNING created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&like.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return like, nil
}

func (r *PostLikeRepository) Exists(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	query, args, err := r.qb.
		Select("1").
		From("post_likes").
		Where(sq.Eq{
			"post_id": postID,
			"user_id": userID,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = r.db.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *PostLikeRepository) CountByPost(ctx context.Context, postID uuid.UUID, postCreatedAt any) (int, error) {
	query, args, err := r.qb.
		Select("COUNT(*)").
		From("post_likes").
		Where(sq.Eq{
			"post_id":         postID,
			"post_created_at": postCreatedAt,
		}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	if err := r.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
