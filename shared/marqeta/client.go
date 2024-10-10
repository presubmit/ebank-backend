package marqeta

import (
	"context"

	"github.com/antihax/optional"
	mq "github.com/ebankro/marqeta"
)

type Api interface {
	CreateBusiness(body mq.BusinessCardholder) (mq.BusinessCardHolderResponse, error)
	CreateUser(body mq.CardHolderModel) (mq.UserCardHolderResponse, error)
	CreateCard(body mq.CardRequest) (mq.CardResponse, error)
}

type Marqeta struct {
	apiClient *mq.APIClient
	ctx       context.Context
}

func NewClient() *Marqeta {
	apiClient := mq.NewAPIClient(&mq.Configuration{
		BasePath: "https://sandbox-api.marqeta.com/v3",
	})
	ctx := context.WithValue(context.Background(), mq.ContextBasicAuth, mq.BasicAuth{
		UserName: "74bde80c-cfc3-47c7-b320-07ef1866e9d3",
		Password: "ea910f5e-9f1e-48dc-b1d5-ff2ecf04fb8b",
	})
	return &Marqeta{
		apiClient: apiClient,
		ctx:       ctx,
	}
}

func (m *Marqeta) CreateBusiness(body mq.BusinessCardholder) (mq.BusinessCardHolderResponse, error) {
	res, _, err := m.apiClient.BusinessesApi.PostBusinesses(m.ctx, &mq.BusinessesApiPostBusinessesOpts{
		Body: optional.NewInterface(body),
	})
	return res, err
}

func (m *Marqeta) CreateUser(body mq.CardHolderModel) (mq.UserCardHolderResponse, error) {
	res, _, err := m.apiClient.UsersApi.PostUsers(m.ctx, &mq.UsersApiPostUsersOpts{
		Body: optional.NewInterface(body),
	})
	return res, err
}

func (m *Marqeta) CreateCard(body mq.CardRequest) (mq.CardResponse, error) {
	res, _, err := m.apiClient.CardsApi.PostCards(m.ctx, &mq.CardsApiPostCardsOpts{
		Body:          optional.NewInterface(body),
		ShowCvvNumber: optional.NewBool(true),
		ShowPan:       optional.NewBool(true),
	})
	return res, err
}
