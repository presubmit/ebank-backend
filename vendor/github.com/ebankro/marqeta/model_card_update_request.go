/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type CardUpdateRequest struct {
	Token       string            `json:"token"`
	UserToken   string            `json:"user_token,omitempty"`
	Expedite    bool              `json:"expedite,omitempty"`
	Fulfillment *Fulfillment      `json:"fulfillment,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}