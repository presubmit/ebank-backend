/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

import (
	"time"
)

type CommandoModeResponse struct {
	Token                            string                        `json:"token,omitempty"`
	ProgramGatewayFundingSourceToken string                        `json:"program_gateway_funding_source_token,omitempty"`
	CurrentState                     *CommandoModeNestedTransition `json:"current_state,omitempty"`
	CommandoModeEnables              *CommandoModeEnables          `json:"commando_mode_enables,omitempty"`
	RealTimeStandinCriteria          *RealTimeStandinCriteria      `json:"real_time_standin_criteria,omitempty"`
	// yyyy-MM-ddTHH:mm:ssZ
	CreatedTime *time.Time `json:"created_time"`
	// yyyy-MM-ddTHH:mm:ssZ
	LastModifiedTime *time.Time `json:"last_modified_time"`
}