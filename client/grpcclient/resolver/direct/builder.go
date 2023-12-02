package direct

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

const name = "direct"

func init() {
	// refer to dns_resolver
	resolver.Register(NewBuilder())
}

type builder struct {
}

// use direct://<authority>/127.0.0.1:9000

func NewBuilder() resolver.Builder {
	return &builder{}
}

// google.golang.org/grpc/resolver" 290 row
/*
// Builder creates a resolver that will be used to watch name resolution updates.
type Builder interface {
	// Build creates a new resolver for the given target.
	//
	// gRPC dial calls Build synchronously, and fails if the returned error is
	// not nil.
	Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
	// Scheme returns the scheme supported by this resolver.  Scheme is defined
	// at https://github.com/grpc/grpc/blob/master/doc/naming.md.  The returned
	// string should not contain uppercase characters, as they will not match
	// the parsed target's scheme as defined in RFC 3986.
	Scheme() string
}
*/

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// Resolve address
	var addresses = make([]resolver.Address, 0)
	// TrimPrefix use
	/*
		Code:
			var s = "¡¡¡Hello, Gophers!!!"
			s = strings.TrimPrefix(s, "¡¡¡Hello, ")
			s = strings.TrimPrefix(s, "¡¡¡Howdy, ")
			fmt.Print(s)
		Output:
			Gophers!!!
	*/
	trimPrefix := strings.TrimPrefix(target.URL.Path, "/")
	for _, addr := range strings.Split(trimPrefix, ",") {
		addresses = append(addresses, resolver.Address{Addr: addr})
	}

	// grpc The logic for establishing connections is all here
	err := cc.UpdateState(resolver.State{Addresses: addresses})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d *builder) Scheme() string {
	return name
}

// Determine which interface to implement
var _ resolver.Builder = &builder{}
