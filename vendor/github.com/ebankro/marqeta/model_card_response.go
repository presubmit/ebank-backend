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

type CardResponse struct {
	// yyyy-MM-ddTHH:mm:ssZ
	CreatedTime *time.Time `json:"created_time"`
	// yyyy-MM-ddTHH:mm:ssZ
	LastModifiedTime *time.Time `json:"last_modified_time"`
	// 36 char max
	Token string `json:"token"`
	// 36 char max
	UserToken string `json:"user_token"`
	// 36 char max
	CardProductToken string `json:"card_product_token"`
	LastFour         string `json:"last_four"`
	Pan              string `json:"pan"`
	Expiration       string `json:"expiration"`
	// yyyy-MM-ddTHH:mm:ssZ
	ExpirationTime                  *time.Time               `json:"expiration_time"`
	CvvNumber                       string                   `json:"cvv_number,omitempty"`
	ChipCvvNumber                   string                   `json:"chip_cvv_number,omitempty"`
	Barcode                         string                   `json:"barcode"`
	PinIsSet                        bool                     `json:"pin_is_set"`
	State                           string                   `json:"state"`
	StateReason                     string                   `json:"state_reason"`
	FulfillmentStatus               string                   `json:"fulfillment_status"`
	ReissuePanFromCardToken         string                   `json:"reissue_pan_from_card_token,omitempty"`
	Fulfillment                     *CardFulfillmentResponse `json:"fulfillment,omitempty"`
	BulkIssuanceToken               string                   `json:"bulk_issuance_token,omitempty"`
	TranslatePinFromCardToken       string                   `json:"translate_pin_from_card_token,omitempty"`
	ActivationActions               *ActivationActions       `json:"activation_actions,omitempty"`
	InstrumentType                  string                   `json:"instrument_type,omitempty"`
	Expedite                        bool                     `json:"expedite,omitempty"`
	Metadata                        map[string]string        `json:"metadata,omitempty"`
	ContactlessExemptionCounter     int32                    `json:"contactless_exemption_counter,omitempty"`
	ContactlessExemptionTotalAmount float32                  `json:"contactless_exemption_total_amount,omitempty"`
}
