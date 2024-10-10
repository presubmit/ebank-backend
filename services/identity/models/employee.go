package models

import pb "ebank/pb/services/identity"

type Employee struct {
	ID             string
	UserId         string
	CompanyId      string
	Email          string
	Role           string
	CreatedAt      string
	InvitationSent bool
	IsActive       bool
}

func (e *Employee) ToProto() *pb.Employee {
	return &pb.Employee{
		Id:        e.ID,
		UserId:    e.UserId,
		CompanyId: e.CompanyId,
		Email:     e.Email,
		Role:      e.Role,
	}
}
