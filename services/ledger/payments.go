package ledger

import (
	"context"
	pb "ebank/pb/services/ledger"
	"ebank/services/ledger/models"
	"ebank/shared/db"
	"ebank/shared/errors"
	"ebank/shared/microservice"

	"github.com/google/uuid"
)

func (s *Service) MakePayment(ctx context.Context, r *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	empId, err := microservice.GetEmployeeId(ctx)
	if err != nil {
		return nil, err
	}
	companyId, err := microservice.GetCompanyId(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range r.GetTransfers() {
		err := s.processTransfer(ctx, t, empId, companyId)
		if err != nil {
			return nil, err
		}
	}

	return &pb.PaymentResponse{}, nil
}

func (s *Service) processTransfer(ctx context.Context, t *pb.Transfer, userId, companyId string) error {
	if t.GetAmount() <= 0 {
		return errors.InvalidArgumentf("invalid amount")
	}

	srcAcc, err := s.GetAccount(ctx, &pb.Account{Id: t.GetFromAccount()})
	if err != nil {
		return err
	}
	if srcAcc.GetCurrency() != t.GetCurrency() {
		return errors.InvalidArgumentf("currency mismatch")
	}
	if srcAcc.GetBalance() < t.GetAmount() {
		return errors.InvalidArgumentf("not enough funds")
	}

	var cpty *pb.Counterparty
	var toAcc *pb.Account
	if len(t.GetCounterpartyId()) != 0 {
		if cpty, err = s.GetCounterparty(ctx, &pb.Counterparty{Id: t.GetCounterpartyId()}); err != nil {
			return err
		}
	}
	if len(t.GetToAccount()) != 0 {
		if toAcc, err = s.GetAccount(ctx, &pb.Account{Id: t.GetToAccount()}); err != nil {
			return err
		}
	}

	// return error if both counterparty and to_account are nil, or both are not nil
	// one, and only one, must be a valid destination
	if (cpty == nil) == (toAcc == nil) {
		return errors.InvalidArgumentf("invalid destination")
	}
	// Transaction id, which should be the same for transaction legs
	tid := uuid.New().String()

	tr := &models.Transaction{
		ID:        tid,
		LegID:     uuid.New().String(),
		AccountID: srcAcc.GetId(),
		CompanyID: companyId,
		CreatedBy: userId,
		Amount:    -t.GetAmount(),
		Currency:  t.GetCurrency(),
		Type:      "INTERNAL",
	}
	if cpty != nil {
		tr.CounterpartyID = cpty.GetId()
	}
	if toAcc != nil {
		tr.OtherAccountID = toAcc.GetId()
	}

	if err := s.db.ExecTx(ctx, func(tx db.Tx) error {
		if err := s.saveTransaction(tx, tr); err != nil {
			return err
		}

		if toAcc != nil {
			if err := s.saveTransaction(tx, &models.Transaction{
				ID:             tid,
				LegID:          uuid.New().String(),
				AccountID:      toAcc.GetId(),
				CompanyID:      companyId,
				CreatedBy:      userId,
				Amount:         t.GetAmount(),
				Currency:       t.GetCurrency(),
				OtherAccountID: srcAcc.GetId(),
				Type:           "INTERNAL",
			}); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// TODO: actually execute transfers from 3rd party provider

	return nil
}

func (s *Service) saveTransaction(tx db.Tx, t *models.Transaction) error {
	// TODO: check if user has access to create transaction from account
	account, err := s.db.AddToAccountBalance(tx, t.Amount, t.AccountID)
	if err != nil {
		return errors.InvalidArgumentf("not enough funds")
	}
	t.AfterBalance = account.Balance
	if _, err := s.db.CreateTransaction(tx, t); err != nil {
		return err
	}
	return nil
}
