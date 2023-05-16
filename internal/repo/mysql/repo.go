package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ctxKey string

// CtxTxKey is key for context transaction.
const CtxTxKey ctxKey = "tx"

type Repo struct {
	conn *sql.DB
}

type Executor interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Connection struct {
	Executor
}

func NewRepo(db *sql.DB) Repo {
	return Repo{
		conn: db,
	}
}

func (r Repo) RunInTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	ctx, err := BeginTransaction(ctx, r.conn, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := f(ctx); err != nil {
		if errRollback := RollbackTransaction(ctx); errRollback != nil {
			return fmt.Errorf("failed to rollback, rollback error %s: %w", errRollback.Error(), err)
		}
		return err
	}

	if err := CommitTransaction(ctx); err != nil {
		if errRollback := RollbackTransaction(ctx); errRollback != nil {
			return fmt.Errorf("failed to rollback after commit, rollback error %s: %w", errRollback.Error(), err)
		}
		return fmt.Errorf("hfpg failed to commit: %w", err)
	}
	return nil
}

// BeginTransaction from the context.
func BeginTransaction(ctx context.Context, db *sql.DB, opts *sql.TxOptions) (context.Context, error) {
	if _, err := GetTransaction(ctx); err == nil {
		return ctx, errors.New("transaction already started")
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, CtxTxKey, tx), nil
}

// GetTransaction will get transaction from context.
func GetTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(CtxTxKey).(*sql.Tx)
	if !ok {
		return nil, errors.New("no transaction found in context")
	}

	return tx, nil
}

// RollbackTransaction from the context.
func RollbackTransaction(ctx context.Context) error {
	tx, err := GetTransaction(ctx)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

// CommitTransaction from the context.
func CommitTransaction(ctx context.Context) error {
	tx, err := GetTransaction(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r Repo) GetExecutor(ctx context.Context) Executor {
	if tx, err := GetTransaction(ctx); err == nil {
		return Connection{Executor: tx}
	}
	return Connection{Executor: r.conn}
}
