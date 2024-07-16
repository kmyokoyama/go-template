package outbound

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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

func (db *PgDatabase) CreateUser(user models.User, hashedPassword string) error {
	tx, err := db.Conn.Begin(context.Background())

	if err != nil {
		return err
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	stmt := `INSERT INTO users(uuid, username, password, role_id) VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))`
	_, err = tx.Exec(context.Background(), stmt, user.Id, user.Username, hashedPassword, user.Role.String())
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
	var username string
	var role string

	err := db.Conn.QueryRow(context.Background(), "SELECT username, role FROM users WHERE users.uuid = $1", id).Scan(&username, &role)
	if err != nil {
		return models.User{}, err
	}

	modelRole, err := models.RoleFromString(role)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Id: id, Username: username, Role: modelRole}, nil
}

func (db *PgDatabase) FindUserAndPasswordByUsername(username string) (models.User, string, error) {
	var id uuid.UUID
	var hashedPassword string
	var role string

	err := db.Conn.QueryRow(
		context.Background(),
		"SELECT u.uuid, u.password, r.name FROM users AS u JOIN roles AS r ON u.role_id = r.id WHERE u.username = $1", username,
	).Scan(&id, &hashedPassword, &role)
	if err != nil {
		return models.User{}, "", err
	}

	modelRole, err := models.RoleFromString(role)
	if err != nil {
		return models.User{}, "", err
	}

	return models.User{Id: id, Username: username, Role: modelRole}, hashedPassword, nil
}

func NewDatabase(lc fx.Lifecycle, logger *slog.Logger, config components.Config) components.Database {
	db := &PgDatabase{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Started on", "port", ":8080")

			dbHost := config.Get("DB_HOST")
			dbDatabase := config.Get("DB_DATABASE")
			dbUser := config.Get("DB_USER")
			dbPassword := config.Get("DB_PASSWORD")

			// postgresql://localhost/postgres?user=postgres&password=mysecretpassword"
			connUri := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", dbHost, dbDatabase, dbUser, dbPassword)
			logger.Info("database.go", "connUri", connUri)
			conn, err := pgxpool.New(context.Background(), connUri)
			if err != nil {
				logger.Error("Unable to connect to database: %v\n", "error", err.Error())
				os.Exit(1)
			}

			logger.Info("connected to the database")
			db.Conn = conn

			pwd, _ := os.Getwd()
			logger.Info("running migrations", "dir", pwd)
			driver, err := postgres.WithInstance(stdlib.OpenDBFromPool(conn), &postgres.Config{})
			if err != nil {
				logger.Error("migrate/postgres instance failed to open")
				return err
			}
			m, err := migrate.NewWithDatabaseInstance(
				"file://migrations/",
				"postgres", driver)
			if err != nil {
				logger.Error("migrate failed to get instance")
				return err
			}
			defer m.Close()

			m.Up()
			logger.Info("migrations ran")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			db.Conn.Close()

			return nil
		},
	})

	return db
}

// Verify interface compliance.
var _ components.Database = (*PgDatabase)(nil)
