/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type AvsControls struct {
	AvMessages   *AvsControlOptions `json:"av_messages,omitempty"`
	AuthMessages *AvsControlOptions `json:"auth_messages,omitempty"`
}