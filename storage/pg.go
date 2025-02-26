package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
	initErr    error

	ErrInvalidSession = errors.New("invalid session")
	ErrNotFound       = errors.New("asset not found")
)

func NewPG(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			initErr = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}
		pgInstance = &Postgres{db}
	})

	// Return error if initialization failed
	if initErr != nil {
		return nil, initErr
	}

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}

func (pg *Postgres) StoreAsset(ctx context.Context, userID int64, name string, data []byte) error {
	query := `INSERT INTO assets (name, uid, data) VALUES ($1, $2, $3) 
              ON CONFLICT (name, uid) DO UPDATE SET data = EXCLUDED.data`
	_, err := pg.db.Exec(ctx, query, name, userID, data)
	return err
}

func (pg *Postgres) GetAsset(ctx context.Context, userID int64, name string) ([]byte, error) {
	var data []byte
	query := `SELECT data FROM assets WHERE name = $1 AND uid = $2`

	err := pg.db.QueryRow(ctx, query, name, userID).Scan(&data)
	if err != nil {
		return nil, ErrNotFound
	}

	return data, nil
}

func (pg *Postgres) DeleteAsset(ctx context.Context, userID int64, assetName string) error {
	var exists bool
	err := pg.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM assets WHERE name = $1 AND uid = $2)`, assetName, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check asset existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("asset not found")
	}

	result, err := pg.db.Exec(ctx, `DELETE FROM assets WHERE name = $1 AND uid = $2`, assetName, userID)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

func (pg *Postgres) ListAssets(ctx context.Context, userID int64) ([]string, error) {
	rows, err := pg.db.Query(ctx, "SELECT name FROM assets WHERE uid = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve assets: %w", err)
	}
	defer rows.Close()

	var assets []string

	for rows.Next() {
		var assetName string
		if err := rows.Scan(&assetName); err != nil {
			return nil, fmt.Errorf("unable to scan asset name: %w", err)
		}
		assets = append(assets, assetName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return assets, nil
}
