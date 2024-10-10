package db

import (
	"context"
	m "ebank/services/ledger/models"
	"ebank/shared/db"
)

type MockLedgerDB struct {
	CreateAccountFunc        func(*m.Account) (string, error)
	GetAccountFunc           func(string, string) (*m.Account, error)
	GetAccountsByCompanyFunc func(string) ([]*m.Account, error)
	AddToAccountBalanceFunc  func(int64, string) (*m.Account, error)
	GetTransactionsFunc      func(string) ([]*m.Transaction, error)
	CreateTransactionFunc    func(*m.Transaction) (string, error)
	CreateCounterpartyFunc   func(*m.Counterparty) (string, error)
	GetCounterpartyFunc      func(string, string) (*m.Counterparty, error)
	GetCounterpartiesFunc    func(string) ([]*m.Counterparty, error)
}

func (*MockLedgerDB) ExecTx(_ context.Context, f func(db.Tx) error) error {
	return f(nil)
}

func (m *MockLedgerDB) CreateAccount(_ db.Tx, a *m.Account) (string, error) {
	return m.CreateAccountFunc(a)
}

func (m *MockLedgerDB) GetAccount(_ db.Tx, id, companyId string) (*m.Account, error) {
	return m.GetAccountFunc(id, companyId)
}

func (m *MockLedgerDB) GetAccountsByCompany(_ db.Tx, c string) ([]*m.Account, error) {
	return m.GetAccountsByCompanyFunc(c)
}

func (m *MockLedgerDB) AddToAccountBalance(_ db.Tx, amount int64, id string) (*m.Account, error) {
	return m.AddToAccountBalanceFunc(amount, id)
}

func (m *MockLedgerDB) GetTransactions(_ db.Tx, c string) ([]*m.Transaction, error) {
	return m.GetTransactionsFunc(c)
}

func (m *MockLedgerDB) CreateTransaction(_ db.Tx, t *m.Transaction) (string, error) {
	return m.CreateTransactionFunc(t)
}

func (m *MockLedgerDB) CreateCounterparty(_ db.Tx, c *m.Counterparty) (string, error) {
	return m.CreateCounterpartyFunc(c)
}

func (m *MockLedgerDB) GetCounterparty(_ db.Tx, id, companyId string) (*m.Counterparty, error) {
	return m.GetCounterpartyFunc(id, companyId)
}

func (m *MockLedgerDB) GetCounterparties(_ db.Tx, companyId string) ([]*m.Counterparty, error) {
	return m.GetCounterpartiesFunc(companyId)
}
