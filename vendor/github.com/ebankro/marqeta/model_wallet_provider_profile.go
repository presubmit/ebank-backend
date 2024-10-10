/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type WalletProviderProfile struct {
	Account               *Account        `json:"account,omitempty"`
	RiskAssessment        *RiskAssessment `json:"risk_assessment,omitempty"`
	DeviceScore           string          `json:"device_score,omitempty"`
	PanSource             string          `json:"pan_source,omitempty"`
	ReasonCode            string          `json:"reason_code,omitempty"`
	RecommendationReasons []string        `json:"recommendation_reasons,omitempty"`
}
