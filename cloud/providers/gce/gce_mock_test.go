package gce

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/evergreen-ci/evergreen/model/host"

	"google.golang.org/api/compute/v1"
)

type clientMock struct {
	// API call options
	failInit   bool
	failCreate bool
	failGet    bool
	failDelete bool

	// Other options
	isActive        bool
	hasAccessConfig bool
}

func (c *clientMock) Init(_ *authOptions) error {
	if c.failInit {
		return errors.New("failed to initialize client")
	}

	return nil
}

// CreateInstance returns a unique identifier for the mock instance.
func (c *clientMock) CreateInstance(_ *host.Host, _ *ProviderSettings) (string, error) {
	if c.failCreate {
		return "", errors.New("failed to create instance")
	}

	return fmt.Sprintf("_%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int()), nil
}

func (c *clientMock) GetInstance(_ string) (*compute.Instance, error) {
	if c.failGet {
		return nil, errors.New("failed to get instance")
	}

	instance := &compute.Instance{Status: "RUNNING"}

	if !c.isActive {
		instance.Status = "STOPPING"
	}

	if c.hasAccessConfig {
		instance.NetworkInterfaces = []*compute.NetworkInterface{&compute.NetworkInterface{
			AccessConfigs: []*compute.AccessConfig{
				&compute.AccessConfig{NatIP: "0.0.0.0"},
			},
		}}
	}

	return instance, nil
}

func (c *clientMock) DeleteInstance(_ string) error {
	if c.failDelete {
		return errors.New("failed to delete instance")
	}

	return nil
}
