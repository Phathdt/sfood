package clerkc

import (
	"flag"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	sctx "github.com/viettranx/service-context"
)

type clerkComponent struct {
	id     string
	name   string
	logger sctx.Logger
	client clerk.Client
	token  string
}

func NewClerkComponent(id string) *clerkComponent {
	return &clerkComponent{id: id}
}

func (c *clerkComponent) ID() string {
	return c.id
}

func (c *clerkComponent) InitFlags() {
	flag.StringVar(&c.token, "clerk_token", "", "clerk token")
}

func (c *clerkComponent) Activate(sc sctx.ServiceContext) error {
	c.logger = sc.Logger(c.id)
	c.name = sc.GetName()

	client, err := clerk.NewClient(c.token)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *clerkComponent) Stop() error {
	return nil
}

func (c *clerkComponent) GetClient() clerk.Client {
	return c.client
}
