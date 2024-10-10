package ledger

import (
	"context"
	"strings"

	pb "ebank/pb/services/ledger"
	m "ebank/services/ledger/models"
	"ebank/shared/errors"
	"ebank/shared/microservice"
)

func (s *Service) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.Account, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}
	// TODO: check if user has access to create company accounts

	acc := &m.Account{
		Name:      strings.TrimSpace(r.GetName()),
		CompanyID: companyId,
		Currency:  r.GetCurrency(),
	}
	if err := acc.ValidateFields(); err != nil {
		return nil, errors.InvalidArgument(err)
	}

	accountId, err := s.db.CreateAccount(nil, acc)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return &pb.Account{
		Id:       accountId,
		Name:     acc.Name,
		Currency: acc.Currency,
	}, nil
}

func (s *Service) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.Accounts, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}
	// TODO: check if user has access read company accounts

	acc, err := s.db.GetAccountsByCompany(nil, companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	var accounts []*pb.Account
	for _, a := range acc {
		accounts = append(accounts, a.ToProto())
	}
	return &pb.Accounts{Accounts: accounts}, nil
}

func (s *Service) GetAccount(ctx context.Context, r *pb.Account) (*pb.Account, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, errors.InvalidArgument(err)
	}
	// TODO: check if user has access read this account

	acc, err := s.db.GetAccount(nil, r.GetId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return acc.ToProto(), nil
}
