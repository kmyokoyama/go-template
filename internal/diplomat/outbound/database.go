package outbound

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
	"go.uber.org/fx"
)

type PgDatabase struct {
	Conn *pgxpool.Pool
}

func (db *PgDatabase) FindVersion() (models.Version, error) {
	var version string

	err := db.Conn.QueryRow(context.Background(), "select 'v2'").Scan(&version)
	if err != nil {
		return models.Version{}, err
	}

	return models.Version{Version: version}, nil
}

func (db *PgDatabase) CreateUser(user models.User) error {
	tx, err := db.Conn.Begin(context.Background())

	if err != nil {
		return err
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	stmt := `INSERT INTO users(uuid, name) VALUES ($1, $2)`
	_, err = tx.Exec(context.Background(), stmt, user.Id, user.Name)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (db *PgDatabase) FindUser(id uuid.UUID) (models.User, error) {
	var name string

	err := db.Conn.QueryRow(context.Background(), "SELECT name FROM users WHERE users.uuid = $1", id).Scan(&name)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Id: id, Name: name}, nil
}

func NewDatabase(lc fx.Lifecycle, logger *slog.Logger) components.Database {
	db := &PgDatabase{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Started on", "port", ":8080")
			
			dbHost := os.Getenv("DB_HOST")
			dbDatabase := os.Getenv("DB_DATABASE")
			dbUser := os.Getenv("DB_USER")
			dbPassword := os.Getenv("DB_PASSWORD")

			// postgresql://localhost/postgres?user=postgres&password=mysecretpassword"
			connUri := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", dbHost, dbDatabase, dbUser, dbPassword)
			logger.Debug("database.go", "connUri", connUri)
			conn, err := pgxpool.New(context.Background(), connUri)
			if err != nil {
				logger.Error("Unable to connect to database: %v\n", "error", err.Error())
				os.Exit(1)
			}

			logger.Info("Connected to the database")
			db.Conn = conn

			return nil
		},
		OnStop: func(ctx context.Context) error {
			db.Conn.Close()

			return nil
		},
	})

	return db
}
