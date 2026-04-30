package uow

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type RepoFactory func(db DBTX) any

type UnitOfWorkPGXInterface interface {
	Register(name string, factory RepoFactory)
	Unregister(name string)
	GetRepositoryTx(name string) (any, error)
	GetRepositoryFromPool(name string) (any, error)
	Do(ctx context.Context, fn func(ctx context.Context, ux *UowPGX) error) error
}

type UowPGX struct {
	Pool         *pgxpool.Pool
	tx           pgx.Tx
	repositories map[string]RepoFactory
}

func NewUowPGX(pool *pgxpool.Pool) *UowPGX {
	return &UowPGX{
		Pool:         pool,
		repositories: make(map[string]RepoFactory),
	}
}

func (u *UowPGX) Register(name string, factory RepoFactory) {
	u.repositories[name] = factory
}

func (u *UowPGX) Unregister(name string) {
	delete(u.repositories, name)
}

// ---------- Resolve repo a partir da Tx (somente dentro de Do) ----------
func (u *UowPGX) GetRepositoryTx(name string) (any, error) {
	fc, ok := u.repositories[name]
	if !ok {
		return nil, fmt.Errorf("uow: repo %q not registered", name)
	}
	if u.tx == nil {
		return nil, fmt.Errorf("uow: GetRepositoryTx called outside of transaction")
	}
	return fc(u.tx), nil
}

// ---------- Resolve repo direto do Pool (leitura/rápidos sem Tx) ----------
func (u *UowPGX) GetRepositoryFromPool(name string) (any, error) {
	fc, ok := u.repositories[name]
	if !ok {
		return nil, fmt.Errorf("uow: repo %q not registered", name)
	}
	// pgxpool.Pool implementa DBTX, então podemos passar Pool diretamente
	return fc(u.Pool), nil
}

// Helpers genéricos como funções (não métodos)
func RepoTx[T any](ux UnitOfWorkPGXInterface, name string) (T, error) {
	var zero T
	r, err := ux.GetRepositoryTx(name)
	if err != nil {
		return zero, err
	}
	inst, ok := r.(T)
	if !ok {
		return zero, fmt.Errorf("uow: repo %q has different type", name)
	}
	return inst, nil
}

func RepoPool[T any](u UnitOfWorkPGXInterface, name string) (T, error) {
	var zero T
	r, err := u.GetRepositoryFromPool(name)
	if err != nil {
		return zero, err
	}
	inst, ok := r.(T)
	if !ok {
		return zero, fmt.Errorf("uow: repo %q has different type", name)
	}
	return inst, nil
}

func (u *UowPGX) Do(ctx context.Context, fn func(ctx context.Context, ux *UowPGX) error) error {
	if u.tx != nil {
		return fmt.Errorf("uow: transaction already started")
	}
	tx, err := u.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("uow.begin: %w", err)
	}
	u.tx = tx
	defer func() { _ = u.rollbackSilently(ctx) }()

	if err := fn(ctx, u); err != nil {
		if rbErr := u.tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("uow.exec: %w (rollback: %v)", err, rbErr)
		}
		u.tx = nil
		return err
	}
	if err := u.tx.Commit(ctx); err != nil {
		if rbErr := u.tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("uow.commit: %w (rollback: %v)", err, rbErr)
		}
		u.tx = nil
		return err
	}
	u.tx = nil
	return nil
}

func (u *UowPGX) rollbackSilently(ctx context.Context) error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Rollback(ctx)
	u.tx = nil
	return err
}
