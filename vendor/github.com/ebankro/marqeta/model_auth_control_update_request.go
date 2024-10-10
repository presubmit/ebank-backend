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

type AuthControlUpdateRequest struct {
	Token         string                   `json:"token"`
	Name          string                   `json:"name,omitempty"`
	Association   *SpendControlAssociation `json:"association,omitempty"`
	MerchantScope *MerchantScope           `json:"merchant_scope,omitempty"`
	StartTime     *time.Time               `json:"start_time,omitempty"`
	EndTime       *time.Time               `json:"end_time,omitempty"`
	Active        bool                     `json:"active,omitempty"`
}
