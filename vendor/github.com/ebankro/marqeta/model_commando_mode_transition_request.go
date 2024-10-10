/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type CommandoModeTransitionRequest struct {
	Token             string                        `json:"token,omitempty"`
	CommandoModeToken string                        `json:"commando_mode_token"`
	Transition        *CommandoModeNestedTransition `json:"transition"`
}
