package autoscaler

import (
	"context"
	"net"

	"gitlab.com/gitlab-org/fleeting/fleeting/connector"
	"gitlab.com/gitlab-org/fleeting/taskscaler"
	"gitlab.com/gitlab-org/gitlab-runner/executors"
)

var _ executors.Environment = (*acquisitionRef)(nil)

type acquisitionRef struct {
	key string
	acq *taskscaler.Acquisition
}

func (ref *acquisitionRef) ID() string {
	return ref.acq.InstanceID()
}

func (ref *acquisitionRef) OS() string {
	return ref.acq.InstanceConnectInfo().OS
}

func (ref *acquisitionRef) Arch() string {
	return ref.acq.InstanceConnectInfo().Arch
}

func (ref *acquisitionRef) Dial(ctx context.Context) (executors.Client, error) {
	info := ref.acq.InstanceConnectInfo()

	dialer, err := connector.Dial(ctx, info, connector.DialOptions{
		// todo: make this configurable
		UseExternalAddr: true,
	})
	if err != nil {
		return nil, err
	}

	return &client{dialer}, nil
}

func (ref *acquisitionRef) set(key string, acq *taskscaler.Acquisition) {
	if ref.acq != nil {
		panic("acquisition ref already set")
	}

	ref.key = key
	ref.acq = acq
}

func (ref *acquisitionRef) get() string {
	return ref.key
}

type client struct {
	client connector.Client
}

func (c *client) Dial(n string, addr string) (net.Conn, error) {
	return c.client.Dial(n, addr)
}

func (c *client) Run(opts executors.RunOptions) error {
	return c.client.Run(connector.RunOptions(opts))
}

func (c *client) Close() error {
	return c.client.Close()
}
