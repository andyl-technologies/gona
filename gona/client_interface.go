package gona

import (
	"context"
)

// ServerOperations defines methods for managing VPS server instances.
type ServerOperations interface {
	CreateServer(ctx context.Context, r *CreateServerRequest) (ServerBuild, error)
	GetServer(ctx context.Context, id int) (Server, error)
	BuildServer(ctx context.Context, id int, r *BuildServerRequest) (ServerBuild, error)
	DeleteServer(ctx context.Context, id int, cancelBilling bool) error
	UnlinkServer(ctx context.Context, id int) error
}

// SSHKeyOperations defines methods for managing SSH keys.
type SSHKeyOperations interface {
	CreateSSHKey(ctx context.Context, name, key string) (SSHKey, error)
	GetSSHKey(ctx context.Context, id int) (SSHKey, error)
	DeleteSSHKey(ctx context.Context, id int) error
}

// BGPSessionOperations defines methods for managing BGP sessions.
type BGPSessionOperations interface {
	CreateBGPSessions(ctx context.Context, mbPkgID, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error)
	GetBGPSessions(ctx context.Context, mbPkgID int) ([]*BGPSession, error)
}

// IPOperations defines methods for managing IP addresses.
type IPOperations interface {
	GetIPs(ctx context.Context, mbPkgID int) (IPs, error)
}

// MetadataOperations defines methods for querying metadata like locations and operating systems.
type MetadataOperations interface {
	GetLocations(ctx context.Context) ([]Location, error)
	GetOSs(ctx context.Context) ([]OS, error)
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
