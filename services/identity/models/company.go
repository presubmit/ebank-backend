package models

import (
	pb "ebank/pb/services/identity"
	"errors"
)

type Company struct {
	ID        string
	Name      string
	CreatedAt string
	CreatedBy string
	IsActive  bool
}

func (c *Company) ValidateFields() error {
	// Validate name
	if len(c.Name) == 0 {
		return errors.New("invalid name")
	}
	return nil
}

func (c *Company) ToProto() *pb.Company {
	return &pb.Company{
		Id:   c.ID,
		Name: c.Name,
	}
}
