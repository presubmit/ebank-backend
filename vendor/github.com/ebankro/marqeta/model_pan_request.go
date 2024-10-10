/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type PanRequest struct {
	Pan        string `json:"pan"`
	CvvNumber  string `json:"cvv_number,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}
