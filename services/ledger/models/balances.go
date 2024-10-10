package models

import (
	"ebank/shared/errors"
)

type Amount struct {
	Value    int64
	Currency string
}

func (a *Amount) Add(am *Amount) (*Amount, error) {
	if a.Currency != am.Currency {
		return nil, errors.Internalf("cannot add amounts in different currencies")
	}
	return &Amount{
		Value:    am.Value + a.Value,
		Currency: a.Currency,
	}, nil
}
