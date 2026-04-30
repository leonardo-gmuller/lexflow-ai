package uow

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type stubTx struct{}

// Pode devolver ele mesmo como “nested”
func (stubTx) Begin(ctx context.Context) (pgx.Tx, error) { return stubTx{}, nil }
func (stubTx) Commit(ctx context.Context) error          { return nil }
func (stubTx) Rollback(ctx context.Context) error        { return nil }

func (stubTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (stubTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}
func (stubTx) LargeObjects() pgx.LargeObjects { return pgx.LargeObjects{} }

func (stubTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (stubTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (stubTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) { return nil, nil }
func (stubTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row        { return nil }

func (stubTx) Conn() *pgx.Conn { return nil }

func NewTxWithFactories(factories map[string]RepoFactory) *UowPGX {
	tx := stubTx{}
	return &UowPGX{
		tx:           tx,
		repositories: factories,
	}
}
