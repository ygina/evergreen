package gce

import (
	"os"

	"github.com/evergreen-ci/evergreen/model/host"
	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/compute/v1"
)

const (
	// GoogleAppCredsEnvVar is the environment variable name used for the
	// Google Default Application Credentials JSON file.
	GoogleAppCredsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"
)

type authOptions struct {
	CredsJSON string
	Project   string
	Zone      string
}

// The client interface wraps the Google Compute client interaction.
type client interface {
	Init(*authOptions) error
	CreateInstance(*host.Host, *ProviderSettings) (string, error)
	GetInstance(string) (*compute.Instance, error)
	DeleteInstance(string) error
}

type clientImpl struct {
	Project          string
	Zone             string
	ImagesService    *compute.ImagesService
	InstancesService *compute.InstancesService
}

// Init establishes a connection to a Google Compute endpoint and creates a client that
// can be used to manage instances.
func (c *clientImpl) Init(ao *authOptions) error {
	ctx := context.TODO()

	defer os.Setenv(GoogleAppCredsEnvVar, os.Getenv(GoogleAppCredsEnvVar))

	// The default client constructor reads the contents of the JSON credentials file
	// based on the path set in the google application credentials environment variable.
	os.Setenv(GoogleAppCredsEnvVar, ao.CredsJSON)
	httpClient, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return errors.Wrap(err, "failed to create an oauth http client to google services")
	}

	// Connect to Google Compute Engine.
	service, err := compute.New(httpClient)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Google Compute Engine service")
	}

	c.Project = ao.Project
	c.Zone = ao.Zone

	// Get handles to specific services for image retrieval and instance configuration.
	c.ImagesService = compute.NewImagesService(service)
	c.InstancesService = compute.NewInstancesService(service)

	return nil
}

// CreateInstance requests an instance to be provisioned.
//
// API calls to an instance refer to the instance by the user-provided name (which must be unique)
// and not the ID. If successful, CreateInstance returns the name of the provisioned instance.
func (c *clientImpl) CreateInstance(h *host.Host, s *ProviderSettings) (string, error) {
	// Create instance options to spawn an instance
	instance := &compute.Instance{
		Name:        h.Id,
		Description: "",
		MachineType: makeMachineType(c.Zone, s.MachineType),
		Labels:      makeTags(h),
	}

	// Add the disk with the image URL
	var imageURL string
	if s.ImageFamily != "" {
		imageURL = makeImageFromFamily(s.ImageFamily)
	} else {
		imageURL = makeImage(s.ImageName)
	}

	instance.Disks = []*compute.AttachedDisk{&compute.AttachedDisk{
		AutoDelete: true,
		Boot:       true,
		DeviceName: instance.Name,
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskSizeGb:  s.DiskSizeGB,
			DiskType:    makeDiskType(c.Zone, s.DiskType),
			SourceImage: imageURL,
		},
	}}

	// Attach a network interface
	instance.NetworkInterfaces = []*compute.NetworkInterface{&compute.NetworkInterface{
		AccessConfigs: []*compute.AccessConfig{&compute.AccessConfig{}},
	}}

	// Add the startup script and ssh keys
	keys := sshKeysToString(s.SSHKeys)
	instance.Metadata = &compute.Metadata{
		Items: []*compute.MetadataItems{
			&compute.MetadataItems{Key: "startup-script", Value: &s.SetupScript},
			&compute.MetadataItems{Key: "ssh-keys", Value: &keys},
		},
	}

	// Make the API call to insert the instance
	if _, err := c.InstancesService.Insert(c.Project, c.Zone, instance).Do(); err != nil {
		return "", errors.Wrap(err, "API call to insert instance failed")
	}

	return instance.Name, nil
}

// GetInstance requests details on a single instance.
func (c *clientImpl) GetInstance(name string) (*compute.Instance, error) {
	instance, err := c.InstancesService.Get(c.Project, c.Zone, name).Do()
	if err != nil {
		return nil, errors.Wrap(err, "API call to get instance failed")
	}

	return instance, nil
}

// DeleteInstance requests an instance previously provisioned to be removed.
func (c *clientImpl) DeleteInstance(name string) error {
	if _, err := c.InstancesService.Delete(c.Project, c.Zone, name).Do(); err != nil {
		return errors.Wrap(err, "API call to delete instance failed")
	}

	return nil
}
