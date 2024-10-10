package db

import (
	"context"
	"database/sql"
	m "ebank/services/card/models"
	"ebank/shared/db"
	"ebank/shared/errors"
	"ebank/shared/utils"

	"github.com/google/uuid"
)

type CardDB interface {
	ExecTx(context.Context, func(db.Tx) error) error
	CreateCard(db.Tx, *m.Card) (string, error)
	GetCardsByCompany(db.Tx, string) ([]*m.Card, error)
	GetCard(db.Tx, string, string) (*m.Card, error)
	FreezeCard(db.Tx, string, string) error
	UnfreezeCard(db.Tx, string, string) error
	CloseCard(db.Tx, string, string) error
}

type DB struct {
	*db.DB
}

func (d *DB) CreateCard(tx db.Tx, c *m.Card) (string, error) {
	cardId := uuid.New().String()
	_, err := d.ExecStmt(tx, createCardStmt, cardId, c.EmployeeId, c.CompanyId, c.CreatedBy, c.Brand,
		c.Number, c.ExpirationMonth, c.ExpirationYear, c.SecurityCode, c.Type)
	if err != nil {
		return "", err
	}
	return cardId, nil
}

func scanCard(row db.Row) (*m.Card, error) {
	c := &m.Card{}
	var fronzenAt sql.NullString
	if err := row.Scan(
		&c.ID,
		&c.Brand,
		&c.Number,
		&c.ExpirationMonth,
		&c.ExpirationYear,
		&c.SecurityCode,
		&c.Type,
		&fronzenAt,
		&c.CreatedBy,
	); err != nil {
		return nil, err
	}
	c.FrozenAt = utils.NullStringToString(fronzenAt)
	return c, nil
}

func (d *DB) GetCardsByCompany(tx db.Tx, companyId string) ([]*m.Card, error) {
	rows, err := d.QueryStmt(tx, getCardsByCompanyStmt, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*m.Card
	for rows.Next() {
		c, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (d *DB) GetCard(tx db.Tx, cardId, companyId string) (*m.Card, error) {
	row := d.QueryRowStmt(tx, getCardByIdStmt, cardId, companyId)
	return scanCard(row)
}

func (d *DB) FreezeCard(tx db.Tx, cardId, companyId string) error {
	res, err := d.ExecStmt(tx, freezeCardByIdStmt, cardId, companyId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.InvalidArgumentf("invalid card state")
	}
	return nil
}

func (d *DB) UnfreezeCard(tx db.Tx, cardId, companyId string) error {
	res, err := d.ExecStmt(tx, unfreezeCardByIdStmt, cardId, companyId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.InvalidArgumentf("invalid card state")
	}
	return nil
}

func (d *DB) CloseCard(tx db.Tx, cardId, companyId string) error {
	res, err := d.ExecStmt(tx, closeCardByIdStmt, cardId, companyId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.InvalidArgumentf("invalid card state")
	}
	return nil
}

type statementKey string

var (
	createCardStmt        = statementKey("createCardStmt")
	getCardsByCompanyStmt = statementKey("getCardsByCompanyStmt")
	getCardByIdStmt       = statementKey("getCardByIdStmt")
	freezeCardByIdStmt    = statementKey("freezeCardByIdStmt")
	unfreezeCardByIdStmt  = statementKey("unfreezeCardByIdStmt")
	closeCardByIdStmt     = statementKey("closeCardByIdStmt")
)

func (d *DB) RegisterStmts() error {
	getCardBaseStmt := `SELECT id, brand, number, expiration_month, expiration_year, security_code,
	type, frozen_at, created_by FROM cards `
	return d.Register(map[interface{}]string{
		createCardStmt: db.InsertStmt("cards", []string{
			"id", "employee_id", "company_id", "created_by", "brand", "number",
			"expiration_month", "expiration_year", "security_code", "type",
		}),
		getCardsByCompanyStmt: getCardBaseStmt + `WHERE company_id = ? AND closed_at IS NULL`,
		getCardByIdStmt:       getCardBaseStmt + `WHERE id = ? AND company_id = ?`,
		freezeCardByIdStmt:    `UPDATE cards SET frozen_at = NOW() WHERE id = ? AND frozen_at IS NULL AND company_id = ?`,
		unfreezeCardByIdStmt:  `UPDATE cards SET frozen_at = NULL WHERE id = ? AND frozen_at IS NOT NULL AND company_id = ?`,
		closeCardByIdStmt:     `UPDATE cards SET closed_at = NOW() WHERE id = ? AND closed_at IS NULL AND company_id = ?`,
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
