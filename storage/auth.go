package storage

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/ctestabu/test_task/models"
	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) GetUser(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := pg.db.QueryRow(ctx, query, login)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (pg *Postgres) ValidateUser(ctx context.Context, login, password string) (int64, error) {
	var user models.User

	row := pg.db.QueryRow(ctx, "SELECT id, password_hash FROM users WHERE login = $1", login)

	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("invalid login/password")
		}
		return 0, fmt.Errorf("unable to fetch user data: %w", err)
	}

	hashInput := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	if user.PasswordHash != hashInput {
		return 0, fmt.Errorf("invalid login/password")
	}

	return user.ID, nil
}

func (pg *Postgres) CreateSession(ctx context.Context, userID int64, ip string) (string, error) {
	var sessionID string

	err := pg.db.QueryRow(ctx, `
        INSERT INTO sessions (uid, ip_address, expires_at)
        VALUES ($1, $2, now() + interval '24 hours') 
        RETURNING id
    `, userID, ip).Scan(&sessionID)

	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

func (pg *Postgres) ValidateSession(ctx context.Context, token string) (int64, error) {
	var userID int64
	var expiresAt time.Time

	err := pg.db.QueryRow(ctx, `
        SELECT uid, expires_at FROM sessions WHERE id = $1
    `, token).Scan(&userID, &expiresAt)

	if err != nil {
		return 0, errors.New("invalid token")
	}

	if time.Now().After(expiresAt) {
		_ = pg.DeleteUserSession(ctx, userID)
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func (pg *Postgres) DeleteUserSession(ctx context.Context, userID int64) error {
	_, err := pg.db.Exec(ctx, `DELETE FROM sessions WHERE uid = $1`, userID)
	return err
}
