package powerdns

import (
	"context"
)

type Client interface {
	Healthcheck(ctx context.Context) (Server, error)
	GetServers(ctx context.Context) ([]Server, error)
}
