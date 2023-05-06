package consul

import (
	"context"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	sctx "github.com/viettranx/service-context"
	"sfood/plugins/common"
	"sfood/plugins/discovery"
)

const ttl = time.Second * 5

type consulComponent struct {
	id          string
	serviceName string
	logger      sctx.Logger
	client      *consul.Client
	instanceID  string
	consulHost  string
	version     string
	port        int
}

func NewConsulComponent(id string, serviceName string, version string, port int) *consulComponent {
	return &consulComponent{id: id, serviceName: serviceName, version: version, port: port}
}

func (c *consulComponent) ID() string {
	return c.id
}

func (c *consulComponent) InitFlags() {
	c.consulHost = common.ConsulHost
}

func (c *consulComponent) Activate(sc sctx.ServiceContext) error {
	c.logger = sctx.GlobalLogger().GetLogger(c.id)

	config := consul.DefaultConfig()
	config.Address = c.consulHost
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}

	c.client = client

	c.instanceID = discovery.GenerateInstanceID(c.serviceName)
	if err = c.Register(context.Background()); err != nil {
		return err
	}

	go func() {
		for {
			if err = c.ReportHealthyState(c.instanceID, c.serviceName); err != nil {
				c.logger.Errorln("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return err
}

func (c *consulComponent) Stop() error {
	return c.Deregister(context.Background(), c.instanceID, c.serviceName)
}

func (c *consulComponent) Register(ctx context.Context) error {
	return c.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: discovery.GetOutboundIP().String(),
		ID:      c.instanceID,
		Name:    c.serviceName,
		Port:    50051,
		Tags:    []string{c.version},
		Check: &consul.AgentServiceCheck{
			DeregisterCriticalServiceAfter: ttl.String(),
			TLSSkipVerify:                  true,
			CheckID:                        c.instanceID,
			TTL:                            ttl.String()},
	})
}

func (c *consulComponent) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	return c.client.Agent().ServiceDeregister(instanceID)
}

func (c *consulComponent) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

func (c *consulComponent) ReportHealthyState(instanceID string, serviceName string) error {
	return c.client.Agent().PassTTL(instanceID, "")
}
