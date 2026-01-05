package gonatesting_test

import (
	"testing"

	gonatesting "github.com/netactuate/gona/testing"
)

func TestServerBuilder(t *testing.T) {
	server := gonatesting.NewServerBuilder().
		WithID(123).
		WithName("test.example.com").
		WithStatus("STOPPED").
		Build()

	if server.ID != 123 {
		t.Errorf("expected ID 123, got %d", server.ID)
	}
	if server.Name != "test.example.com" {
		t.Errorf("expected name 'test.example.com', got %q", server.Name)
	}
	if server.ServerStatus != "STOPPED" {
		t.Errorf("expected status 'STOPPED', got %q", server.ServerStatus)
	}
}

func TestServerBuilder_Defaults(t *testing.T) {
	server := gonatesting.NewServerBuilder().Build()

	if server.ID != 1000 {
		t.Errorf("expected default ID 1000, got %d", server.ID)
	}
	if server.ServerStatus != "RUNNING" {
		t.Errorf("expected default status 'RUNNING', got %q", server.ServerStatus)
	}
	if server.PowerStatus != "on" {
		t.Errorf("expected default power status 'on', got %q", server.PowerStatus)
	}
	// Verify RFC-reserved IP ranges
	if server.PrimaryIPv4 != "192.0.2.100" {
		t.Errorf("expected IPv4 192.0.2.100, got %q", server.PrimaryIPv4)
	}
	if server.PrimaryIPv6 != "2001:db8::100" {
		t.Errorf("expected IPv6 2001:db8::100, got %q", server.PrimaryIPv6)
	}
}

func TestServerBuilder_AllMethods(t *testing.T) {
	server := gonatesting.NewServerBuilder().
		WithID(999).
		WithName("full-test.example.com").
		WithOS("Debian 12", 2).
		WithLocation("LAX Los Angeles", 2).
		WithStatus("INSTALLING").
		WithIPv4("192.0.2.50").
		WithIPv6("2001:db8::50").
		WithPackageBilling("monthly", "contract-123").
		WithCloudPool("AMD EPYC").
		Build()

	if server.ID != 999 {
		t.Errorf("expected ID 999, got %d", server.ID)
	}
	if server.OS != "Debian 12" {
		t.Errorf("expected OS 'Debian 12', got %q", server.OS)
	}
	if server.OSID != 2 {
		t.Errorf("expected OSID 2, got %d", server.OSID)
	}
	if server.Location != "LAX Los Angeles" {
		t.Errorf("expected location 'LAX Los Angeles', got %q", server.Location)
	}
	if server.ServerStatus != "INSTALLING" {
		t.Errorf("expected status 'INSTALLING', got %q", server.ServerStatus)
	}
	if server.CloudPool != "AMD EPYC" {
		t.Errorf("expected cloud pool 'AMD EPYC', got %q", server.CloudPool)
	}
}

func TestNewRunningServer(t *testing.T) {
	server := gonatesting.NewRunningServer(555)

	if server.ID != 555 {
		t.Errorf("expected ID 555, got %d", server.ID)
	}
	if server.ServerStatus != "RUNNING" {
		t.Errorf("expected status 'RUNNING', got %q", server.ServerStatus)
	}
	if server.PowerStatus != "on" {
		t.Errorf("expected power 'on', got %q", server.PowerStatus)
	}
}

func TestSSHKeyBuilder(t *testing.T) {
	key := gonatesting.NewSSHKeyBuilder().
		WithID(50).
		WithName("deploy-key").
		WithKey("ssh-ed25519 AAAA...").
		Build()

	if key.ID != 50 {
		t.Errorf("expected ID 50, got %d", key.ID)
	}
	if key.Name != "deploy-key" {
		t.Errorf("expected name 'deploy-key', got %q", key.Name)
	}
}

func TestSSHKeyBuilder_Defaults(t *testing.T) {
	key := gonatesting.NewSSHKeyBuilder().Build()

	if key.ID != 100 {
		t.Errorf("expected default ID 100, got %d", key.ID)
	}
	if key.Name != "test-key" {
		t.Errorf("expected default name 'test-key', got %q", key.Name)
	}
}

func TestNewTestSSHKey(t *testing.T) {
	key := gonatesting.NewTestSSHKey(77, "my-ssh-key")

	if key.ID != 77 {
		t.Errorf("expected ID 77, got %d", key.ID)
	}
	if key.Name != "my-ssh-key" {
		t.Errorf("expected name 'my-ssh-key', got %q", key.Name)
	}
}

func TestBGPSessionBuilder(t *testing.T) {
	// Note: WithIPv6() must be called before WithCustomerIP() since it overwrites CustomerIP
	session := gonatesting.NewBGPSessionBuilder().
		WithID(200).
		WithIPv6().
		WithCustomerIP("2001:db8::custom").
		WithGroupID(5).
		WithState("pending").
		Build()

	if session.ID != 200 {
		t.Errorf("expected ID 200, got %d", session.ID)
	}
	if session.CustomerIP != "2001:db8::custom" {
		t.Errorf("expected customer IP '2001:db8::custom', got %q", session.CustomerIP)
	}
	if session.State != "pending" {
		t.Errorf("expected state 'pending', got %q", session.State)
	}
	if session.ProviderIPType != "IPv6" {
		t.Errorf("expected IPv6 type, got %q", session.ProviderIPType)
	}
}

func TestBGPSessionBuilder_Defaults(t *testing.T) {
	session := gonatesting.NewBGPSessionBuilder().Build()

	if session.State != "established" {
		t.Errorf("expected default state 'established', got %q", session.State)
	}
	if session.ProviderAsn != 64512 {
		t.Errorf("expected default provider ASN 64512, got %d", session.ProviderAsn)
	}
}

func TestNewTestBGPSession(t *testing.T) {
	session := gonatesting.NewTestBGPSession(300)

	if session.ID != 300 {
		t.Errorf("expected ID 300, got %d", session.ID)
	}
	if session.ProviderIPType != "IPv4" {
		t.Errorf("expected IPv4, got %q", session.ProviderIPType)
	}
}

func TestNewTestBGPSessionIPv6(t *testing.T) {
	session := gonatesting.NewTestBGPSessionIPv6(301)

	if session.ID != 301 {
		t.Errorf("expected ID 301, got %d", session.ID)
	}
	if session.ProviderIPType != "IPv6" {
		t.Errorf("expected IPv6, got %q", session.ProviderIPType)
	}
}

func TestIPsBuilder(t *testing.T) {
	// Note: NewIPsBuilder() starts with default IPs, AddIPv4/AddIPv6 adds to them
	ips := gonatesting.NewIPsBuilder().
		AddIPv4("192.0.2.10", "192.0.2.1", "255.255.255.0", true).
		AddIPv4("192.0.2.11", "192.0.2.1", "255.255.255.0", false).
		AddIPv6("2001:db8::10", "2001:db8::1", "ffff:ffff:ffff:ffff::", true).
		Build()

	// Builder has 1 default IPv4 + 2 added = 3
	if len(ips.IPv4) != 3 {
		t.Errorf("expected 3 IPv4 addresses (1 default + 2 added), got %d", len(ips.IPv4))
	}
	// Builder has 1 default IPv6 + 1 added = 2
	if len(ips.IPv6) != 2 {
		t.Errorf("expected 2 IPv6 addresses (1 default + 1 added), got %d", len(ips.IPv6))
	}

	// Check that first (default) IPv4 is primary
	if ips.IPv4[0].Primary != 1 {
		t.Error("expected first (default) IPv4 to be primary")
	}
}

func TestNewTestIPs(t *testing.T) {
	ips := gonatesting.NewTestIPs(42)

	if len(ips.IPv4) != 1 {
		t.Fatalf("expected 1 IPv4, got %d", len(ips.IPv4))
	}
	if len(ips.IPv6) != 1 {
		t.Fatalf("expected 1 IPv6, got %d", len(ips.IPv6))
	}

	// Verify IPs use server ID
	if ips.IPv4[0].IP != "192.0.2.42" {
		t.Errorf("expected IPv4 192.0.2.42, got %q", ips.IPv4[0].IP)
	}
}

func TestLocationBuilder(t *testing.T) {
	loc := gonatesting.NewLocationBuilder().
		WithID(5).
		WithName("NYC New York").
		WithIATACode("NYC").
		WithContinent("NA").
		Build()

	if loc.ID != 5 {
		t.Errorf("expected ID 5, got %d", loc.ID)
	}
	if loc.Name != "NYC New York" {
		t.Errorf("expected name 'NYC New York', got %q", loc.Name)
	}
	if loc.IATACode != "NYC" {
		t.Errorf("expected IATA 'NYC', got %q", loc.IATACode)
	}
}

func TestLocationBuilder_Defaults(t *testing.T) {
	loc := gonatesting.NewLocationBuilder().Build()

	if loc.ID != 1 {
		t.Errorf("expected default ID 1, got %d", loc.ID)
	}
	if loc.Continent != "EU" {
		t.Errorf("expected default continent 'EU', got %q", loc.Continent)
	}
}

func TestOSBuilder(t *testing.T) {
	os := gonatesting.NewOSBuilder().
		WithID(10).
		WithName("Arch Linux").
		WithType("linux").
		Build()

	if os.ID != 10 {
		t.Errorf("expected ID 10, got %d", os.ID)
	}
	if os.Os != "Arch Linux" {
		t.Errorf("expected name 'Arch Linux', got %q", os.Os)
	}
	if os.Type != "linux" {
		t.Errorf("expected type 'linux', got %q", os.Type)
	}
}

func TestOSBuilder_Defaults(t *testing.T) {
	os := gonatesting.NewOSBuilder().Build()

	if os.ID != 1 {
		t.Errorf("expected default ID 1, got %d", os.ID)
	}
	if os.Os != "Ubuntu 22.04 LTS" {
		t.Errorf("expected default OS 'Ubuntu 22.04 LTS', got %q", os.Os)
	}
	if os.Bits != "64" {
		t.Errorf("expected default bits '64', got %q", os.Bits)
	}
}
