package client

import "context"

type MoneyClient interface {
	GetIpDetails(ctx context.Context, ip string) (*IpDetails, error)
}
