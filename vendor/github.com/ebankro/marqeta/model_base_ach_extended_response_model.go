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

type BaseAchExtendedResponseModel struct {
	// yyyy-MM-ddTHH:mm:ssZ
	CreatedTime *time.Time `json:"created_time"`
	// yyyy-MM-ddTHH:mm:ssZ
	LastModifiedTime        *time.Time `json:"last_modified_time"`
	Token                   string     `json:"token"`
	AccountSuffix           string     `json:"account_suffix"`
	VerificationStatus      string     `json:"verification_status,omitempty"`
	AccountType             string     `json:"account_type"`
	NameOnAccount           string     `json:"name_on_account"`
	BankName                string     `json:"bank_name,omitempty"`
	Active                  bool       `json:"active"`
	DateSentForVerification *time.Time `json:"date_sent_for_verification,omitempty"`
	IsDefaultAccount        bool       `json:"is_default_account,omitempty"`
	DateVerified            *time.Time `json:"date_verified,omitempty"`
	VerificationOverride    bool       `json:"verification_override,omitempty"`
	VerificationNotes       string     `json:"verification_notes,omitempty"`
	RoutingNumber           string     `json:"routing_number,omitempty"`
	AccountNumber           string     `json:"account_number,omitempty"`
}
