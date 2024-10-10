/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type Msa struct {
	CampaignToken string  `json:"campaign_token"`
	TriggerAmount float32 `json:"trigger_amount"`
	ReloadAmount  float32 `json:"reload_amount"`
}
