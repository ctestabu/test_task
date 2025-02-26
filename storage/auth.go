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

// GetUser fetches a user by login
func (pg *Postgres) GetUser(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := pg.db.QueryRow(ctx, query, login)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// ValidateUser проверяет, существует ли пользователь с данным логином и паролем
func (pg *Postgres) ValidateUser(ctx context.Context, login, password string) (int64, error) {
	var user models.User

	// Запрос для получения данных пользователя по логину
	row := pg.db.QueryRow(ctx, "SELECT id, password_hash FROM users WHERE login = $1", login)

	// Сканируем результат запроса в структуру User
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Пользователь с таким логином не найден
			return 0, fmt.Errorf("invalid login/password")
		}
		return 0, fmt.Errorf("unable to fetch user data: %w", err)
	}

	// Хэшируем введённый пароль в MD5
	hashInput := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	// Сравнение хэшированного пароля с паролем в базе
	if user.PasswordHash != hashInput {
		// Если пароли не совпадают
		return 0, fmt.Errorf("invalid login/password")
	}

	// Если все прошло успешно, возвращаем ID пользователя
	return user.ID, nil
}

func (pg *Postgres) CreateSession(ctx context.Context, userID int64, ip string) (string, error) {
	var sessionID string

	// Создаем новую сессию с IP-адресом
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
		// Сессия истекла — удаляем её
		_ = pg.DeleteUserSession(ctx, userID)
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func (pg *Postgres) DeleteUserSession(ctx context.Context, userID int64) error {
	_, err := pg.db.Exec(ctx, `DELETE FROM sessions WHERE uid = $1`, userID)
	return err
}
