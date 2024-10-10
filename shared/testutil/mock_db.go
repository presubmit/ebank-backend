package testutil

import "database/sql"

type MockTx struct{}

func (*MockTx) Commit() error {
	return nil
}

func (*MockTx) Rollback() error {
	return nil
}

func (*MockTx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return stmt
}
