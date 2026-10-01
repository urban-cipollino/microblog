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
	"micro-blog/internal/exception"
)

type UserRepository struct {
	db database.DBTX
	qb sq.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool, qb sq.StatementBuilderType) *UserRepository {
	return &UserRepository{
		db: db,
		qb: qb,
	}
}

func (r *UserRepository) WithTx(tx pgx.Tx) *UserRepository {
	return &UserRepository{db: tx, qb: r.qb}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query, args, err := r.qb.
		Insert("users").
		Columns("id", "name", "password").
		Values(user.ID, user.Name, user.Password).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	createdUser := &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&createdUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (r *UserRepository) IsNameTaken(ctx context.Context, name string) bool {
	query, args, err := r.qb.
		Select("COUNT(1) > 0").
		From("users").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return false
	}

	var exists bool

	err = r.db.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query, args, err := r.qb.
		Select("id", "name", "password", "created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	u := &domain.User{}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &exception.NotFoundError{Object: "User"}
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (*domain.User, error) {
	query, args, err := r.qb.
		Select("id", "name", "password", "created_at").
		From("users").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	u := &domain.User{}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &exception.NotFoundError{Object: "User"}
		}

		return nil, err
	}

	return u, nil
}
