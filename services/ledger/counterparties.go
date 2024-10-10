package ledger

import (
	"context"
	pb "ebank/pb/services/ledger"
	m "ebank/services/ledger/models"
	"ebank/shared/errors"
	"ebank/shared/microservice"
)

func (s *Service) GetCounterpartyFields(ctx context.Context, r *pb.Counterparty) (*pb.CounterpartyFields, error) {
	c := m.CounterpartyFromProto(r)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &pb.CounterpartyFields{Fields: c.Fields()}, nil
}

func (s *Service) CreateCounterparty(ctx context.Context, r *pb.Counterparty) (*pb.Counterparty, error) {
	c := m.CounterpartyFromProto(r)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := c.ValidateFields(); err != nil {
		return nil, err
	}

	var err error
	c.CreatedBy, err = microservice.GetEmployeeId(ctx)
	if err != nil {
		return nil, err
	}
	c.CompanyID, err = microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	c.ID, err = s.db.CreateCounterparty(nil, c)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return c.ToProto(), nil
}

func (s *Service) GetCounterparties(ctx context.Context, _ *pb.Empty) (*pb.Counterparties, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.db.GetCounterparties(nil, companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	var cps []*pb.Counterparty
	for _, c := range res {
		cps = append(cps, c.ToProto())
	}
	return &pb.Counterparties{Counterparties: cps}, nil
}

func (s *Service) GetCounterparty(ctx context.Context, r *pb.Counterparty) (*pb.Counterparty, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.db.GetCounterparty(nil, r.GetId(), companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return c.ToProto(), nil
}
