package postgres

import (
	"context"
	"errors"
	"github.com/VanceMichael/go-base-airbridge/internal/auth"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"github.com/VanceMichael/go-base-airbridge/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct{ db *pgxpool.Pool }

func NewStore(db *DB) *Store { return &Store{db: db.Pool} }

var _ repository.Store = (*Store)(nil)

func (s *Store) Ping(ctx context.Context) error { return s.db.Ping(ctx) }
func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.user(ctx, `SELECT id,tenant_id,email,password_hash,role,active,created_at,deactivated_at FROM users WHERE email=$1`, email)
}
func (s *Store) GetUser(ctx context.Context, id string) (domain.User, error) {
	return s.user(ctx, `SELECT id,tenant_id,email,password_hash,role,active,created_at,deactivated_at FROM users WHERE id=$1`, id)
}
func (s *Store) user(ctx context.Context, q string, args ...any) (domain.User, error) {
	var u domain.User
	var role string
	var hash string
	err := s.db.QueryRow(ctx, q, args...).Scan(&u.ID, &u.TenantID, &u.Email, &hash, &role, &u.Active, &u.CreatedAt, &u.DeactivatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, domain.ErrNotFound
	}
	if err != nil {
		return u, err
	}
	u.PasswordHash = []byte(hash)
	u.Role = domain.Role(role)
	return u, nil
}
func (s *Store) CreateUser(ctx context.Context, u domain.User) error {
	_, err := s.db.Exec(ctx, `INSERT INTO users(id,tenant_id,email,password_hash,role,active,created_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, u.ID, u.TenantID, u.Email, string(u.PasswordHash), u.Role, u.Active, u.CreatedAt)
	if err != nil {
		if isUnique(err) {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}
func (s *Store) DeactivateUser(ctx context.Context, id string, at time.Time) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, `UPDATE users SET active=false,deactivated_at=$2 WHERE id=$1`, id, at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	if _, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=$2 WHERE user_id=$1 AND revoked_at IS NULL`, id, at); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (s *Store) CreateSession(ctx context.Context, v domain.Session) error {
	_, err := s.db.Exec(ctx, `INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES($1,$2,$3,$4,$5)`, v.ID, v.UserID, v.TokenHash, v.ExpiresAt, v.CreatedAt)
	return err
}
func (s *Store) GetSession(ctx context.Context, token string) (domain.Session, error) {
	var v domain.Session
	var hash string
	err := s.db.QueryRow(ctx, `SELECT id,user_id,token_hash,expires_at,revoked_at,created_at FROM sessions WHERE token_hash=$1`, auth.HashToken(token)).Scan(&v.ID, &v.UserID, &hash, &v.ExpiresAt, &v.RevokedAt, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	v.TokenHash = hash
	return v, err
}
func (s *Store) RevokeSession(ctx context.Context, token string, at time.Time) error {
	tag, err := s.db.Exec(ctx, `UPDATE sessions SET revoked_at=$2 WHERE token_hash=$1 AND revoked_at IS NULL`, auth.HashToken(token), at)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (s *Store) RevokeUserSessions(ctx context.Context, id string, at time.Time) error {
	_, err := s.db.Exec(ctx, `UPDATE sessions SET revoked_at=$2 WHERE user_id=$1 AND revoked_at IS NULL`, id, at)
	return err
}
func isUnique(err error) bool {
	return err != nil && len(err.Error()) > 0 && contains(err.Error(), "duplicate key")
}
func contains(v, w string) bool {
	for i := 0; i+len(w) <= len(v); i++ {
		if v[i:i+len(w)] == w {
			return true
		}
	}
	return false
}
