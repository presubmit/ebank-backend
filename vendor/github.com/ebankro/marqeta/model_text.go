/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type Text struct {
	NameLine1 *TextValue `json:"name_line_1"`
	NameLine2 *TextValue `json:"name_line_2,omitempty"`
	NameLine3 *TextValue `json:"name_line_3,omitempty"`
}
