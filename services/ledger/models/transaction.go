package models

import (
	pb "ebank/pb/services/ledger"
	"ebank/shared/errors"
	"ebank/shared/money"
)

type Transaction struct {
	ID             string
	LegID          string
	CompanyID      string
	Amount         int64
	Currency       string
	AccountID      string
	CounterpartyID string
	OtherAccountID string
	Description    string
	Type           string
	AfterBalance   int64
	CreatedBy      string
	CreatedAt      string
}

func (t *Transaction) Validate() error {
	if len(t.AccountID) == 0 {
		return errors.InvalidArgumentf("invalid account_id")
	}
	if !money.IsCurrencyValid(t.Currency) {
		return errors.InvalidArgumentf("invalid currency")
	}
	return nil
}

func (t *Transaction) ToProto() *pb.Transaction {
	return &pb.Transaction{
		Id:             t.ID,
		LegId:          t.LegID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		AccountId:      t.AccountID,
		CounterpartyId: t.CounterpartyID,
		OtherAccountId: t.OtherAccountID,
		Description:    t.Description,
		Type:           t.Type,
		AfterBalance:   t.AfterBalance,
		CreatedBy:      t.CreatedBy,
		CreatedAt:      t.CreatedAt,
	}
}
