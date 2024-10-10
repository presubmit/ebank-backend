package db

import (
	"context"
	"ebank/services/auth/tokens"
	"ebank/shared/db"

	"github.com/golang/glog"
)

type AuthDB interface {
	BeginTx(ctx context.Context) (db.Tx, error)
	SaveRefreshToken(db.Tx, *tokens.RefreshToken) error
	GetRefreshToken(db.Tx, string) (*tokens.RefreshToken, error)
	DeleteRefreshToken(db.Tx, string) error
}

type DB struct {
	*db.DB
}

func (d *DB) SaveRefreshToken(tx db.Tx, rt *tokens.RefreshToken) error {
	_, err := d.ExecStmt(tx, createTokenStmt, rt.ID, rt.UserID, rt.AccessTokenID, rt.ExpiresAt)
	return err
}

func (d *DB) DeleteRefreshToken(tx db.Tx, id string) error {
	info, err := d.ExecStmt(tx, deleteTokenStmt, id)
	aff, e := info.RowsAffected()
	if e != nil {
		glog.Error(e)
	} else {
		glog.Info("Rows affected: ", aff)
	}
	return err
}

func (d *DB) GetRefreshToken(tx db.Tx, id string) (*tokens.RefreshToken, error) {
	row := d.QueryRowStmt(tx, getTokenStmt, id)
	rt := &tokens.RefreshToken{}
	err := row.Scan(&rt.ID, &rt.UserID, &rt.AccessTokenID, &rt.ExpiresAt)
	return rt, err
}

type statementKey string

var (
	createTokenStmt = statementKey("createTokenStmt")
	getTokenStmt    = statementKey("getTokenStmt")
	deleteTokenStmt = statementKey("deleteTokenStmt")
)

func (d *DB) RegisterStmts() error {
	getTokenBaseStmt := `SELECT id, user_id, access_token_id, expires_at FROM refresh_tokens `
	return d.Register(map[interface{}]string{
		createTokenStmt: `INSERT INTO refresh_tokens (id, user_id, access_token_id, expires_at) VALUES (?, ?, ?, ?)`,
		getTokenStmt:    getTokenBaseStmt + `WHERE id = ?`,
		deleteTokenStmt: `DELETE FROM refresh_tokens WHERE id = ?`,
	})
}

func New() *DB {
	db := &DB{
		DB: db.New(),
	}
	if err := db.RegisterStmts(); err != nil {
		panic(err)
	}
	return db
}
