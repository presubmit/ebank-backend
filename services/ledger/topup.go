package ledger

import (
	"context"
	pb "ebank/pb/services/ledger"
	"ebank/services/ledger/models"
	"ebank/shared/db"
	"ebank/shared/errors"
	"ebank/shared/microservice"
	"ebank/shared/money"

	"github.com/google/uuid"
)

func (s *Service) TopUp(ctx context.Context, r *pb.TopUpRequest) (*pb.Transaction, error) {
	empId, err := microservice.GetEmployeeId(ctx)
	if err != nil {
		return nil, err
	}
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	if r.GetAmount() <= 0 || !money.IsCurrencyValid(r.GetCurrency()) {
		return nil, errors.InvalidArgumentf("invalid topup")
	}

	acc, err := s.GetAccount(ctx, &pb.Account{Id: r.GetAccountId()})
	if err != nil {
		return nil, err
	}
	if acc.Currency != r.GetCurrency() {
		return nil, errors.InvalidArgumentf("invalid currency")
	}

	t := &models.Transaction{
		ID:        uuid.New().String(),
		LegID:     uuid.New().String(),
		AccountID: acc.GetId(),
		CompanyID: companyId,
		CreatedBy: empId,
		Amount:    r.GetAmount(),
		Currency:  r.GetCurrency(),
		Type:      "TOPUP",
	}

	if err := s.db.ExecTx(ctx, func(tx db.Tx) error {
		// TODO: check if user has access to create transaction from account
		account, err := s.db.AddToAccountBalance(tx, t.Amount, t.AccountID)
		if err != nil {
			// should never happen since we're always adding money
			// but left here just in case
			return errors.InvalidArgumentf("not enough funds")
		}
		t.AfterBalance = account.Balance

		_, err = s.db.CreateTransaction(tx, t)
		return err
	}); err != nil {
		return nil, err
	}

	return t.ToProto(), nil
}
