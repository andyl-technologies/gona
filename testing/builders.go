package gonatesting

import (
	"fmt"

	"github.com/netactuate/gona/gona"
)

// ServerBuilder provides a fluent API for creating test Server instances.
//
// Example:
//
//	server := NewServerBuilder().
//	    WithID(123).
//	    WithName("test.example.com").
//	    WithStatus("RUNNING").
//	    Build()
type ServerBuilder struct {
	server gona.Server
}

// NewServerBuilder creates a new ServerBuilder with sensible defaults
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		server: gona.Server{
			ID:                       1000,
			Name:                     "test-server.example.com",
			OS:                       "Ubuntu 22.04 LTS",
			OSID:                     1,
			PrimaryIPv4:              "192.0.2.100",
			PrimaryIPv6:              "2001:db8::100",
			PlanID:                   1,
			Package:                  "test-plan",
			PackageBilling:           "usage",
			PackageBillingContractId: "",
			Location:                 "AMS Amsterdam",
			LocationID:               1,
			ServerStatus:             "RUNNING",
			PowerStatus:              "on",
			Installed:                1,
			CloudPool:                "",
		},
	}
}

// WithID sets the server ID (mbpkgid)
func (b *ServerBuilder) WithID(id int) *ServerBuilder {
	b.server.ID = id
	return b
}

// WithName sets the server hostname/FQDN
func (b *ServerBuilder) WithName(name string) *ServerBuilder {
	b.server.Name = name
	return b
}

// WithOS sets the OS name and ID
func (b *ServerBuilder) WithOS(osName string, osID int) *ServerBuilder {
	b.server.OS = osName
	b.server.OSID = osID
	return b
}

// WithLocation sets the location name and ID
func (b *ServerBuilder) WithLocation(locationName string, locationID int) *ServerBuilder {
	b.server.Location = locationName
	b.server.LocationID = locationID
	return b
}

// WithStatus sets the server status (RUNNING, TERMINATED, etc.)
func (b *ServerBuilder) WithStatus(status string) *ServerBuilder {
	b.server.ServerStatus = status
	return b
}

// WithIPv4 sets the primary IPv4 address
func (b *ServerBuilder) WithIPv4(ip string) *ServerBuilder {
	b.server.PrimaryIPv4 = ip
	return b
}

// WithIPv6 sets the primary IPv6 address
func (b *ServerBuilder) WithIPv6(ip string) *ServerBuilder {
	b.server.PrimaryIPv6 = ip
	return b
}

// WithPackageBilling sets the package billing mode and contract ID
func (b *ServerBuilder) WithPackageBilling(billing, contractID string) *ServerBuilder {
	b.server.PackageBilling = billing
	b.server.PackageBillingContractId = contractID
	return b
}

// WithCloudPool sets the cloud pool name
func (b *ServerBuilder) WithCloudPool(pool string) *ServerBuilder {
	b.server.CloudPool = pool
	return b
}

// NotInstalled marks the server as not installed (Installed = 0)
func (b *ServerBuilder) NotInstalled() *ServerBuilder {
	b.server.Installed = 0
	b.server.Name = ""
	b.server.OS = ""
	return b
}

// Build returns the constructed Server
func (b *ServerBuilder) Build() gona.Server {
	return b.server
}

// SSHKeyBuilder provides a fluent API for creating test SSHKey instances.
type SSHKeyBuilder struct {
	key gona.SSHKey
}

// NewSSHKeyBuilder creates a new SSHKeyBuilder with sensible defaults
func NewSSHKeyBuilder() *SSHKeyBuilder {
	return &SSHKeyBuilder{
		key: gona.SSHKey{
			ID:          100,
			Name:        "test-key",
			Key:         "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... test@example.com",
			Fingerprint: "SHA256:test-fingerprint-hash",
		},
	}
}

// WithID sets the SSH key ID
func (b *SSHKeyBuilder) WithID(id int) *SSHKeyBuilder {
	b.key.ID = id
	return b
}

// WithName sets the SSH key name
func (b *SSHKeyBuilder) WithName(name string) *SSHKeyBuilder {
	b.key.Name = name
	return b
}

// WithKey sets the SSH public key content
func (b *SSHKeyBuilder) WithKey(key string) *SSHKeyBuilder {
	b.key.Key = key
	return b
}

// WithFingerprint sets the SSH key fingerprint
func (b *SSHKeyBuilder) WithFingerprint(fingerprint string) *SSHKeyBuilder {
	b.key.Fingerprint = fingerprint
	return b
}

// Build returns the constructed SSHKey
func (b *SSHKeyBuilder) Build() gona.SSHKey {
	return b.key
}

// BGPSessionBuilder provides a fluent API for creating test BGPSession instances.
type BGPSessionBuilder struct {
	session gona.BGPSession
}

// NewBGPSessionBuilder creates a new BGPSessionBuilder with sensible defaults
func NewBGPSessionBuilder() *BGPSessionBuilder {
	return &BGPSessionBuilder{
		session: gona.BGPSession{
			ID:             500,
			CustomerIP:     "198.51.100.10",
			GroupID:        1,
			Locked:         0,
			Description:    "Test BGP Session",
			State:          "established",
			RoutesReceived: 100,
			ConfigStatus:   1,
			ProviderPeerIP: "198.51.100.1",
			Location:       "AMS Amsterdam",
			GroupName:      "test-group",
			ProviderIPType: "IPv4",
			ProviderAsn:    64512,
			CustomerAsn:    65000,
		},
	}
}

// WithID sets the BGP session ID
func (b *BGPSessionBuilder) WithID(id int) *BGPSessionBuilder {
	b.session.ID = id
	return b
}

// WithCustomerIP sets the customer IP address
func (b *BGPSessionBuilder) WithCustomerIP(ip string) *BGPSessionBuilder {
	b.session.CustomerIP = ip
	return b
}

// WithGroupID sets the group ID
func (b *BGPSessionBuilder) WithGroupID(groupID int) *BGPSessionBuilder {
	b.session.GroupID = groupID
	return b
}

// WithIPv6 configures the session for IPv6
func (b *BGPSessionBuilder) WithIPv6() *BGPSessionBuilder {
	b.session.ProviderIPType = "IPv6"
	b.session.CustomerIP = "2001:db8:bgp::10"
	b.session.ProviderPeerIP = "2001:db8:bgp::1"
	return b
}

// WithState sets the BGP session state
func (b *BGPSessionBuilder) WithState(state string) *BGPSessionBuilder {
	b.session.State = state
	return b
}

// WithASNs sets the provider and customer ASNs
func (b *BGPSessionBuilder) WithASNs(providerASN, customerASN int) *BGPSessionBuilder {
	b.session.ProviderAsn = providerASN
	b.session.CustomerAsn = customerASN
	return b
}

// Build returns the constructed BGPSession
func (b *BGPSessionBuilder) Build() gona.BGPSession {
	return b.session
}

// BuildPtr returns a pointer to the constructed BGPSession
func (b *BGPSessionBuilder) BuildPtr() *gona.BGPSession {
	session := b.Build()
	return &session
}

// IPsBuilder provides a fluent API for creating test IPs instances.
type IPsBuilder struct {
	ips gona.IPs
}

// NewIPsBuilder creates a new IPsBuilder with sensible defaults
func NewIPsBuilder() *IPsBuilder {
	return &IPsBuilder{
		ips: gona.IPs{
			IPv4: []gona.IP{
				{
					ID:        1,
					IP:        "192.0.2.100",
					Primary:   1,
					Gateway:   "192.0.2.1",
					Netmask:   "255.255.255.0",
					Broadcast: "192.0.2.255",
					Reverse:   "100.2.0.192.in-addr.arpa",
				},
			},
			IPv6: []gona.IP{
				{
					ID:      2,
					IP:      "2001:db8::100",
					Primary: 1,
					Gateway: "2001:db8::1",
					Netmask: "ffff:ffff:ffff:ffff::",
					Reverse: "",
				},
			},
		},
	}
}

// AddIPv4 adds an IPv4 address to the IPs
func (b *IPsBuilder) AddIPv4(ip, gateway, netmask string, primary bool) *IPsBuilder {
	primaryInt := 0
	if primary {
		primaryInt = 1
	}

	newIP := gona.IP{
		ID:        len(b.ips.IPv4) + 1,
		IP:        ip,
		Primary:   primaryInt,
		Gateway:   gateway,
		Netmask:   netmask,
		Broadcast: "",
	}
	b.ips.IPv4 = append(b.ips.IPv4, newIP)
	return b
}

// AddIPv6 adds an IPv6 address to the IPs
func (b *IPsBuilder) AddIPv6(ip, gateway, netmask string, primary bool) *IPsBuilder {
	primaryInt := 0
	if primary {
		primaryInt = 1
	}

	newIP := gona.IP{
		ID:      len(b.ips.IPv6) + 1,
		IP:      ip,
		Primary: primaryInt,
		Gateway: gateway,
		Netmask: netmask,
	}
	b.ips.IPv6 = append(b.ips.IPv6, newIP)
	return b
}

// ClearIPv4 removes all IPv4 addresses
func (b *IPsBuilder) ClearIPv4() *IPsBuilder {
	b.ips.IPv4 = []gona.IP{}
	return b
}

// ClearIPv6 removes all IPv6 addresses
func (b *IPsBuilder) ClearIPv6() *IPsBuilder {
	b.ips.IPv6 = []gona.IP{}
	return b
}

// Build returns the constructed IPs
func (b *IPsBuilder) Build() gona.IPs {
	return b.ips
}

// LocationBuilder provides a fluent API for creating test Location instances.
type LocationBuilder struct {
	location gona.Location
}

// NewLocationBuilder creates a new LocationBuilder with sensible defaults
func NewLocationBuilder() *LocationBuilder {
	return &LocationBuilder{
		location: gona.Location{
			ID:        1,
			Name:      "AMS Amsterdam",
			IATACode:  "AMS",
			Continent: "EU",
			Flag:      "🇳🇱",
			Disabled:  0,
		},
	}
}

// WithID sets the location ID
func (b *LocationBuilder) WithID(id int) *LocationBuilder {
	b.location.ID = id
	return b
}

// WithName sets the location name
func (b *LocationBuilder) WithName(name string) *LocationBuilder {
	b.location.Name = name
	return b
}

// WithIATACode sets the IATA code
func (b *LocationBuilder) WithIATACode(code string) *LocationBuilder {
	b.location.IATACode = code
	return b
}

// WithContinent sets the continent
func (b *LocationBuilder) WithContinent(continent string) *LocationBuilder {
	b.location.Continent = continent
	return b
}

// Disabled marks the location as disabled
func (b *LocationBuilder) Disabled() *LocationBuilder {
	b.location.Disabled = 1
	return b
}

// Build returns the constructed Location
func (b *LocationBuilder) Build() gona.Location {
	return b.location
}

// OSBuilder provides a fluent API for creating test OS instances.
type OSBuilder struct {
	os gona.OS
}

// NewOSBuilder creates a new OSBuilder with sensible defaults
func NewOSBuilder() *OSBuilder {
	return &OSBuilder{
		os: gona.OS{
			ID:      1,
			Os:      "Ubuntu 22.04 LTS",
			Type:    "linux",
			Subtype: "ubuntu",
			Size:    "10GB",
			Bits:    "64",
			Tech:    "kvm",
		},
	}
}

// WithID sets the OS ID
func (b *OSBuilder) WithID(id int) *OSBuilder {
	b.os.ID = id
	return b
}

// WithName sets the OS name
func (b *OSBuilder) WithName(name string) *OSBuilder {
	b.os.Os = name
	return b
}

// WithType sets the OS type
func (b *OSBuilder) WithType(osType string) *OSBuilder {
	b.os.Type = osType
	return b
}

// WithSubtype sets the OS subtype
func (b *OSBuilder) WithSubtype(subtype string) *OSBuilder {
	b.os.Subtype = subtype
	return b
}

// Build returns the constructed OS
func (b *OSBuilder) Build() gona.OS {
	return b.os
}

// Helper functions for creating quick test instances

// NewTestServer creates a test Server with the given ID and name
func NewTestServer(id int, name string) gona.Server {
	return NewServerBuilder().WithID(id).WithName(name).Build()
}

// NewRunningServer creates a test Server in RUNNING status
func NewRunningServer(id int) gona.Server {
	return NewServerBuilder().
		WithID(id).
		WithName(fmt.Sprintf("server-%d.example.com", id)).
		WithStatus("RUNNING").
		Build()
}

// NewTerminatedServer creates a test Server in TERMINATED status
func NewTerminatedServer(id int) gona.Server {
	return NewServerBuilder().
		WithID(id).
		WithStatus("TERMINATED").
		Build()
}

// NewTestSSHKey creates a test SSHKey with the given ID and name
func NewTestSSHKey(id int, name string) gona.SSHKey {
	return NewSSHKeyBuilder().WithID(id).WithName(name).Build()
}

// NewTestBGPSession creates a test BGPSession with the given ID
func NewTestBGPSession(id int) *gona.BGPSession {
	return NewBGPSessionBuilder().WithID(id).BuildPtr()
}

// NewTestBGPSessionIPv6 creates a test IPv6 BGP session with the given ID
func NewTestBGPSessionIPv6(id int) *gona.BGPSession {
	return NewBGPSessionBuilder().WithID(id).WithIPv6().BuildPtr()
}

// NewTestLocation creates a test Location with the given ID and name
func NewTestLocation(id int, name string) gona.Location {
	return NewLocationBuilder().WithID(id).WithName(name).Build()
}

// NewTestOS creates a test OS with the given ID and name
func NewTestOS(id int, name string) gona.OS {
	return NewOSBuilder().WithID(id).WithName(name).Build()
}

// NewTestIPs creates test IPs for a server
func NewTestIPs(serverID int) gona.IPs {
	return NewIPsBuilder().
		ClearIPv4().
		ClearIPv6().
		AddIPv4(fmt.Sprintf("192.0.2.%d", serverID%256), "192.0.2.1", "255.255.255.0", true).
		AddIPv6(fmt.Sprintf("2001:db8::%x", serverID), "2001:db8::1", "ffff:ffff:ffff:ffff::", true).
		Build()
}
