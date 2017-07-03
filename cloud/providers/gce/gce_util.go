package gce

import (
	"fmt"
	"strings"

	"github.com/evergreen-ci/evergreen/cloud"
	"github.com/evergreen-ci/evergreen/model/host"
)

const (
	// NameTimeFormat is the format in which to log times like start time.
	NameTimeFormat = "20060102150405"
	// ImageReadyStatus indicates that an image is ready to be used.
	ImageReadyStatus = "READY"
	// StatusProvisioning means resources are being reserved for the instance.
	StatusProvisioning = "PROVISIONING"
	// StatusStaging means resources have been acquired and the instance is
	// preparing for launch.
	StatusStaging = "STAGING"
	// StatusRunning means the instance is booting up or running. You should be
	// able to SSH into the instance soon, though not immediately, after it
	// enters this state.
	StatusRunning = "RUNNING"
	// StatusStopping means the instance is being stopped either due to a
	// failure, or the instance is being shut down. This is a temporary status
	// before it becomes terminated.
	StatusStopping = "STOPPING"
	// StatusTerminated means the instance was shut down or encountered a
	// failure, either through the API or from inside the guest.
	StatusTerminated = "TERMINATED"
)

// SSHKey is the ssh key information to add to a VM instance.
type sshKey struct {
	Username  string `mapstructure:"username"`
	PublicKey string `mapstructure:"public_key"`
}

// Formats the ssh key in "username:publickey" format.
func (key sshKey) String() string {
	return fmt.Sprintf("%s:%s", key.Username, key.PublicKey)
}

// Formats the ssh keys in "username:publickey" format joined by newlines.
func sshKeysToString(keys []sshKey) string {
	arr := make([]string, len(keys))
	for i, key := range keys {
		arr[i] = key.String()
	}
	return strings.Join(arr, "\n")
}

func toEvgStatus(status string) cloud.CloudStatus {
	switch status {
	case StatusProvisioning:
		return cloud.StatusInitializing
	case StatusStaging:
		return cloud.StatusInitializing
	case StatusRunning:
		return cloud.StatusRunning
	case StatusStopping:
		return cloud.StatusStopped
	case StatusTerminated:
		return cloud.StatusTerminated
	default:
		return cloud.StatusUnknown
	}
}

// Returns a machine type URL for the given the zone
func makeMachineType(zone, machineType string) string {
	return fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
}

// Returns a disk type URL given the zone
func makeDiskType(zone, disk string) string {
	return fmt.Sprintf("zones/%s/diskTypes/%s", zone, disk)
}

// Returns an image source URL for a private image family. The URL refers to
// the newest image version associated with the given family.
func makeImageFromFamily(family string) string {
	return fmt.Sprintf("global/images/family/%s", family)
}

// Returns an image source URL for a private image.
func makeImage(name string) string {
	return fmt.Sprintf("global/images/%s", name)
}

func makeTags(intent *host.Host) map[string]string {
	tags := map[string]string{
		"distro":     intent.Distro.Id,
		"owner":      intent.StartedBy,
		"mode":       "production",
		"start-time": intent.CreationTime.Format(NameTimeFormat),
	}

	if intent.UserHost {
		tags["mode"] = "testing"
	}
	return tags
}
