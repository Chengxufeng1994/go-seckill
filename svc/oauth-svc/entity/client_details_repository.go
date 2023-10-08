package entity

import "context"

type ClientDetailsRepository interface {
	GetClientDetailsByClientId(context.Context, string) (*ClientDetails, error)
	CreateClientDetails(context.Context, *ClientDetails) (uint, error)
}
