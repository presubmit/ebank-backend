/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type AddressVerificationModel struct {
	Request  *AvsInformation `json:"request,omitempty"`
	OnFile   *AvsInformation `json:"on_file,omitempty"`
	Response *Response       `json:"response,omitempty"`
}