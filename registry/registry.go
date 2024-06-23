package registry

// IRegistrar Service registration interface
type IRegistrar interface {
	Register()
	Deregister()
}
