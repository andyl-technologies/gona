package gonatesting

import (
	"context"
	"fmt"
	"sync"

	"github.com/netactuate/gona/gona"
)

var _ gona.ClientInterface = (*FakeClient)(nil)

// FakeClient implements ClientInterface with in-memory state management.
// Use this for integration-style tests where you need realistic stateful behavior.
//
// Example:
//
//	fake := NewFakeClient()
//	fake.AddServer(Server{ID: 123, Name: "test-server", ServerStatus: "RUNNING"})
//	server, err := fake.GetServer(ctx, 123)  // Returns the added server
type FakeClient struct {
	mu               sync.RWMutex
	servers          map[int]gona.Server
	sshKeys          map[int]gona.SSHKey
	bgpSessions      map[int][]*gona.BGPSession
	ips              map[int]gona.IPs
	locations        []gona.Location
	oses             []gona.OS
	nextServerID     int
	nextSSHKeyID     int
	nextBGPSessionID int

	// Configurable errors - set these to inject errors for specific operations
	CreateServerError      error
	GetServerError         error
	BuildServerError       error
	DeleteServerError      error
	UnlinkServerError      error
	CreateSSHKeyError      error
	GetSSHKeyError         error
	DeleteSSHKeyError      error
	CreateBGPSessionsError error
	GetBGPSessionsError    error
	GetIPsError            error
	GetLocationsError      error
	GetOSsError            error

	// Call tracking
	calltrack struct {
		sync.Mutex
		calls []string
	}
}

// NewFakeClient creates a new FakeClient with default test data.
//
// Default IPs use RFC-reserved documentation ranges that will never
// appear in production:
//   - IPv4: 192.0.2.0/24 (RFC 5737 TEST-NET-1)
//   - IPv6: 2001:db8::/32 (RFC 3849)
func NewFakeClient() *FakeClient {
	return &FakeClient{
		servers:          make(map[int]gona.Server),
		sshKeys:          make(map[int]gona.SSHKey),
		bgpSessions:      make(map[int][]*gona.BGPSession),
		ips:              make(map[int]gona.IPs),
		nextServerID:     1000,
		nextSSHKeyID:     100,
		nextBGPSessionID: 500,
		locations: []gona.Location{
			{ID: 1, Name: "AMS Amsterdam", IATACode: "AMS", Continent: "EU"},
			{ID: 2, Name: "LAX Los Angeles", IATACode: "LAX", Continent: "NA"},
			{ID: 3, Name: "SJC San Jose", IATACode: "SJC", Continent: "NA"},
		},
		oses: []gona.OS{
			{ID: 1, Os: "Ubuntu 22.04 LTS", Type: "linux", Bits: "64"},
			{ID: 2, Os: "Debian 12", Type: "linux", Bits: "64"},
			{ID: 3, Os: "Rocky Linux 9", Type: "linux", Bits: "64"},
		},
	}
}

// CreateServer implements ClientInterface
func (f *FakeClient) CreateServer(ctx context.Context, r *gona.CreateServerRequest) (gona.ServerBuild, error) {
	f.trackCall("CreateServer")
	if f.CreateServerError != nil {
		return gona.ServerBuild{}, f.CreateServerError
	}

	// Validate request
	if r == nil {
		return gona.ServerBuild{}, fmt.Errorf("CreateServerRequest cannot be nil")
	}
	if r.Plan == "" {
		return gona.ServerBuild{}, fmt.Errorf("Plan is required")
	}
	if r.Location <= 0 {
		return gona.ServerBuild{}, fmt.Errorf("invalid Location ID: %d", r.Location)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	serverID := f.nextServerID
	f.nextServerID++

	server := gona.Server{
		ID:                       serverID,
		Name:                     r.FQDN,
		OSID:                     r.Image,
		LocationID:               r.Location,
		PlanID:                   1,
		Package:                  r.Plan,
		PackageBilling:           r.PackageBilling,
		PackageBillingContractId: r.PackageBillingContractId,
		ServerStatus:             "RUNNING",
		PowerStatus:              "on",
		Installed:                1,
		PrimaryIPv4:              fmt.Sprintf("192.0.2.%d", serverID%256),
		PrimaryIPv6:              fmt.Sprintf("2001:db8::%x", serverID),
	}

	// Set OS name from OSID
	for _, os := range f.oses {
		if os.ID == r.Image {
			server.OS = os.Os
			break
		}
	}

	// Set location name from LocationID
	for _, loc := range f.locations {
		if loc.ID == r.Location {
			server.Location = loc.Name
			break
		}
	}

	f.servers[serverID] = server

	return gona.ServerBuild{
		ServerID: serverID,
		Status:   "building",
		Build:    1,
	}, nil
}

// GetServer implements ClientInterface
func (f *FakeClient) GetServer(ctx context.Context, id int) (gona.Server, error) {
	f.trackCall(fmt.Sprintf("GetServer(%d)", id))
	if f.GetServerError != nil {
		return gona.Server{}, f.GetServerError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	server, exists := f.servers[id]
	if !exists {
		return gona.Server{}, fmt.Errorf("server %d not found", id)
	}

	return server, nil
}

// BuildServer implements ClientInterface
func (f *FakeClient) BuildServer(ctx context.Context, id int, r *gona.BuildServerRequest) (gona.ServerBuild, error) {
	f.trackCall(fmt.Sprintf("BuildServer(%d)", id))
	if f.BuildServerError != nil {
		return gona.ServerBuild{}, f.BuildServerError
	}

	// Validate parameters
	if id <= 0 {
		return gona.ServerBuild{}, fmt.Errorf("invalid server ID: %d", id)
	}
	if r == nil {
		return gona.ServerBuild{}, fmt.Errorf("BuildServerRequest cannot be nil")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	server, exists := f.servers[id]
	if !exists {
		return gona.ServerBuild{}, fmt.Errorf("server %d not found", id)
	}

	// Update server with new build parameters
	server.Name = r.FQDN
	server.OSID = r.Image
	server.LocationID = r.Location
	server.PackageBilling = r.PackageBilling
	server.PackageBillingContractId = r.PackageBillingContractId
	server.ServerStatus = "RUNNING"

	// Update OS name
	for _, os := range f.oses {
		if os.ID == r.Image {
			server.OS = os.Os
			break
		}
	}

	// Update location name
	for _, loc := range f.locations {
		if loc.ID == r.Location {
			server.Location = loc.Name
			break
		}
	}

	f.servers[id] = server

	return gona.ServerBuild{
		ServerID: id,
		Status:   "building",
		Build:    1,
	}, nil
}

// DeleteServer implements ClientInterface
func (f *FakeClient) DeleteServer(ctx context.Context, id int, cancelBilling bool) error {
	f.trackCall(fmt.Sprintf("DeleteServer(%d, %v)", id, cancelBilling))
	if f.DeleteServerError != nil {
		return f.DeleteServerError
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.servers[id]; !exists {
		return fmt.Errorf("server %d not found", id)
	}

	// Mark as terminated instead of deleting (more realistic)
	server := f.servers[id]
	server.ServerStatus = "TERMINATED"
	f.servers[id] = server

	return nil
}

// UnlinkServer implements ClientInterface
func (f *FakeClient) UnlinkServer(ctx context.Context, id int) error {
	f.trackCall(fmt.Sprintf("UnlinkServer(%d)", id))
	if f.UnlinkServerError != nil {
		return f.UnlinkServerError
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.servers, id)
	return nil
}

// CreateSSHKey implements ClientInterface
func (f *FakeClient) CreateSSHKey(ctx context.Context, name, key string) (gona.SSHKey, error) {
	f.trackCall(fmt.Sprintf("CreateSSHKey(%s)", name))
	if f.CreateSSHKeyError != nil {
		return gona.SSHKey{}, f.CreateSSHKeyError
	}

	// Validate parameters
	if name == "" {
		return gona.SSHKey{}, fmt.Errorf("SSH key name cannot be empty")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	keyID := f.nextSSHKeyID
	f.nextSSHKeyID++

	sshKey := gona.SSHKey{
		ID:          keyID,
		Name:        name,
		Key:         key,
		Fingerprint: fmt.Sprintf("SHA256:fake-fingerprint-%d", keyID),
	}

	f.sshKeys[keyID] = sshKey
	return sshKey, nil
}

// GetSSHKey implements ClientInterface
func (f *FakeClient) GetSSHKey(ctx context.Context, id int) (gona.SSHKey, error) {
	f.trackCall(fmt.Sprintf("GetSSHKey(%d)", id))
	if f.GetSSHKeyError != nil {
		return gona.SSHKey{}, f.GetSSHKeyError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	key, exists := f.sshKeys[id]
	if !exists {
		return gona.SSHKey{}, fmt.Errorf("ssh key %d not found", id)
	}

	return key, nil
}

// DeleteSSHKey implements ClientInterface
func (f *FakeClient) DeleteSSHKey(ctx context.Context, id int) error {
	f.trackCall(fmt.Sprintf("DeleteSSHKey(%d)", id))
	if f.DeleteSSHKeyError != nil {
		return f.DeleteSSHKeyError
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.sshKeys[id]; !exists {
		return fmt.Errorf("ssh key %d not found", id)
	}

	delete(f.sshKeys, id)
	return nil
}

// CreateBGPSessions implements ClientInterface
func (f *FakeClient) CreateBGPSessions(ctx context.Context, mbPkgID int, groupID int, isIPV6 bool, redundant bool) (*gona.BGPSession, error) {
	f.trackCall(fmt.Sprintf("CreateBGPSessions(%d, %d, %v, %v)", mbPkgID, groupID, isIPV6, redundant))
	if f.CreateBGPSessionsError != nil {
		return nil, f.CreateBGPSessionsError
	}

	// Validate parameters
	if mbPkgID <= 0 {
		return nil, fmt.Errorf("invalid mbPkgID: %d", mbPkgID)
	}
	if groupID <= 0 {
		return nil, fmt.Errorf("invalid groupID: %d", groupID)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	sessionID := f.nextBGPSessionID
	f.nextBGPSessionID++

	ipType := "IPv4"
	customerIP := fmt.Sprintf("198.51.100.%d", sessionID%256)
	if isIPV6 {
		ipType = "IPv6"
		customerIP = fmt.Sprintf("2001:db8:bgp::%x", sessionID)
	}

	session := &gona.BGPSession{
		ID:             sessionID,
		CustomerIP:     customerIP,
		GroupID:        groupID,
		ProviderIPType: ipType,
		ProviderAsn:    64512,
		CustomerAsn:    65000,
		State:          "established",
		ConfigStatus:   1,
	}

	// Inline append to avoid deadlock (AddBGPSession also acquires f.mu)
	sessions, _ := f.bgpSessions[mbPkgID]
	f.bgpSessions[mbPkgID] = append(sessions, session)

	return session, nil
}

// GetBGPSessions implements ClientInterface
func (f *FakeClient) GetBGPSessions(ctx context.Context, mbPkgID int) ([]*gona.BGPSession, error) {
	f.trackCall(fmt.Sprintf("GetBGPSessions(%d)", mbPkgID))
	if f.GetBGPSessionsError != nil {
		return nil, f.GetBGPSessionsError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	sessions, exists := f.bgpSessions[mbPkgID]
	if !exists {
		return []*gona.BGPSession{}, nil
	}

	// Return defensive copy
	result := make([]*gona.BGPSession, len(sessions))
	copy(result, sessions)
	return result, nil
}

// GetIPs implements ClientInterface
func (f *FakeClient) GetIPs(ctx context.Context, mbPkgID int) (gona.IPs, error) {
	f.trackCall(fmt.Sprintf("GetIPs(%d)", mbPkgID))
	if f.GetIPsError != nil {
		return gona.IPs{}, f.GetIPsError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	ips, exists := f.ips[mbPkgID]
	if !exists {
		// Return default IPs if server exists
		if _, serverExists := f.servers[mbPkgID]; serverExists {
			return gona.IPs{
				IPv4: []gona.IP{
					{ID: 1, IP: fmt.Sprintf("192.0.2.%d", mbPkgID%256), Primary: 1, Gateway: "192.0.2.1", Netmask: "255.255.255.0"},
				},
				IPv6: []gona.IP{
					{ID: 2, IP: fmt.Sprintf("2001:db8::%x", mbPkgID), Primary: 1, Gateway: "2001:db8::1", Netmask: "ffff:ffff:ffff:ffff::"},
				},
			}, nil
		}
		return gona.IPs{}, nil
	}

	return ips, nil
}

// GetLocations implements ClientInterface
func (f *FakeClient) GetLocations(ctx context.Context) ([]gona.Location, error) {
	f.trackCall("GetLocations")
	if f.GetLocationsError != nil {
		return nil, f.GetLocationsError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Return defensive copy
	locations := make([]gona.Location, len(f.locations))
	copy(locations, f.locations)
	return locations, nil
}

// GetOSs implements ClientInterface
func (f *FakeClient) GetOSs(ctx context.Context) ([]gona.OS, error) {
	f.trackCall("GetOSs")
	if f.GetOSsError != nil {
		return nil, f.GetOSsError
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Return defensive copy
	oses := make([]gona.OS, len(f.oses))
	copy(oses, f.oses)
	return oses, nil
}

// Helper methods for FakeClient

// AddServer adds a server to the fake client's state
func (f *FakeClient) AddServer(server gona.Server) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.servers[server.ID] = server
}

// AddSSHKey adds an SSH key to the fake client's state
func (f *FakeClient) AddSSHKey(key gona.SSHKey) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sshKeys[key.ID] = key
}

// AddBGPSession adds a BGP session to the fake client's state
func (f *FakeClient) AddBGPSession(mbPkgID int, session *gona.BGPSession) {
	f.mu.Lock()
	defer f.mu.Unlock()
	sessions, _ := f.bgpSessions[mbPkgID]
	f.bgpSessions[mbPkgID] = append(sessions, session)
}

// SetIPs sets the IPs for a server in the fake client's state
func (f *FakeClient) SetIPs(mbPkgID int, ips gona.IPs) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ips[mbPkgID] = ips
}

// AddLocation adds a location to the fake client's state
func (f *FakeClient) AddLocation(location gona.Location) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.locations = append(f.locations, location)
}

// AddOS adds an OS to the fake client's state
func (f *FakeClient) AddOS(os gona.OS) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.oses = append(f.oses, os)
}

// trackCall records method calls for verification
func (f *FakeClient) trackCall(call string) {
	f.calltrack.Lock()
	defer f.calltrack.Unlock()
	f.calltrack.calls = append(f.calltrack.calls, call)
}

// GetCalls returns all tracked method calls
func (f *FakeClient) GetCalls() []string {
	f.calltrack.Lock()
	defer f.calltrack.Unlock()
	calls := make([]string, len(f.calltrack.calls))
	copy(calls, f.calltrack.calls)
	return calls
}

// ResetCalls clears the call tracking
func (f *FakeClient) ResetCalls() {
	f.calltrack.Lock()
	defer f.calltrack.Unlock()
	f.calltrack.calls = nil
}

// Reset clears all state from the FakeClient, returning it to a fresh state.
// This is useful for reusing a FakeClient instance across multiple tests.
func (f *FakeClient) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clear all state maps
	f.servers = make(map[int]gona.Server)
	f.sshKeys = make(map[int]gona.SSHKey)
	f.bgpSessions = make(map[int][]*gona.BGPSession)
	f.ips = make(map[int]gona.IPs)

	// Reset counters to initial values
	f.nextServerID = 1000
	f.nextSSHKeyID = 100
	f.nextBGPSessionID = 500

	// Clear all error fields
	f.CreateServerError = nil
	f.GetServerError = nil
	f.BuildServerError = nil
	f.DeleteServerError = nil
	f.UnlinkServerError = nil
	f.CreateSSHKeyError = nil
	f.GetSSHKeyError = nil
	f.DeleteSSHKeyError = nil
	f.CreateBGPSessionsError = nil
	f.GetBGPSessionsError = nil
	f.GetIPsError = nil
	f.GetLocationsError = nil
	f.GetOSsError = nil

	// Restore default locations and OSes
	f.locations = []gona.Location{
		{ID: 1, Name: "AMS Amsterdam", IATACode: "AMS", Continent: "EU"},
		{ID: 2, Name: "LAX Los Angeles", IATACode: "LAX", Continent: "NA"},
		{ID: 3, Name: "SJC San Jose", IATACode: "SJC", Continent: "NA"},
	}
	f.oses = []gona.OS{
		{ID: 1, Os: "Ubuntu 22.04 LTS", Type: "linux", Bits: "64"},
		{ID: 2, Os: "Debian 12", Type: "linux", Bits: "64"},
		{ID: 3, Os: "Rocky Linux 9", Type: "linux", Bits: "64"},
	}

	// Also reset call tracking
	f.ResetCalls()
}
