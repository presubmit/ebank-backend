package models

import (
	pb "ebank/pb/services/ledger"
	"ebank/shared/errors"
	"ebank/shared/utils"
)

type Counterparty struct {
	ID          string
	Country     string
	Currency    string
	FirstName   string
	LastName    string
	CompanyName string
	Type        string
	IBAN        string
	CompanyID   string
	CreatedBy   string
	CreatedAt   string
}

type Country struct {
	Currencies       []string
	IndividualFields []string
	CompanyFields    []string
}

// Country codes use ISO 3166-1 alpha-2 format. (eg. RO)
// Currency codes use ISO 4217 format. (eg. RON)
var countries = map[string]*Country{
	"RO": {
		Currencies: []string{"RON", "EUR"},
		IndividualFields: []string{
			"firstName",
			"lastName",
			"iban",
		},
		CompanyFields: []string{
			"companyName",
			"iban",
		},
	},
}

func (c *Counterparty) Validate() error {
	country := countries[c.Country]
	if country == nil {
		return errors.InvalidArgumentf("invalid country")
	}
	if !utils.ContainsString(country.Currencies, c.Currency) {
		return errors.InvalidArgumentf("invalid currency")
	}
	if c.Type != "individual" && c.Type != "company" {
		return errors.InvalidArgumentf("invalid type")
	}
	return nil
}

func (c *Counterparty) Fields() []string {
	country := countries[c.Country]
	if country == nil {
		return []string{}
	}
	if c.Type == "individual" {
		return country.IndividualFields
	} else if c.Type == "company" {
		return country.CompanyFields
	}
	return []string{}
}

func (c *Counterparty) ValidateFields() error {
	fields := c.Fields()
	if len(fields) == 0 {
		return errors.Internalf("invalid counterparty")
	}
	for _, f := range fields {
		switch f {
		case "firstName":
			if len(c.FirstName) == 0 {
				return errors.InvalidArgumentf("invalid first name")
			}
		case "lastName":
			if len(c.LastName) == 0 {
				return errors.InvalidArgumentf("invalid last name")
			}
		case "companyName":
			if len(c.CompanyName) == 0 {
				return errors.InvalidArgumentf("invalid company name")
			}
		case "iban":
			if len(c.IBAN) != 16 {
				// TODO: Check if iban is valid.
				return errors.InvalidArgumentf("invalid iban")
			}
		}
	}
	return nil
}

func CounterpartyFromProto(r *pb.Counterparty) *Counterparty {
	return &Counterparty{
		Country:     r.GetCountry(),
		Currency:    r.GetCurrency(),
		Type:        r.GetType(),
		IBAN:        r.GetIban(),
		FirstName:   r.GetFirstName(),
		LastName:    r.GetLastName(),
		CompanyName: r.GetCompanyName(),
	}
}

func (c *Counterparty) ToProto() *pb.Counterparty {
	return &pb.Counterparty{
		Id:          c.ID,
		Country:     c.Country,
		Currency:    c.Currency,
		Type:        c.Type,
		Iban:        c.IBAN,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		CompanyName: c.CompanyName,
		CreatedBy:   c.CreatedBy,
	}
}
