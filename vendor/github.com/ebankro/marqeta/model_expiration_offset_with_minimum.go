/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type ExpirationOffsetWithMinimum struct {
	// specify if a value is provided; default is YEARS
	Unit string `json:"unit,omitempty"`
	// specify if unit is provided; default is 4
	Value     int32      `json:"value,omitempty"`
	MinOffset *MinOffset `json:"min_offset,omitempty"`
}