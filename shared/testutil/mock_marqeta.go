package testutil

import (
	mq "github.com/ebankro/marqeta"
)

type MockMarqeta struct {
	CreateBusinessFunc func(body mq.BusinessCardholder) (mq.BusinessCardHolderResponse, error)
	CreateUserFunc     func(body mq.CardHolderModel) (mq.UserCardHolderResponse, error)
	CreateCardFunc     func(body mq.CardRequest) (mq.CardResponse, error)
}

func (m *MockMarqeta) CreateBusiness(body mq.BusinessCardholder) (mq.BusinessCardHolderResponse, error) {
	return m.CreateBusinessFunc(body)
}

func (m *MockMarqeta) CreateUser(body mq.CardHolderModel) (mq.UserCardHolderResponse, error) {
	return m.CreateUserFunc(body)
}

func (m *MockMarqeta) CreateCard(body mq.CardRequest) (mq.CardResponse, error) {
	return m.CreateCardFunc(body)
}
