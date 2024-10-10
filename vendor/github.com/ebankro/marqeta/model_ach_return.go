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

type AchReturn struct {
	Amount        float32    `json:"amount,omitempty"`
	Date          *time.Time `json:"date,omitempty"`
	DateInitiated *time.Time `json:"dateInitiated,omitempty"`
	OrderId       string     `json:"orderId,omitempty"`
	ReasonCode    string     `json:"reasonCode,omitempty"`
	DirectDeposit bool       `json:"directDeposit,omitempty"`
	AchType       string     `json:"achType,omitempty"`
}