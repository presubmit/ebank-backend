/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type JitFundingApi struct {
	Token                                    string                       `json:"token"`
	Method                                   string                       `json:"method"`
	UserToken                                string                       `json:"user_token"`
	ActingUserToken                          string                       `json:"acting_user_token,omitempty"`
	BusinessToken                            string                       `json:"business_token,omitempty"`
	Amount                                   float32                      `json:"amount"`
	Memo                                     string                       `json:"memo,omitempty"`
	Tags                                     string                       `json:"tags,omitempty"`
	OriginalJitFundingToken                  string                       `json:"original_jit_funding_token,omitempty"`
	IncrementalAuthorizationJitFundingTokens []string                     `json:"incremental_authorization_jit_funding_tokens,omitempty"`
	AddressVerification                      *JitAddressVerification      `json:"address_verification,omitempty"`
	DeclineReason                            string                       `json:"decline_reason,omitempty"`
	Balances                                 map[string]CardholderBalance `json:"balances,omitempty"`
}
