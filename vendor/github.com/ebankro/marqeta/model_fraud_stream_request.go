/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

import (
	"time"
)

type FraudStreamRequest struct {
	Program                        string                              `json:"program,omitempty"`
	Type_                          string                              `json:"type,omitempty"`
	State                          string                              `json:"state,omitempty"`
	Itc                            string                              `json:"itc,omitempty"`
	Token                          string                              `json:"token,omitempty"`
	UserToken                      string                              `json:"user_token,omitempty"`
	ActingUserToken                string                              `json:"acting_user_token,omitempty"`
	CardToken                      string                              `json:"card_token,omitempty"`
	UserTransactionTime            *time.Time                          `json:"user_transaction_time,omitempty"`
	RequestAmount                  float32                             `json:"request_amount,omitempty"`
	Amount                         float32                             `json:"amount,omitempty"`
	CurrencyCode                   string                              `json:"currency_code,omitempty"`
	Network                        string                              `json:"network,omitempty"`
	AccountRiskScore               string                              `json:"account_risk_score,omitempty"`
	AccountRiskScoreReasonCode     string                              `json:"account_risk_score_reason_code,omitempty"`
	TransactionRiskScore           int32                               `json:"transaction_risk_score,omitempty"`
	TransactionRiskScoreReasonCode string                              `json:"transaction_risk_score_reason_code,omitempty"`
	CardAcceptor                   *TransactionCardAcceptorViewModelV1 `json:"card_acceptor,omitempty"`
	AddressVerification            *AddressVerificationModel           `json:"address_verification,omitempty"`
}