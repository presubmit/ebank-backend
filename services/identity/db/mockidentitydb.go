package db

import (
	"context"
	m "ebank/services/identity/models"
	"ebank/shared/db"
)

type MockIdentityDB struct {
	CreateUserFunc           func(*m.User) (string, error)
	GetUserFunc              func(string) (*m.User, error)
	GetUserByEmailFunc       func(string) (*m.User, error)
	CreateCompanyFunc        func(*m.Company) (string, error)
	GetCompaniesByUserIdFunc func(string) ([]*m.Company, error)
	CreateEmployeeFunc       func(*m.Employee) (string, error)
	GetEmployeeByUserIdFunc  func(string, string) (*m.Employee, error)
}

func (*MockIdentityDB) ExecTx(_ context.Context, f func(db.Tx) error) error {
	return f(nil)
}

func (m *MockIdentityDB) CreateUser(_ db.Tx, user *m.User) (string, error) {
	return m.CreateUserFunc(user)
}

func (m *MockIdentityDB) GetUserByEmail(_ db.Tx, email string) (*m.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockIdentityDB) GetUser(_ db.Tx, userId string) (*m.User, error) {
	return m.GetUserFunc(userId)
}
func (m *MockIdentityDB) CreateCompany(_ db.Tx, company *m.Company) (string, error) {
	return m.CreateCompanyFunc(company)
}

func (m *MockIdentityDB) CreateEmployee(_ db.Tx, e *m.Employee) (string, error) {
	return m.CreateEmployeeFunc(e)
}

func (m *MockIdentityDB) GetEmployeeByUserId(_ db.Tx, userId, companyId string) (*m.Employee, error) {
	return m.GetEmployeeByUserIdFunc(userId, companyId)
}

func (m *MockIdentityDB) GetCompaniesByUserId(_ db.Tx, userId string) ([]*m.Company, error) {
	return m.GetCompaniesByUserIdFunc(userId)
}
