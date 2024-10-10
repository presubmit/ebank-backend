package identity

import (
	"context"
	pb "ebank/pb/services/identity"
	"ebank/shared/errors"
)

func (s *Service) VerifyEmployee(ctx context.Context, r *pb.Employee) (*pb.Employee, error) {
	if len(r.GetCompanyId()) == 0 || len(r.GetUserId()) == 0 {
		return nil, errors.InvalidArgumentf("invalid company or user")
	}

	e, err := s.db.GetEmployeeByUserId(nil, r.GetUserId(), r.GetCompanyId())
	if err != nil {
		return nil, errors.NotFoundf("employee not found")
	}

	return e.ToProto(), nil
}
