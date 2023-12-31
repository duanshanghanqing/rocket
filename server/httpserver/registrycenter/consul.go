package registrycenter

import (
	"context"
	"errors"
	"fmt"
	"github.com/duanshanghanqing/rocket/registry"
	consulApi "github.com/hashicorp/consul/api"
)

type ConsulRegisterCenter struct {
	client       *consulApi.Client
	clientConfig *consulApi.Config
}

type ConsulRegisterCenterOptions func(o *ConsulRegisterCenter)

func WithConsulRegisterCenterClientConfig(clientConfig *consulApi.Config) ConsulRegisterCenterOptions {
	return func(o *ConsulRegisterCenter) {
		o.clientConfig = clientConfig
	}
}

func NewConsulRegisterCenter(opts ...ConsulRegisterCenterOptions) (registry.IRegistrar, error) {
	r := &ConsulRegisterCenter{
		clientConfig: consulApi.DefaultConfig(),
	}

	for _, o := range opts {
		o(r)
	}

	client, err := consulApi.NewClient(r.clientConfig)
	if err != nil {
		return nil, err
	}

	r.client = client
	return r, err
}

func (r *ConsulRegisterCenter) Register(ctx context.Context, service *registry.ServiceRegisterInfo) error {
	if service == nil {
		return errors.New("ServiceRegisterInfo is nil")
	}
	// 1.Registration Service Information
	reg := consulApi.AgentServiceRegistration{}
	reg.ID = service.ID
	reg.Name = service.Name
	reg.Address = service.Host
	reg.Port = service.Port
	reg.Tags = service.Tags

	// 2.Set up health check ups
	check := &consulApi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", reg.Address, reg.Port),
		Timeout:                        "2s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s",
	}

	// 3.Add check
	reg.Check = check
	fmt.Println(reg.Check.HTTP)

	// 4.Registration Services
	err := r.client.Agent().ServiceRegister(&reg)
	if err != nil {
		return err
	}

	return nil
}

func (r *ConsulRegisterCenter) Deregister(ctx context.Context, service *registry.ServiceRegisterInfo) error {
	if service == nil {
		return errors.New("ServiceRegisterInfo is nil")
	}
	return r.client.Agent().ServiceDeregister(service.ID) // Unregister service
}
