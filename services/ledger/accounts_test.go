package ledger

import (
	pb "ebank/pb/services/ledger"
	"ebank/services/ledger/db"
	m "ebank/services/ledger/models"
	"ebank/shared/errors"
	"ebank/shared/testutil"
	"testing"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		req              *pb.CreateAccountRequest
		userId           string
		companyId        string
		accountId        string
		createAccountErr error
		res              *pb.Account
		err              error
	}{
		{
			req: &pb.CreateAccountRequest{
				Name:     "Main",
				Currency: "RON",
			},
			res:       &pb.Account{Id: "acc-1", Name: "Main", Currency: "RON"},
			accountId: "acc-1",
			userId:    "user-1",
			companyId: "c-1",
		},
		{
			req: &pb.CreateAccountRequest{
				Name:     " ",
				Currency: "RON",
			},
			err:       errors.InvalidArgumentf("invalid name"),
			userId:    "user-1",
			companyId: "c-1",
		},
		{
			req: &pb.CreateAccountRequest{
				Name: "Main",
			},
			err:       errors.InvalidArgumentf("invalid currency"),
			userId:    "user-1",
			companyId: "c-1",
		},
		{
			req: &pb.CreateAccountRequest{
				Name:     "Main",
				Currency: "RON",
			},
			err:    errors.InvalidArgumentf("invalid company id"),
			userId: "user-1",
		},
		{
			req: &pb.CreateAccountRequest{
				Name:     "Main",
				Currency: "RON",
			},
			createAccountErr: errors.Internal(),
			err:              errors.Internal(),
			userId:           "user-1",
			companyId:        "c-1",
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &db.MockLedgerDB{
				CreateAccountFunc: func(_ *m.Account) (string, error) {
					return test.accountId, test.createAccountErr
				},
			},
		}
		got, err := s.CreateAccount(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("CreateAccount(%v) didn't return error; expected %v; got %v", test.req, test.err, err)
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("CreateAccount(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestGetAccountsByCompany(t *testing.T) {
	tests := []struct {
		req                     *pb.GetAccountsRequest
		userId                  string
		companyId               string
		err                     error
		res                     *pb.Accounts
		getAccountsByCompanyRes []*m.Account
		getAccountsByCompanyErr error
	}{
		{
			req:       &pb.GetAccountsRequest{},
			userId:    "user-1",
			companyId: "c-1",
			getAccountsByCompanyRes: []*m.Account{
				{
					ID:       "acc-1",
					Name:     "Main",
					Balance:  0,
					Currency: "RON",
				},
				{
					ID:       "acc-2",
					Name:     "Other",
					Balance:  500,
					Currency: "EUR",
				},
			},
			res: &pb.Accounts{
				Accounts: []*pb.Account{
					{
						Id:       "acc-1",
						Name:     "Main",
						Balance:  0,
						Currency: "RON",
					},
					{
						Id:       "acc-2",
						Name:     "Other",
						Balance:  500,
						Currency: "EUR",
					},
				},
			},
		},
		{
			req:    &pb.GetAccountsRequest{},
			userId: "user-1",
			err:    errors.InvalidArgumentf("no company id"),
		},
		{
			req:                     &pb.GetAccountsRequest{},
			companyId:               "c-1",
			userId:                  "user-1",
			getAccountsByCompanyErr: errors.Internalf("db error"),
			err:                     errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mctx := &testutil.MockContext{
			UserID:    test.userId,
			CompanyID: test.companyId,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &db.MockLedgerDB{
				GetAccountsByCompanyFunc: func(companyId string) ([]*m.Account, error) {
					if companyId != test.companyId {
						t.Errorf("GetCompanyAccounts(%v); different companyId: expected %v; got %v", test.req, test.companyId, companyId)
					}
					return test.getAccountsByCompanyRes, test.getAccountsByCompanyErr
				},
			},
		}
		got, err := s.GetAccounts(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("GetCompanyAccounts(%v) didn't return error; expected %v; got %v", test.req, test.err, err)
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("GetCompanyAccounts(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}
