/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type KycQuestion struct {
	Key      string   `json:"key,omitempty"`
	Question string   `json:"question,omitempty"`
	Answers  []string `json:"answers,omitempty"`
}