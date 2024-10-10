package card

import (
	"context"
	pb "ebank/pb/services/card"
	m "ebank/services/card/models"
	"ebank/shared/db"
	"ebank/shared/errors"
	"ebank/shared/testutil"
	"flag"
	"testing"
	"time"

	mq "github.com/ebankro/marqeta"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type MockDB struct {
	createCardFunc        func() (string, error)
	getCardsByCompanyFunc func(string) ([]*m.Card, error)
	getCardByIdFunc       func(string, string) (*m.Card, error)
	freezeCardFunc        func(string, string) error
	unfreezeCardFunc      func(string, string) error
	closeCardFunc         func(string, string) error
}

func (*MockDB) ExecTx(_ context.Context, f func(db.Tx) error) error {
	return f(nil)
}

func (d *MockDB) CreateCard(db.Tx, *m.Card) (string, error) {
	return d.createCardFunc()
}

func (d *MockDB) GetCardsByCompany(_ db.Tx, companyId string) ([]*m.Card, error) {
	return d.getCardsByCompanyFunc(companyId)
}

func (d *MockDB) GetCard(_ db.Tx, cardId, companyId string) (*m.Card, error) {
	return d.getCardByIdFunc(cardId, companyId)
}

func (d *MockDB) FreezeCard(_ db.Tx, cardId, companyId string) error {
	return d.freezeCardFunc(cardId, companyId)
}

func (d *MockDB) UnfreezeCard(_ db.Tx, cardId, companyId string) error {
	return d.unfreezeCardFunc(cardId, companyId)
}

func (d *MockDB) CloseCard(_ db.Tx, cardId, companyId string) error {
	return d.closeCardFunc(cardId, companyId)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		req        *pb.CreateRequest
		employeeId string
		companyId  string
		cardId     string
		createErr  error
		res        *pb.Card
		err        error
	}{
		{
			req: &pb.CreateRequest{
				CardProductId: "mq-test-1",
				CardholderId:  "emp-1",
			},
			res: &pb.Card{
				Id:    "card-1",
				Brand: "sandbox",
				Type:  "virtual",
			},
			cardId:     "card-1",
			employeeId: "emp-1",
			companyId:  "c-1",
		},
		// invalid card product id
		{
			req: &pb.CreateRequest{
				CardProductId: "invalid-666",
				CardholderId:  "emp-1",
			},
			employeeId: "emp-1",
			companyId:  "c-1",
			err:        errors.InvalidArgument(),
		},
		// db internal error
		{
			req: &pb.CreateRequest{
				CardProductId: "mq-test-1",
				CardholderId:  "emp-1",
			},
			employeeId: "emp-1",
			companyId:  "c-1",
			createErr:  errors.Internalf("db error"),
			err:        errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			CompanyID:  test.companyId,
			EmployeeID: test.employeeId,
		}
		mockMarqeta := &testutil.MockMarqeta{
			CreateCardFunc: func(body mq.CardRequest) (mq.CardResponse, error) {
				if body.UserToken != test.employeeId {
					t.Errorf("Register(%v): marqeta.createCard expected token %s; got %s", test.req, test.req.CardholderId, body.Token)
				}
				now := time.Now()
				return mq.CardResponse{
					UserToken:      test.employeeId,
					Pan:            "1234567812345678",
					ExpirationTime: &now,
					CvvNumber:      "567",
				}, nil
			},
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				createCardFunc: func() (string, error) {
					return test.cardId, test.createErr
				},
			},
			marqeta: mockMarqeta,
		}
		got, err := s.Create(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("Create(%v) expected error %v; got %v", test.req, test.err, err)
		}
		if test.res != nil {
			if got.Id != test.res.GetId() {
				t.Errorf("Create(%v); expected %v; got %v", test.req, test.res.GetId(), got.Id)
			}
			if got.Brand != test.res.GetBrand() {
				t.Errorf("Create(%v); expected %v; got %v", test.req, test.res.GetBrand(), got.Brand)
			}
			if got.Type != test.res.GetType() {
				t.Errorf("Create(%v); expected %v; got %v", test.req, test.res.GetType(), got.Type)
			}
			if len(got.Number) != 16 {
				t.Errorf("Create(%v); expected 16 characters card number; got %v", test.req, got.Type)
			}
			if got.ExpirationMonth < 1 || got.ExpirationMonth > 12 {
				t.Errorf("Create(%v); invalid expiration month; got %v", test.req, got.Type)
			}
			if got.ExpirationYear < 2020 || got.ExpirationYear > 2025 {
				t.Errorf("Create(%v); invalid expiration year; got %v", test.req, got.Type)
			}
			if len(got.SecurityCode) != 3 {
				t.Errorf("Create(%v); invalid security code; got %v", test.req, got.Type)
			}
		} else {
			if !proto.Equal(got, test.res) {
				t.Errorf("Create(%v); expected %v; got %v", test.req, test.res, got)
			}
		}
	}
}

func TestGetCardsByCompany(t *testing.T) {
	flag.Parse()
	tests := []struct {
		req                  *pb.GetCompanyCardsRequest
		employeeId           string
		companyId            string
		err                  error
		res                  *pb.Cards
		getCardsByCompanyRes []*m.Card
		getCardsByCompanyErr error
	}{
		{
			req:        &pb.GetCompanyCardsRequest{},
			employeeId: "emp-1",
			companyId:  "company-1",
			getCardsByCompanyRes: []*m.Card{
				{
					ID:              "card-1",
					Brand:           "mastercard",
					Type:            "physical",
					Number:          "5700123456789012",
					ExpirationMonth: 5,
					ExpirationYear:  21,
					SecurityCode:    "789",
				},
				{
					ID:              "card-2",
					Brand:           "visa",
					Type:            "virtual",
					Number:          "4500123456789012",
					ExpirationMonth: 6,
					ExpirationYear:  22,
					SecurityCode:    "456",
				},
			},
			res: &pb.Cards{
				Cards: []*pb.Card{
					{
						Id:              "card-1",
						Brand:           "mastercard",
						Type:            "physical",
						Number:          "5700123456789012",
						ExpirationMonth: 5,
						ExpirationYear:  21,
						SecurityCode:    "789",
					},
					{
						Id:              "card-2",
						Brand:           "visa",
						Type:            "virtual",
						Number:          "4500123456789012",
						ExpirationMonth: 6,
						ExpirationYear:  22,
						SecurityCode:    "456",
					},
				},
			},
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			CompanyID:  test.companyId,
			EmployeeID: test.employeeId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				getCardsByCompanyFunc: func(companyId string) ([]*m.Card, error) {
					if companyId != test.companyId {
						t.Errorf("GetCompanyCards(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					return test.getCardsByCompanyRes, test.getCardsByCompanyErr
				},
			},
		}
		got, err := s.GetCompanyCards(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("GetCompanyCards(%v) didn't return error; expected %v; got %v", test.req,
				status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("GetCompanyCards(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestGetCard(t *testing.T) {
	tests := []struct {
		req            *pb.GetCardRequest
		userId         string
		companyId      string
		cardId         string
		err            error
		res            *pb.Card
		getCardByIdRes *m.Card
		getCardByIdErr error
	}{
		{
			req: &pb.GetCardRequest{
				CardId: "card-1",
			},
			userId:    "user-1",
			companyId: "company-1",
			cardId:    "card-1",
			getCardByIdRes: &m.Card{
				ID:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
			res: &pb.Card{
				Id:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
		},
		{
			req: &pb.GetCardRequest{
				CardId: "card-1",
			},
			userId:         "user-1",
			companyId:      "company-1",
			cardId:         "card-1",
			getCardByIdErr: errors.Internalf("db error"),
			err:            errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				getCardByIdFunc: func(cardId, companyId string) (*m.Card, error) {
					if companyId != test.companyId {
						t.Errorf("GetCard(%v); different cardId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("GetCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.getCardByIdRes, test.getCardByIdErr
				},
			},
		}
		got, err := s.GetCard(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("GetCard(%v) didn't return error; expected %v; got %v", test.req,
				status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("GetCard(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestFreezeCard(t *testing.T) {
	tests := []struct {
		req            *pb.FreezeRequest
		userId         string
		companyId      string
		cardId         string
		err            error
		res            *pb.Card
		freezeCardErr  error
		getCardByIdRes *m.Card
		getCardByIdErr error
	}{
		{
			req: &pb.FreezeRequest{
				CardId: "card-1",
			},
			userId:    "user-1",
			companyId: "company-1",
			cardId:    "card-1",
			getCardByIdRes: &m.Card{
				ID:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
			res: &pb.Card{
				Id:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
		},
		{
			req: &pb.FreezeRequest{
				CardId: "card-1",
			},
			userId:        "user-1",
			companyId:     "company-1",
			cardId:        "card-1",
			freezeCardErr: errors.Internalf("db error"),
			err:           errors.Internalf("db error"),
		},
		{
			req: &pb.FreezeRequest{
				CardId: "card-1",
			},
			userId:         "user-1",
			companyId:      "company-1",
			cardId:         "card-1",
			getCardByIdErr: errors.Internalf("db error"),
			err:            errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				freezeCardFunc: func(cardId, companyId string) error {
					if companyId != test.companyId {
						t.Errorf("FreezeCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("FreezeCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.freezeCardErr
				},
				getCardByIdFunc: func(cardId, companyId string) (*m.Card, error) {
					if companyId != test.companyId {
						t.Errorf("FreezeCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("FreezeCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.getCardByIdRes, test.getCardByIdErr
				},
			},
		}
		got, err := s.Freeze(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("FreezeCard(%v) didn't return error; expected %v; got %v", test.req,
				status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("FreezeCard(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestUnfreezeCard(t *testing.T) {
	tests := []struct {
		req             *pb.UnfreezeRequest
		userId          string
		companyId       string
		cardId          string
		err             error
		res             *pb.Card
		unfreezeCardErr error
		getCardByIdRes  *m.Card
		getCardByIdErr  error
	}{
		{
			req: &pb.UnfreezeRequest{
				CardId: "card-1",
			},
			userId:    "user-1",
			companyId: "company-1",
			cardId:    "card-1",
			getCardByIdRes: &m.Card{
				ID:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
			res: &pb.Card{
				Id:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
		},
		{
			req: &pb.UnfreezeRequest{
				CardId: "card-1",
			},
			userId:          "user-1",
			companyId:       "company-1",
			cardId:          "card-1",
			unfreezeCardErr: errors.Internalf("db error"),
			err:             errors.Internalf("db error"),
		},
		{
			req: &pb.UnfreezeRequest{
				CardId: "card-1",
			},
			userId:         "user-1",
			companyId:      "company-1",
			cardId:         "card-1",
			getCardByIdErr: errors.Internalf("db error"),
			err:            errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				unfreezeCardFunc: func(cardId, companyId string) error {
					if companyId != test.companyId {
						t.Errorf("UnfreezeCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("UnfreezeCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.unfreezeCardErr
				},
				getCardByIdFunc: func(cardId, companyId string) (*m.Card, error) {
					if companyId != test.companyId {
						t.Errorf("UnfreezeCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("UnfreezeCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.getCardByIdRes, test.getCardByIdErr
				},
			},
		}
		got, err := s.Unfreeze(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("UnfreezeCard(%v) didn't return error; expected %v; got %v", test.req,
				status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("UnfreezeCard(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestCloseCard(t *testing.T) {
	tests := []struct {
		req            *pb.CloseRequest
		userId         string
		companyId      string
		cardId         string
		err            error
		res            *pb.Card
		closeCardErr   error
		getCardByIdRes *m.Card
		getCardByIdErr error
	}{
		{
			req: &pb.CloseRequest{
				CardId: "card-1",
			},
			userId:    "user-1",
			companyId: "company-1",
			cardId:    "card-1",
			getCardByIdRes: &m.Card{
				ID:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
			res: &pb.Card{
				Id:              "card-1",
				Brand:           "mastercard",
				Type:            "physical",
				Number:          "5700123456789012",
				ExpirationMonth: 5,
				ExpirationYear:  21,
				SecurityCode:    "789",
			},
		},
		{
			req: &pb.CloseRequest{
				CardId: "card-1",
			},
			userId:       "user-1",
			companyId:    "company-1",
			cardId:       "card-1",
			closeCardErr: errors.Internalf("db error"),
			err:          errors.Internalf("db error"),
		},
		{
			req: &pb.CloseRequest{
				CardId: "card-1",
			},
			userId:         "user-1",
			companyId:      "company-1",
			cardId:         "card-1",
			getCardByIdErr: errors.Internalf("db error"),
			err:            errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &MockDB{
				closeCardFunc: func(cardId, companyId string) error {
					if companyId != test.companyId {
						t.Errorf("CloseCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("CloseCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.closeCardErr
				},
				getCardByIdFunc: func(cardId, companyId string) (*m.Card, error) {
					if companyId != test.companyId {
						t.Errorf("CloseCard(%v); different companyId: expected %v; got %v",
							test.req, test.companyId, companyId)
					}
					if cardId != test.cardId {
						t.Errorf("CloseCard(%v); different cardId: expected %v; got %v",
							test.req, test.cardId, cardId)
					}
					return test.getCardByIdRes, test.getCardByIdErr
				},
			},
		}
		got, err := s.Close(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("CloseCard(%v) didn't return error; expected %v; got %v", test.req,
				status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("CloseCard(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}
