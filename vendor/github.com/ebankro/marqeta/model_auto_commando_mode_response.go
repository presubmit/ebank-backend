/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type AutoCommandoModeResponse struct {
	Response                     *Response                                     `json:"response,omitempty"`
	CommandoModeResponse         *CommandoModeResponse                         `json:"commando_mode_response,omitempty"`
	VelocityControlResponse      *VelocityControlCheckResponse                 `json:"velocity_control_response,omitempty"`
	ProgramFundingSourceResponse *AutoCommandoModeProgramFundingSourceResponse `json:"program_funding_source_response,omitempty"`
}