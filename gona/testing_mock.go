package gona

import (
	"context"
	"fmt"
	"sync"
)

var _ ClientInterface = (*MockClient)(nil)

// MockClient implements ClientInterface with configurable function fields.
// Use this for unit tests where you need precise control over responses.
//
// Example:
//
//	mock := &MockClient{
//	    GetServerFunc: func(ctx context.Context, id int) (Server, error) {
//	        return Server{ID: id, Name: "test-server", ServerStatus: "RUNNING"}, nil
//	    },
//	}
type MockClient struct {
	// Server operation functions
	CreateServerFunc func(ctx context.Context, r *CreateServerRequest) (ServerBuild, error)
	GetServerFunc    func(ctx context.Context, id int) (Server, error)
	BuildServerFunc  func(ctx context.Context, id int, r *BuildServerRequest) (ServerBuild, error)
	DeleteServerFunc func(ctx context.Context, id int, cancelBilling bool) error
	UnlinkServerFunc func(ctx context.Context, id int) error

	// SSH Key operation functions
	CreateSSHKeyFunc func(ctx context.Context, name, key string) (SSHKey, error)
	GetSSHKeyFunc    func(ctx context.Context, id int) (SSHKey, error)
	DeleteSSHKeyFunc func(ctx context.Context, id int) error

	// BGP Session operation functions
	CreateBGPSessionsFunc func(ctx context.Context, mbPkgID int, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error)
	GetBGPSessionsFunc    func(ctx context.Context, mbPkgID int) ([]*BGPSession, error)

	// IP operation functions
	GetIPsFunc func(ctx context.Context, mbPkgID int) (IPs, error)

	// Metadata operation functions
	GetLocationsFunc func(ctx context.Context) ([]Location, error)
	GetOSsFunc       func(ctx context.Context) ([]OS, error)

	// Call tracking
	mu    sync.Mutex
	calls []string
}

// CreateServer implements ClientInterface
func (m *MockClient) CreateServer(ctx context.Context, r *CreateServerRequest) (ServerBuild, error) {
	m.trackCall("CreateServer")
	if m.CreateServerFunc != nil {
		return m.CreateServerFunc(ctx, r)
	}
	return ServerBuild{}, fmt.Errorf("CreateServerFunc not configured")
}

// GetServer implements ClientInterface
func (m *MockClient) GetServer(ctx context.Context, id int) (Server, error) {
	m.trackCall(fmt.Sprintf("GetServer(%d)", id))
	if m.GetServerFunc != nil {
		return m.GetServerFunc(ctx, id)
	}
	return Server{}, fmt.Errorf("GetServerFunc not configured")
}

// BuildServer implements ClientInterface
func (m *MockClient) BuildServer(ctx context.Context, id int, r *BuildServerRequest) (ServerBuild, error) {
	m.trackCall(fmt.Sprintf("BuildServer(%d)", id))
	if m.BuildServerFunc != nil {
		return m.BuildServerFunc(ctx, id, r)
	}
	return ServerBuild{}, fmt.Errorf("BuildServerFunc not configured")
}

// DeleteServer implements ClientInterface
func (m *MockClient) DeleteServer(ctx context.Context, id int, cancelBilling bool) error {
	m.trackCall(fmt.Sprintf("DeleteServer(%d, %v)", id, cancelBilling))
	if m.DeleteServerFunc != nil {
		return m.DeleteServerFunc(ctx, id, cancelBilling)
	}
	return fmt.Errorf("DeleteServerFunc not configured")
}

// UnlinkServer implements ClientInterface
func (m *MockClient) UnlinkServer(ctx context.Context, id int) error {
	m.trackCall(fmt.Sprintf("UnlinkServer(%d)", id))
	if m.UnlinkServerFunc != nil {
		return m.UnlinkServerFunc(ctx, id)
	}
	return fmt.Errorf("UnlinkServerFunc not configured")
}

// CreateSSHKey implements ClientInterface
func (m *MockClient) CreateSSHKey(ctx context.Context, name, key string) (SSHKey, error) {
	m.trackCall(fmt.Sprintf("CreateSSHKey(%s)", name))
	if m.CreateSSHKeyFunc != nil {
		return m.CreateSSHKeyFunc(ctx, name, key)
	}
	return SSHKey{}, fmt.Errorf("CreateSSHKeyFunc not configured")
}

// GetSSHKey implements ClientInterface
func (m *MockClient) GetSSHKey(ctx context.Context, id int) (SSHKey, error) {
	m.trackCall(fmt.Sprintf("GetSSHKey(%d)", id))
	if m.GetSSHKeyFunc != nil {
		return m.GetSSHKeyFunc(ctx, id)
	}
	return SSHKey{}, fmt.Errorf("GetSSHKeyFunc not configured")
}

// DeleteSSHKey implements ClientInterface
func (m *MockClient) DeleteSSHKey(ctx context.Context, id int) error {
	m.trackCall(fmt.Sprintf("DeleteSSHKey(%d)", id))
	if m.DeleteSSHKeyFunc != nil {
		return m.DeleteSSHKeyFunc(ctx, id)
	}
	return fmt.Errorf("DeleteSSHKeyFunc not configured")
}

// CreateBGPSessions implements ClientInterface
func (m *MockClient) CreateBGPSessions(ctx context.Context, mbPkgID int, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error) {
	m.trackCall(fmt.Sprintf("CreateBGPSessions(%d, %d, %v, %v)", mbPkgID, groupID, isIPV6, redundant))
	if m.CreateBGPSessionsFunc != nil {
		return m.CreateBGPSessionsFunc(ctx, mbPkgID, groupID, isIPV6, redundant)
	}
	return nil, fmt.Errorf("CreateBGPSessionsFunc not configured")
}

// GetBGPSessions implements ClientInterface
func (m *MockClient) GetBGPSessions(ctx context.Context, mbPkgID int) ([]*BGPSession, error) {
	m.trackCall(fmt.Sprintf("GetBGPSessions(%d)", mbPkgID))
	if m.GetBGPSessionsFunc != nil {
		return m.GetBGPSessionsFunc(ctx, mbPkgID)
	}
	return nil, fmt.Errorf("GetBGPSessionsFunc not configured")
}

// GetIPs implements ClientInterface
func (m *MockClient) GetIPs(ctx context.Context, mbPkgID int) (IPs, error) {
	m.trackCall(fmt.Sprintf("GetIPs(%d)", mbPkgID))
	if m.GetIPsFunc != nil {
		return m.GetIPsFunc(ctx, mbPkgID)
	}
	return IPs{}, fmt.Errorf("GetIPsFunc not configured")
}

// GetLocations implements ClientInterface
func (m *MockClient) GetLocations(ctx context.Context) ([]Location, error) {
	m.trackCall("GetLocations")
	if m.GetLocationsFunc != nil {
		return m.GetLocationsFunc(ctx)
	}
	return nil, fmt.Errorf("GetLocationsFunc not configured")
}

// GetOSs implements ClientInterface
func (m *MockClient) GetOSs(ctx context.Context) ([]OS, error) {
	m.trackCall("GetOSs")
	if m.GetOSsFunc != nil {
		return m.GetOSsFunc(ctx)
	}
	return nil, fmt.Errorf("GetOSsFunc not configured")
}

// trackCall records method calls for verification
func (m *MockClient) trackCall(call string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, call)
}

// GetCalls returns all tracked method calls
func (m *MockClient) GetCalls() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]string, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// ResetCalls clears the call tracking
func (m *MockClient) ResetCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = nil
}
