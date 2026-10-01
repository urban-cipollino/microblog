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
	"micro-blog/internal/dto"
	"micro-blog/internal/exception"
)

type PostRepository struct {
	db database.DBTX
	qb sq.StatementBuilderType
}

func NewPostRepository(db *pgxpool.Pool, qb sq.StatementBuilderType) *PostRepository {
	return &PostRepository{
		db: db,
		qb: qb,
	}
}

func (r *PostRepository) WithTx(tx pgx.Tx) *PostRepository {
	return &PostRepository{db: tx, qb: r.qb}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	query, args, err := r.qb.
		Insert("posts").
		Columns("id", "author_id", "text").
		Values(post.ID, post.AuthorID, post.Text).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, query, args...).Scan(&post.CreatedAt); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	query, args, err := r.qb.
		Update("posts").
		Set("text", post.Text).
		Where(sq.Eq{"id": post.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() == 0 {
		return nil, &exception.NotFoundError{Object: "Post"}
	}

	return post, nil
}

func (r *PostRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	query, args, err := r.qb.
		Select("id", "author_id", "text", "created_at").
		From("posts").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	p := &domain.Post{}
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&p.ID,
		&p.AuthorID,
		&p.Text,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &exception.NotFoundError{Object: "Post"}
		}
		return nil, err
	}

	return p, nil
}

func (r *PostRepository) FindAll(ctx context.Context) ([]*dto.PostDTO, error) {
	query, args, err := r.qb.
		Select(
			"p.id",
			"p.text",
			"p.created_at",
			"u.id AS author_id",
			"u.name AS author_name",
			"COUNT(pl.user_id) AS likes",
		).
		From("posts AS p").
		Join("users AS u ON p.author_id = u.id").
		LeftJoin("post_likes AS pl ON p.id = pl.post_id AND p.created_at = pl.post_created_at").
		GroupBy("p.id", "p.text", "p.created_at", "u.id", "u.name").
		OrderBy("p.created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*dto.PostDTO, 0)
	for rows.Next() {
		post := &domain.Post{}
		author := &domain.User{}
		var likes int

		if err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.CreatedAt,
			&author.ID,
			&author.Name,
			&likes,
		); err != nil {
			return nil, err
		}

		post.AuthorID = author.ID

		posts = append(posts, &dto.PostDTO{
			Post:   post,
			Author: author,
			Likes:  likes,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
