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

type PeerTransferResponse struct {
	Token                  string     `json:"token"`
	Amount                 float32    `json:"amount"`
	Tags                   string     `json:"tags,omitempty"`
	Memo                   string     `json:"memo,omitempty"`
	CurrencyCode           string     `json:"currency_code"`
	SenderUserToken        string     `json:"sender_user_token,omitempty"`
	RecipientUserToken     string     `json:"recipient_user_token,omitempty"`
	SenderBusinessToken    string     `json:"sender_business_token,omitempty"`
	RecipientBusinessToken string     `json:"recipient_business_token,omitempty"`
	CreatedTime            *time.Time `json:"created_time"`
}