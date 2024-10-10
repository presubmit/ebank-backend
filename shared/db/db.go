package db

import (
	"context"
	"database/sql"
	"ebank/shared/microservice"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	pqDuplicateEntryErrorCode = pq.ErrorCode("23505")

	ErrNotFound       = errors.New("not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

type Tx interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback rolls back the transaction.
	Rollback() error
	// Stmt returns a transaction-specific prepared statement from
	// an existing statement.
	Stmt(stmt *sql.Stmt) *sql.Stmt
}

type DB struct {
	DB    *sql.DB
	stmts map[interface{}]*sql.Stmt
}

func New() *DB {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		microservice.DbUser,
		microservice.DbPass,
		microservice.DbHost,
		microservice.DbPort,
		microservice.DbName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return &DB{
		DB:    db,
		stmts: map[interface{}]*sql.Stmt{},
	}
}

// Close all database statements.
func (d *DB) Close() {
	for _, stmt := range d.stmts {
		stmt.Close()
	}
}

func (d *DB) Register(stmts map[interface{}]string) error {
	for key, stmtStr := range stmts {
		if _, found := d.stmts[key]; found {
			d.Close()
			return fmt.Errorf("statement has duplicate key: %q", key)
		}
		stmt, err := d.DB.Prepare(postgresStmt(stmtStr))
		if err != nil {
			d.Close()
			return fmt.Errorf("prepare failed for statement with key %q: %v", key, err.Error())
		}
		d.stmts[key] = stmt
	}
	return nil
}

func postgresStmt(str string) string {
	// Postgres uses $1, $2, $3 ... for prepared statement parameters.
	// However, most other SQL databases use ? so we need to convert
	// unnamed parameters (?) to numbered parameters ($1).
	k := 1
	for {
		newStr := strings.Replace(str, "?", fmt.Sprintf("$%d", k), 1)
		if newStr == str {
			break
		}
		str = newStr
		k += 1
	}
	return str
}

type Row interface {
	Scan(dest ...interface{}) error
}

type rowErr struct {
	err error
}

func (r *rowErr) Scan(dest ...interface{}) error {
	return r.err
}

func maybeNull(v interface{}) interface{} {
	s, ok := v.(string)
	if !ok {
		return v
	}
	if len(s) == 0 {
		return sql.NullString{}
	}
	return s
}

func (d *DB) ExecStmt(tx Tx, key interface{}, args ...interface{}) (sql.Result, error) {
	stmt, err := d.stmtForKey(key)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == pqDuplicateEntryErrorCode {
			return nil, ErrDuplicateEntry
		}
		return nil, err
	}
	for i := range args {
		args[i] = maybeNull(args[i])
	}
	if tx != nil {
		return tx.Stmt(stmt).Exec(args...)
	}
	return stmt.Exec(args...)
}

func (d *DB) QueryStmt(tx Tx, key interface{}, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.stmtForKey(key)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.Stmt(stmt).Query(args...)
	}
	return stmt.Query(args...)
}

func (d *DB) QueryRowStmt(tx Tx, key interface{}, args ...interface{}) Row {
	stmt, err := d.stmtForKey(key)
	if err != nil {
		return &rowErr{err}
	}
	if tx != nil {
		return tx.Stmt(stmt).QueryRow(args...)
	}
	return stmt.QueryRow(args...)
}

func (d *DB) stmtForKey(key interface{}) (*sql.Stmt, error) {
	stmt, found := d.stmts[key]
	if !found {
		return nil, fmt.Errorf("no prepared statemet found for key %q", key)
	}
	return stmt, nil
}

func InsertStmt(table string, cols []string) string {
	colStr := strings.Join(cols, ", ")
	var vals []string
	for range cols {
		vals = append(vals, "?")
	}
	valStr := strings.Join(vals, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, colStr, valStr)
}

func (d *DB) BeginTx(ctx context.Context) (Tx, error) {
	return d.DB.BeginTx(ctx, nil)
}

func (d *DB) ExecTx(ctx context.Context, f func(Tx) error) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	if err := f(tx); err != nil {
		return err
	}

	return tx.Commit()
}
