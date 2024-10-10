package identity

import (
	apb "ebank/pb/services/auth"
	pb "ebank/pb/services/identity"
	"ebank/services/identity/db"
	m "ebank/services/identity/models"
	"ebank/shared/errors"
	"ebank/shared/testutil"
	"testing"

	mq "github.com/ebankro/marqeta"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestCreateCompany(t *testing.T) {
	tests := []struct {
		req               *pb.CreateCompanyRequest
		companyID         string
		userID            string
		employeeID        string
		createCompanyErr  error
		createCompanyRes  string
		createUserRoleErr error
		res               *pb.Company
		err               error
	}{
		{
			req: &pb.CreateCompanyRequest{
				Name: "Tesla",
			},
			companyID: "test-1",
			res: &pb.Company{
				Id:   "test-1",
				Name: "Tesla",
			},
			userID: "user-1",
		},
		{
			req: &pb.CreateCompanyRequest{
				Name: "Tesla",
			},
			createCompanyErr: errors.Internalf("db error"),
			err:              errors.Internalf("db error"),
			userID:           "user-1",
		},
		{
			req: &pb.CreateCompanyRequest{
				Name: "",
			},
			err:    errors.InvalidArgument(),
			userID: "user-1",
		},
		{
			req: &pb.CreateCompanyRequest{
				Name: "Tesla",
			},
			companyID:         "test-1",
			userID:            "user-1",
			createUserRoleErr: errors.Internalf("db error"),
			err:               errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mockAuthService := &apb.MockAuthServiceClient{}
		mockMarqeta := &testutil.MockMarqeta{
			CreateBusinessFunc: func(body mq.BusinessCardholder) (mq.BusinessCardHolderResponse, error) {
				if body.Token != test.companyID {
					t.Errorf("Register(%v): marqeta.createBusiness expected token %s; got %s", test.req, test.companyID, body.Token)
				}
				return mq.BusinessCardHolderResponse{
					Token: test.companyID,
				}, nil
			},
			CreateUserFunc: func(body mq.CardHolderModel) (mq.UserCardHolderResponse, error) {
				if body.Token != test.employeeID {
					t.Errorf("Register(%v): marqeta.createUser expected token %s; got %s", test.req, test.userID, body.Token)
				}
				return mq.UserCardHolderResponse{
					Token: test.employeeID,
				}, nil
			},
		}
		mctx := &testutil.MockContext{
			UserID: test.userID,
		}
		ctx := mctx.Build()
		s := &Service{
			db: &db.MockIdentityDB{
				GetUserFunc: func(userId string) (*m.User, error) {
					if userId != test.userID {
						t.Errorf("CreateCompany(%v).GetUser; expected userId %v; got %v", test.req, test.userID, userId)
					}
					return &m.User{ID: userId}, nil
				},
				CreateCompanyFunc: func(_ *m.Company) (string, error) {
					return test.companyID, test.createCompanyErr
				},
				CreateEmployeeFunc: func(e *m.Employee) (string, error) {
					if e.UserId != test.userID {
						t.Errorf("CreateCompany(%v); different userId: expected %v; got %v", e, test.userID, e.UserId)
					}
					if e.CompanyId != test.companyID {
						t.Errorf("CreateCompany(%v); different companyId: expected %v; got %v", e, test.companyID, e.CompanyId)
					}
					if e.Role != "OWNER" {
						t.Errorf("CreateCompany(%v); different companyId: expected OWNER; got %v", e, e.Role)
					}
					return test.employeeID, test.createUserRoleErr
				},
			},
			authService: mockAuthService,
			marqeta:     mockMarqeta,
		}
		got, err := s.CreateCompany(ctx, test.req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("Register(%v) expected error %v; got %v", test.req, test.err, err)
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("Register(%v); expected %v; got %v", test.req, test.res, got)
		}
	}
}

func TestGetCompanies(t *testing.T) {
	tests := []struct {
		userID          string
		err             error
		res             *pb.Companies
		getCompaniesRes []*m.Company
		getCompaniesErr error
	}{
		{
			userID: "user-1",
			getCompaniesRes: []*m.Company{
				{
					ID:   "company-1",
					Name: "tesla",
				},
				{
					ID:   "company-2",
					Name: "spacex",
				},
			},
			res: &pb.Companies{
				Companies: []*pb.Company{
					{
						Id:   "company-1",
						Name: "tesla",
					},
					{
						Id:   "company-2",
						Name: "spacex",
					},
				},
			},
		},
		{
			userID:          "user-1",
			getCompaniesErr: errors.Internalf("db error"),
			err:             errors.Internalf("db error"),
		},
	}

	for _, test := range tests {
		mockAuthService := &apb.MockAuthServiceClient{}
		mctx := &testutil.MockContext{
			UserID: test.userID,
		}
		req := &pb.Empty{}
		ctx := mctx.Build()
		s := &Service{
			db: &db.MockIdentityDB{
				GetCompaniesByUserIdFunc: func(userId string) ([]*m.Company, error) {
					if userId != test.userID {
						t.Errorf("GetCompanies(%v); different userId: expected %v; got %v", req, test.userID, userId)
					}
					return test.getCompaniesRes, test.getCompaniesErr
				},
			},
			authService: mockAuthService,
		}
		got, err := s.GetCompanies(ctx, req)
		if status.Code(err) != status.Code(test.err) {
			t.Errorf("GetCompanies(%v) didn't return error; expected %v; got %v", req, status.Code(test.err), status.Code(err))
		}
		if !proto.Equal(got, test.res) {
			t.Errorf("GetCompanies(%v); expected %v; got %v", req, test.res, got)
		}
	}
}
