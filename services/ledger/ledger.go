package ledger

import (
	"context"

	pb "ebank/pb/services/ledger"
	ldb "ebank/services/ledger/db"
	"ebank/shared/errors"
	"ebank/shared/microservice"
)

type Service struct {
	db ldb.LedgerDB
}

func NewService() *Service {
	return &Service{
		db: ldb.New(),
	}
}

func (s *Service) GetTransactions(ctx context.Context, r *pb.GetTransactionsRequest) (*pb.Transactions, error) {
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.db.GetTransactions(nil, companyId)
	if err != nil {
		return nil, errors.Internal(err)
	}
	var transactions []*pb.Transaction
	for _, t := range res {
		transactions = append(transactions, t.ToProto())
	}
	return &pb.Transactions{Transactions: transactions}, nil
}
