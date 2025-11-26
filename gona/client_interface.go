package gona

import (
	"context"
)

// ServerOperations defines methods for managing VPS server instances.
type ServerOperations interface {
	CreateServer(context.Context, *CreateServerRequest) (ServerBuild, error)
	GetServer(_ context.Context, id int) (Server, error)
	BuildServer(_ context.Context, id int, r *BuildServerRequest) (ServerBuild, error)
	DeleteServer(_ context.Context, id int, cancelBilling bool) error
	UnlinkServer(_ context.Context, id int) error
}

// SSHKeyOperations defines methods for managing SSH keys.
type SSHKeyOperations interface {
	CreateSSHKey(_ context.Context, name, key string) (SSHKey, error)
	GetSSHKey(_ context.Context, id int) (SSHKey, error)
	DeleteSSHKey(_ context.Context, id int) error
}

// BGPSessionOperations defines methods for managing BGP sessions.
type BGPSessionOperations interface {
	CreateBGPSessions(_ context.Context, mbPkgID, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error)
	GetBGPSessions(_ context.Context, mbPkgID int) ([]*BGPSession, error)
}

// IPOperations defines methods for managing IP addresses.
type IPOperations interface {
	GetIPs(_ context.Context, mbPkgID int) (IPs, error)
}

// MetadataOperations defines methods for querying metadata like locations and operating systems.
type MetadataOperations interface {
	GetLocations(context.Context) ([]Location, error)
	GetOSs(context.Context) ([]OS, error)
}

// ClientInterface is composed of all operation-specific interfaces.
// This enables mocking and faking for tests, and allows consumers to depend
// on only the subset of operations they need.
type ClientInterface interface {
	ServerOperations
	SSHKeyOperations
	BGPSessionOperations
	IPOperations
	MetadataOperations
}

var _ ClientInterface = (*Client)(nil)
