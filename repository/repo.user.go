package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/androsyz/nexus-user-svc/model"
	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepoUser struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewUserRepository(db *sqlx.DB, rdb *redis.Client) *RepoUser {
	return &RepoUser{
		db:  db,
		rdb: rdb,
	}
}

func (r *RepoUser) Create(ctx context.Context, user *model.UserDB) (*string, error) {
	user.ID = uuid.New()
	now := time.Now()

	query := `
		INSERT INTO users (id, email, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.Password, now, now)
	if err != nil {
		return nil, err
	}

	idStr := user.ID.String()

	return &idStr, nil
}

func (r *RepoUser) GetByID(ctx context.Context, id string) (*model.UserDB, error) {
	var user model.UserDB
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}

func (r *RepoUser) GetByEmail(ctx context.Context, email string) (*model.UserDB, error) {
	var user model.UserDB
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}
