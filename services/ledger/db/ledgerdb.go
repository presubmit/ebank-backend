package db

import (
	"context"
	"database/sql"
	m "ebank/services/ledger/models"
	"ebank/shared/db"

	"github.com/google/uuid"
)

type LedgerDB interface {
	ExecTx(context.Context, func(db.Tx) error) error

	CreateAccount(db.Tx, *m.Account) (string, error)
	GetAccount(db.Tx, string, string) (*m.Account, error)
	GetAccountsByCompany(db.Tx, string) ([]*m.Account, error)
	AddToAccountBalance(db.Tx, int64, string) (*m.Account, error)

	GetTransactions(db.Tx, string) ([]*m.Transaction, error)
	CreateTransaction(db.Tx, *m.Transaction) (string, error)

	CreateCounterparty(db.Tx, *m.Counterparty) (string, error)
	GetCounterparty(db.Tx, string, string) (*m.Counterparty, error)
	GetCounterparties(db.Tx, string) ([]*m.Counterparty, error)
}

type DB struct {
	*db.DB
}

func (d *DB) CreateAccount(tx db.Tx, a *m.Account) (string, error) {
	accountId := uuid.New().String()
	_, err := d.ExecStmt(tx, createAccountStmt, accountId, a.Name, a.Currency, a.CompanyID)
	if err != nil {
		return "", err
	}
	return accountId, nil
}

func scanAccount(row db.Row) (*m.Account, error) {
	a := &m.Account{}
	if err := row.Scan(&a.ID, &a.Name, &a.Balance, &a.Currency, &a.CreatedAt); err != nil {
		return nil, err
	}
	return a, nil
}

func (d *DB) GetAccount(tx db.Tx, id, companyId string) (*m.Account, error) {
	row := d.QueryRowStmt(tx, getAccountStmt, id, companyId)
	return scanAccount(row)
}

func (d *DB) GetAccountsByCompany(tx db.Tx, companyId string) ([]*m.Account, error) {
	rows, err := d.QueryStmt(tx, getAccountsByCompanyStmt, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*m.Account
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (d *DB) AddToAccountBalance(tx db.Tx, amount int64, accountId string) (*m.Account, error) {
	row := d.QueryRowStmt(tx, addToAccountBalanceStmt, amount, accountId, amount)
	return scanAccount(row)
}

func scanTransaction(row db.Row) (*m.Transaction, error) {
	t := &m.Transaction{}
	var otherAccountId sql.NullString
	var counterpartyId sql.NullString
	if err := row.Scan(
		&t.ID,
		&t.LegID,
		&t.AccountID,
		&otherAccountId,
		&counterpartyId,
		&t.Description,
		&t.Amount,
		&t.Currency,
		&t.Type,
		&t.AfterBalance,
		&t.CreatedBy,
		&t.CreatedAt,
	); err != nil {
		return nil, err
	}
	t.OtherAccountID = otherAccountId.String
	t.CounterpartyID = counterpartyId.String
	return t, nil
}

func (d *DB) GetTransactions(tx db.Tx, companyId string) ([]*m.Transaction, error) {
	rows, err := d.QueryStmt(tx, getTransactionsStmt, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*m.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (d *DB) CreateTransaction(tx db.Tx, t *m.Transaction) (string, error) {
	_, err := d.ExecStmt(tx, createTransactionStmt,
		t.ID,
		t.LegID,
		t.CompanyID,
		t.AccountID,
		t.OtherAccountID,
		t.CounterpartyID,
		t.Amount,
		t.Currency,
		t.Description,
		t.Type,
		"completed",
		t.AfterBalance,
		t.CreatedBy,
	)
	if err != nil {
		return "", err
	}
	return t.ID, nil
}

func (d *DB) CreateCounterparty(tx db.Tx, c *m.Counterparty) (string, error) {
	id := uuid.New().String()
	_, err := d.ExecStmt(tx, createCounterpartyStmt,
		id,
		c.CompanyID,
		c.FirstName,
		c.LastName,
		c.CompanyName,
		c.Type,
		c.IBAN,
		c.Currency,
		c.Country,
		c.CreatedBy,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func scanCounterparty(row db.Row) (*m.Counterparty, error) {
	c := &m.Counterparty{}
	if err := row.Scan(
		&c.ID,
		&c.CompanyID,
		&c.FirstName,
		&c.LastName,
		&c.CompanyName,
		&c.Type,
		&c.IBAN,
		&c.Currency,
		&c.Country,
		&c.CreatedBy,
		&c.CreatedAt,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (d *DB) GetCounterparty(tx db.Tx, id, companyId string) (*m.Counterparty, error) {
	row := d.QueryRowStmt(tx, getCounterpartyStmt, id, companyId)
	return scanCounterparty(row)
}

func (d *DB) GetCounterparties(tx db.Tx, companyId string) ([]*m.Counterparty, error) {
	rows, err := d.QueryStmt(tx, getCounterpartiesStmt, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counterparties []*m.Counterparty
	for rows.Next() {
		c, err := scanCounterparty(rows)
		if err != nil {
			return nil, err
		}
		counterparties = append(counterparties, c)
	}
	return counterparties, nil
}

type statementKey string

var (
	createAccountStmt        = statementKey("createAccountStmt")
	addToAccountBalanceStmt  = statementKey("addToAccountBalanceStmt")
	getAccountStmt           = statementKey("getAccountStmt")
	getAccountsByCompanyStmt = statementKey("getAccountsByCompanyStmt")
	getTransactionsStmt      = statementKey("getTransactionsStmt")
	createTransactionStmt    = statementKey("createTransactionStmt")
	createCounterpartyStmt   = statementKey("createCounterpartyStmt")
	getCounterpartyStmt      = statementKey("getCounterpartyStmt")
	getCounterpartiesStmt    = statementKey("getCounterpartiesStmt")
)

func (d *DB) RegisterStmts() error {
	getTransactionBase := `
		SELECT 
			id, 
			leg_id,
			account_id,
			other_account_id,
			counterparty_id,
			COALESCE(description, ''),
			amount, 
			currency,
			type,
			after_balance,
			created_by, 
			created_at
		FROM transactions
	`
	accountCols := `id, name, balance, currency, created_at`
	getAccountBase := `SELECT ` + accountCols + ` FROM accounts`
	getCounterpartyBaseStmt := `
		SELECT 
			id, 
			company_id, 
			COALESCE(first_name, ''),
			COALESCE(last_name, ''),
			COALESCE(company_name, ''),
			type, 
			iban, 
			currency, 
			country, 
			created_by, 
			created_at
		FROM counterparties
	`

	return d.Register(map[interface{}]string{
		createAccountStmt:        `INSERT INTO accounts (id, name, currency, company_id) VALUES (?, ?, ?, ?)`,
		getAccountStmt:           getAccountBase + ` WHERE id = ? AND company_id = ? FOR UPDATE`,
		getAccountsByCompanyStmt: getAccountBase + ` WHERE company_id = ?`,
		addToAccountBalanceStmt:  `UPDATE accounts SET balance = balance + ? WHERE id = ? AND balance + ? >= 0 RETURNING ` + accountCols,

		getTransactionsStmt: getTransactionBase + " WHERE company_id = ? ORDER BY created_at DESC",
		createTransactionStmt: db.InsertStmt("transactions", []string{
			"id", "leg_id", "company_id", "account_id", "other_account_id", "counterparty_id",
			"amount", "currency", "description", "type", "status",
			"after_balance", "created_by",
		}),

		createCounterpartyStmt: db.InsertStmt("counterparties", []string{
			"id", "company_id", "first_name",
			"last_name", "company_name", "type",
			"iban", "currency", "country", "created_by",
		}),
		getCounterpartyStmt:   getCounterpartyBaseStmt + " WHERE id = ? AND company_id = ?",
		getCounterpartiesStmt: getCounterpartyBaseStmt + " WHERE company_id = ?",
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
