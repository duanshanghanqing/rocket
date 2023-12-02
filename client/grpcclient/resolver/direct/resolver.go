package direct

import (
	"google.golang.org/grpc/resolver"
)

type directResolver struct{}

func newDirectResolver() *directResolver {
	return new(directResolver)
}

func (r *directResolver) Close() {}

func (r *directResolver) ResolveNow(_ resolver.ResolveNowOptions) {}
