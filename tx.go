package main

import (
	"context"
	"database/sql"
)

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TransactionScope define contract of a scope of a sql transaction
type TransactionScope interface {
	CreateNew(fn TxFn) error
}

// NewTransactionScope creates TransactionScope that connects to Db
func NewTransactionScope(db *sql.DB) TransactionScope {
	return &SQLTransactionScope{db}
}

// SQLTransactionScope is a representation of TransactionScope for Sql-based database
type SQLTransactionScope struct {
	db *sql.DB
}

// CreateNew creates new scope for this transaction
func (m *SQLTransactionScope) CreateNew(fn TxFn) error {
	tx, err := m.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	//tx, err := m.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error