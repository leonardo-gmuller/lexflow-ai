package postgres

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations
var MigrationsFS embed.FS

type Client struct {
	Pool *pgxpool.Pool
}

// DBTX é a interface comum que tanto *pgxpool.Pool quanto pgx.Tx implementam.
// É o tipo que vamos passar pros repositórios (sqlc.New(db)).
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (c *Client) Close() {
	if c != nil && c.Pool != nil {
		c.Pool.Close()
	}
}

// New cria o pool PGX, roda as migrations e devolve o Client.
func New(ctx context.Context, cfg config.Postgres) (*Client, error) {
	const operation = "Postgres.New"

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// ---- Migrations (mesma lógica que você tinha, só com alias no pacote) ----

	sqlStdlibDB := stdlib.OpenDB(*pgxConfig.ConnConfig)

	driver, err := pgxmigrate.WithInstance(sqlStdlibDB, &pgxmigrate.Config{
		DatabaseName: pgxConfig.ConnConfig.Database,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	source, err := httpfs.New(http.FS(MigrationsFS), "migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	migration, err := migrate.NewWithInstance("httpfs", source, pgxConfig.ConnConfig.Database, driver)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// se já estiver up-to-date, Up() retorna ErrNoChange, que podemos ignorar
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	srcErr, dbErr := migration.Close()
	if srcErr != nil {
		return nil, fmt.Errorf("%s: %w", operation, srcErr)
	}
	if dbErr != nil {
		return nil, fmt.Errorf("%s: %w", operation, dbErr)
	}

	return &Client{Pool: pool}, nil
}

// WithTx abre uma transação e injeta um DBTX (a tx) no callback.
// Aqui é o ponto que vamos usar no webhook/worker:
//
//	err := pgClient.WithTx(ctx, func(ctx context.Context, db postgres.DBTX) error {
//	    waUC := app.NewWaUseCase(db)
//	    return waUC.HandleIncomingMessage(ctx, in)
//	})
func (c *Client) WithTx(
	ctx context.Context,
	fn func(ctx context.Context, db DBTX) error,
) error {
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // se der erro no callback, rollback automático

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DB expõe o pool como DBTX, para operações sem transação explícita (se você quiser usar).
func (c *Client) DB() DBTX {
	return c.Pool
}
