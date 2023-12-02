package registry

import "context"

// Service registration interface
type IRegistrar interface {
	Register(ctx context.Context, service *ServiceRegisterInfo) error
	Deregister(ctx context.Context, service *ServiceRegisterInfo) error
}

type ServiceRegisterInfo struct {
	ID   string
	Name string
	Tags []string

	// Used during service registration, it is a Host:Port
	Host string
	Port int

	SocketPath string
	Meta       map[string]string
	Namespace  string
	Version    string
	// Service Metadata
	Metadata map[string]string `json:"metadata"`
}
