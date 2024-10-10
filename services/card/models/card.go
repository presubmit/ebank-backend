package models

import (
	pb "ebank/pb/services/card"
	"math/rand"
	"strconv"
)

type Card struct {
	ID              string
	CompanyId       string
	EmployeeId      string
	Brand           string
	Number          string
	ExpirationMonth int32
	ExpirationYear  int32
	SecurityCode    string
	Type            string
	FrozenAt        string
	ClosedAt        string
	CreatedAt       string
	CreatedBy       string
}

func (c *Card) GenerateRandomValues() {
	c.Number = "54"
	for i := 0; i < 14; i++ {
		c.Number += strconv.Itoa(rand.Intn(9))
	}
	c.ExpirationMonth = 1 + int32(rand.Intn(11))
	c.ExpirationYear = 2020 + int32(rand.Intn(4))
	c.SecurityCode = strconv.Itoa(rand.Intn(999))
}

func (c *Card) ToProto() *pb.Card {
	return &pb.Card{
		Id:              c.ID,
		Brand:           c.Brand,
		Number:          c.Number,
		ExpirationMonth: c.ExpirationMonth,
		ExpirationYear:  c.ExpirationYear,
		SecurityCode:    c.SecurityCode,
		Type:            c.Type,
		FrozenAt:        c.FrozenAt,
	}
}
