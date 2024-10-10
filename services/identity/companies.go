package identity

import (
	"context"
	pb "ebank/pb/services/identity"
	m "ebank/services/identity/models"
	"ebank/shared/db"
	"ebank/shared/errors"
	"ebank/shared/microservice"
	"strings"

	mq "github.com/ebankro/marqeta"
)

func (s *Service) CreateCompany(ctx context.Context, r *pb.CreateCompanyRequest) (*pb.Company, error) {
	userId, err := microservice.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	// company
	company := &m.Company{
		Name:      strings.TrimSpace(r.GetName()),
		CreatedBy: userId,
	}
	if err := company.ValidateFields(); err != nil {
		return nil, errors.InvalidArgument(err)
	}

	user, err := s.db.GetUser(nil, userId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	var employeeId string
	if err := s.db.ExecTx(ctx, func(tx db.Tx) error {
		var err error
		// save company in db
		company.ID, err = s.db.CreateCompany(tx, company)
		if err != nil {
			return errors.Internal(err)
		}
		// save employee in db
		employeeId, err = s.db.CreateEmployee(tx, &m.Employee{
			UserId:    user.ID,
			CompanyId: company.ID,
			Email:     user.Email,
			Role:      "OWNER",
		})
		if err != nil {
			return errors.Internal(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Marqeta
	if _, err := s.marqeta.CreateBusiness(mq.BusinessCardholder{
		Token: company.ID,
	}); err != nil {
		return nil, errors.Internal(err)
	}
	if _, err := s.marqeta.CreateUser(mq.CardHolderModel{
		Token: employeeId,
	}); err != nil {
		return nil, errors.Internal(err)
	}

	return company.ToProto(), nil
}

func (s *Service) GetCompanies(ctx context.Context, r *pb.Empty) (*pb.Companies, error) {
	userId, err := microservice.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	companies, err := s.db.GetCompaniesByUserId(nil, userId)
	if err != nil {
		return nil, errors.Internal(err)
	}

	var companiesRes []*pb.Company
	for _, c := range companies {
		companiesRes = append(companiesRes, &pb.Company{
			Id:   c.ID,
			Name: c.Name,
		})
	}
	return &pb.Companies{Companies: companiesRes}, nil
}
