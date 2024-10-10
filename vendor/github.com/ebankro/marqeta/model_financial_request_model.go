/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type FinancialRequestModel struct {
	Amount             float32             `json:"amount"`
	CardToken          string              `json:"card_token"`
	Pin                string              `json:"pin,omitempty"`
	Mid                string              `json:"mid"`
	CashBackAmount     float32             `json:"cash_back_amount,omitempty"`
	IsPreAuth          bool                `json:"is_pre_auth,omitempty"`
	CardAcceptor       *CardAcceptorModel  `json:"card_acceptor"`
	TransactionOptions *TransactionOptions `json:"transaction_options,omitempty"`
	Webhook            *Webhook            `json:"webhook,omitempty"`
}
