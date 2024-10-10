/*
 * Marqeta Core API
 *
 * Simplified management of your payment programs
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package marqeta

type CardHolderAddressModel struct {
	// Required if 'business_token' is not specified
	UserToken string `json:"user_token,omitempty"`
	// Required if 'user_token' is not specified
	BusinessToken string `json:"business_token,omitempty"`
	Token         string `json:"token,omitempty"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Address1      string `json:"address_1"`
	Address2      string `json:"address_2,omitempty"`
	City          string `json:"city"`
	State         string `json:"state"`
	// Required if 'postal_code' is not specified
	Zip              string `json:"zip,omitempty"`
	Country          string `json:"country"`
	Phone            string `json:"phone,omitempty"`
	IsDefaultAddress bool   `json:"is_default_address,omitempty"`
	Active           bool   `json:"active,omitempty"`
	// Required if 'zip' is not specified
	PostalCode string `json:"postal_code,omitempty"`
}