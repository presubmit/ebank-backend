package models

import (
	pb "ebank/pb/services/ledger"
	"ebank/shared/money"
	"errors"
)

type Account struct {
	ID        string
	Name      string
	Currency  string
	Balance   int64
	CompanyID string
	CreatedAt string
}

func (a *Account) ValidateFields() error {
	// validate company id
	if len(a.CompanyID) == 0 {
		return errors.New("invalid company id")
	}
	// validate currency
	if !money.IsCurrencyValid(a.Currency) {
		return errors.New("invalid currency")
	}
	// validate name
	if len(a.Name) == 0 {
		return errors.New("invalid name")
	}
	return nil
}

func (a *Account) ToProto() *pb.Account {
	return &pb.Account{
		Id:       a.ID,
		Name:     a.Name,
		Balance:  a.Balance,
		Currency: a.Currency,
	}
}
