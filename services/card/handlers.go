package card

import (
	"context"

	apb "ebank/pb/services/auth"
	pb "ebank/pb/services/card"
	"ebank/services/card/db"
	m "ebank/services/card/models"
	"ebank/shared/errors"
	"ebank/shared/marqeta"
	"ebank/shared/microservice"

	mq "github.com/ebankro/marqeta"
)

type Service struct {
	db          db.CardDB
	authService apb.AuthServiceClient
	marqeta     marqeta.Api
}

func NewService() *Service {
	return &Service{
		db:          db.New(),
		authService: microservice.Auth(),
		marqeta:     marqeta.NewClient(),
	}
}

func (s *Service) GetCardProducts(ctx context.Context, r *pb.GetCardProductsRequest) (*pb.CardProducts, error) {
	return &pb.CardProducts{
		Products: []*pb.CardProduct{
			{
				Id:    "mq-test-1",
				Brand: "sandbox",
				Type:  "virtual",
			},
		},
	}, nil
}

func (s *Service) Create(ctx context.Context, r *pb.CreateRequest) (*pb.Card, error) {
	empId, err := microservice.GetEmployeeId(ctx)
	if err != nil {
		return nil, errors.Internal(err)
	}
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}
	cardProducts, err := s.GetCardProducts(ctx, &pb.GetCardProductsRequest{})
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}

	var cardProduct *pb.CardProduct
	for _, p := range cardProducts.Products {
		if p.Id == r.GetCardProductId() {
			cardProduct = p
			break
		}
	}
	if cardProduct == nil {
		return nil, errors.InvalidArgumentf("Invalid card product id")
	}

	mqRes, err := s.marqeta.CreateCard(mq.CardRequest{
		CardProductToken: cardProduct.Id,
		UserToken:        r.GetCardholderId(),
		Fulfillment: &mq.Fulfillment{
			Shipping: &mq.Shipping{
				RecipientAddress: &mq.FulfillmentAddressRequest{
					Address1:   "2221 Broadway Street",
					City:       "Redwood City",
					State:      "CA",
					PostalCode: "94063",
				},
			},
		},
	})
	if err != nil {
		return nil, errors.Internalf("marqeta.CreateCard error: %v", err)
	}

	card := &m.Card{
		EmployeeId:      r.GetCardholderId(),
		CompanyId:       companyId,
		CreatedBy:       empId,
		Brand:           cardProduct.GetBrand(),
		Type:            cardProduct.GetType(),
		Number:          mqRes.Pan,
		ExpirationMonth: int32(mqRes.ExpirationTime.Month()),
		ExpirationYear:  int32(mqRes.ExpirationTime.Year()),
		SecurityCode:    mqRes.CvvNumber,
	}

	// save card in db
	if card.ID, err = s.db.CreateCard(nil, card); err != nil {
		return nil, errors.Internal(err)
	}
	return card.ToProto(), nil
}

func (s *Service) GetCompanyCards(ctx context.Context, r *pb.GetCompanyCardsRequest) (*pb.Cards, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}

	cc, err := s.db.GetCardsByCompany(nil, companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	var cards []*pb.Card
	for _, c := range cc {
		cards = append(cards, c.ToProto())
	}
	return &pb.Cards{Cards: cards}, nil
}

func (s *Service) GetCard(ctx context.Context, r *pb.GetCardRequest) (*pb.Card, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}
	c, err := s.db.GetCard(nil, r.GetCardId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return c.ToProto(), nil
}

func (s *Service) Freeze(ctx context.Context, r *pb.FreezeRequest) (*pb.Card, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}

	if err := s.db.FreezeCard(nil, r.GetCardId(), companyId); err != nil {
		return nil, errors.Internal(err)
	}

	c, err := s.db.GetCard(nil, r.GetCardId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return c.ToProto(), nil
}

func (s *Service) Unfreeze(ctx context.Context, r *pb.UnfreezeRequest) (*pb.Card, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}

	if err := s.db.UnfreezeCard(nil, r.GetCardId(), companyId); err != nil {
		return nil, errors.Internal(err)
	}

	c, err := s.db.GetCard(nil, r.GetCardId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return c.ToProto(), nil
}

func (s *Service) Close(ctx context.Context, r *pb.CloseRequest) (*pb.Card, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}

	if err := s.db.CloseCard(nil, r.GetCardId(), companyId); err != nil {
		return nil, errors.Internal(err)
	}

	c, err := s.db.GetCard(nil, r.GetCardId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return c.ToProto(), nil
}
