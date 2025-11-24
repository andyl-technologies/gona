# Testing Guide

This document describes the testing infrastructure available in the gona library.

## Overview

The gona library provides comprehensive testing utilities to make it easy to write tests for code that depends on the NetActuate API:

- **Interfaces** - Operation-specific interfaces for flexible mocking (in `gona` package)
- **Test Builders** - Fluent builders for creating test data (in `testing` package)
- **FakeClient** - In-memory fake for integration-style tests (in `testing` package)
- **MockClient** - Configurable mock for unit tests (in `testing` package)

## Package Structure

Test utilities are organized in a dedicated subpackage to separate test code from production code:

```go
import (
    "github.com/netactuate/gona/gona"           // Production types and Client
    gonatesting "github.com/netactuate/gona/testing"  // Test utilities
)
```

The main `gona` package contains:
- Production `Client` implementation
- Operation interfaces (ServerOperations, SSHKeyOperations, etc.)
- All domain types (Server, SSHKey, BGPSession, etc.)

The `gona/testing` subpackage contains:
- Test data builders (ServerBuilder, SSHKeyBuilder, etc.)
- FakeClient implementation
- MockClient implementation

## Interfaces

The library defines focused interfaces that group related operations:

### ServerOperations

Methods for managing VPS server instances:
- `CreateServer(ctx, *CreateServerRequest) (ServerBuild, error)`
- `GetServer(ctx, id) (Server, error)`
- `BuildServer(ctx, id, *BuildServerRequest) (ServerBuild, error)`
- `DeleteServer(ctx, id, cancelBilling) error`
- `UnlinkServer(ctx, id) error`

### SSHKeyOperations

Methods for managing SSH keys:
- `CreateSSHKey(ctx, name, key) (SSHKey, error)`
- `GetSSHKey(ctx, id) (SSHKey, error)`
- `DeleteSSHKey(ctx, id) error`

### BGPSessionOperations

Methods for managing BGP sessions:
- `CreateBGPSessions(ctx, mbPkgID, groupID, isIPV6, redundant) (*BGPSession, error)`
- `GetBGPSessions(ctx, mbPkgID) ([]*BGPSession, error)`

### IPOperations

Methods for managing IP addresses:
- `GetIPs(ctx, mbPkgID) (IPs, error)`

### MetadataOperations

Methods for querying metadata:
- `GetLocations(ctx) ([]Location, error)`
- `GetOSs(ctx) ([]OS, error)`

### ClientInterface

Composed interface that includes all operation interfaces. The real `*Client` type implements this interface.

## Test Builders

Fluent builders make it easy to create test data with sensible defaults.

### ServerBuilder

```go
server := gonatesting.NewServerBuilder().
    WithID(123).
    WithName("test-server.example.com").
    WithStatus("RUNNING").
    WithIPv4("192.0.2.10").
    WithLocation("AMS Amsterdam", 1).
    Build()
```

Quick helpers:
```go
server := gonatesting.NewRunningServer(123)       // Server in RUNNING state
server := gonatesting.NewTerminatedServer(456)    // Server in TERMINATED state
server := gonatesting.NewTestServer(789, "test")  // Custom ID and name
```

### SSHKeyBuilder

```go
key := gonatesting.NewSSHKeyBuilder().
    WithID(100).
    WithName("deploy-key").
    WithKey("ssh-rsa AAAAB3...").
    Build()

// Quick helper
key := gonatesting.NewTestSSHKey(100, "my-key")
```

### BGPSessionBuilder

```go
session := gonatesting.NewBGPSessionBuilder().
    WithID(500).
    WithCustomerIP("198.51.100.10").
    WithGroupID(1).
    WithState("established").
    BuildPtr()

// Quick helpers
session := gonatesting.NewTestBGPSession(500)       // IPv4 session
session := gonatesting.NewTestBGPSessionIPv6(501)   // IPv6 session
```

### IPsBuilder

```go
ips := gonatesting.NewIPsBuilder().
    ClearIPv4().
    ClearIPv6().
    AddIPv4("192.0.2.10", "192.0.2.1", "255.255.255.0", true).
    AddIPv6("2001:db8::10", "2001:db8::1", "ffff:ffff:ffff:ffff::", true).
    Build()

// Quick helper
ips := gonatesting.NewTestIPs(123)  // IPs for server ID 123
```

### LocationBuilder

```go
location := gonatesting.NewLocationBuilder().
    WithID(1).
    WithName("AMS Amsterdam").
    WithIATACode("AMS").
    WithContinent("EU").
    Build()

// Quick helper
location := gonatesting.NewTestLocation(1, "AMS Amsterdam")
```

### OSBuilder

```go
os := gonatesting.NewOSBuilder().
    WithID(1).
    WithName("Ubuntu 22.04 LTS").
    WithType("linux").
    WithSubtype("ubuntu").
    Build()

// Quick helper
os := gonatesting.NewTestOS(1, "Ubuntu 22.04 LTS")
```

## FakeClient

`FakeClient` is an in-memory fake implementation that maintains state. Use this for integration-style tests where you need realistic stateful behavior.

### Basic Usage

```go
import (
    "context"
    "testing"

    "github.com/netactuate/gona/gona"
    gonatesting "github.com/netactuate/gona/testing"
)

func TestServerLifecycle(t *testing.T) {
    // Create fake client with default test data
    fake := gonatesting.NewFakeClient()
    ctx := context.Background()

    // Create a server
    req := &gona.CreateServerRequest{
        Plan:     "test-plan",
        Location: 1,
        Image:    1,
        FQDN:     "test.example.com",
    }
    build, err := fake.CreateServer(ctx, req)
    if err != nil {
        t.Fatalf("CreateServer failed: %v", err)
    }

    // Retrieve the server
    server, err := fake.GetServer(ctx, build.ServerID)
    if err != nil {
        t.Fatalf("GetServer failed: %v", err)
    }

    if server.Name != "test.example.com" {
        t.Errorf("expected name test.example.com, got %s", server.Name)
    }
}
```

### Pre-populating State

```go
func TestWithExistingServer(t *testing.T) {
    fake := gonatesting.NewFakeClient()
    ctx := context.Background()

    // Add server to fake's state
    server := gonatesting.NewRunningServer(123)
    fake.AddServer(server)

    // Now GetServer will return it
    retrieved, err := fake.GetServer(ctx, 123)
    if err != nil {
        t.Fatalf("GetServer failed: %v", err)
    }

    if retrieved.ID != 123 {
        t.Errorf("expected ID 123, got %d", retrieved.ID)
    }
}
```

### Error Injection

```go
func TestCreateServerError(t *testing.T) {
    fake := gonatesting.NewFakeClient()
    ctx := context.Background()

    // Inject error for CreateServer
    fake.CreateServerError = fmt.Errorf("quota exceeded")

    req := &gona.CreateServerRequest{
        Plan:     "test-plan",
        Location: 1,
        Image:    1,
    }

    _, err := fake.CreateServer(ctx, req)
    if err == nil {
        t.Fatal("expected error, got nil")
    }

    if err.Error() != "quota exceeded" {
        t.Errorf("expected 'quota exceeded', got '%s'", err.Error())
    }
}
```

### Call Tracking

```go
func TestCallTracking(t *testing.T) {
    fake := gonatesting.NewFakeClient()
    ctx := context.Background()

    fake.GetServer(ctx, 123)
    fake.GetServer(ctx, 456)

    calls := fake.GetCalls()
    if len(calls) != 2 {
        t.Errorf("expected 2 calls, got %d", len(calls))
    }

    expectedCalls := []string{"GetServer(123)", "GetServer(456)"}
    for i, call := range calls {
        if call != expectedCalls[i] {
            t.Errorf("call %d: expected %s, got %s", i, expectedCalls[i], call)
        }
    }

    // Reset for next test
    fake.ResetCalls()
}
```

### Helper Methods

FakeClient provides helper methods to pre-populate state:

- `AddServer(server Server)` - Add a server
- `AddSSHKey(key SSHKey)` - Add an SSH key
- `AddBGPSession(mbPkgID int, session *BGPSession)` - Add a BGP session
- `SetIPs(mbPkgID int, ips IPs)` - Set IPs for a server
- `AddLocation(location Location)` - Add a location
- `AddOS(os OS)` - Add an OS

### Default Test Data

FakeClient comes pre-loaded with:

**3 Locations:**
- ID 1: "AMS Amsterdam" (EU)
- ID 2: "LAX Los Angeles" (NA)
- ID 3: "SJC San Jose" (NA)

**3 Operating Systems:**
- ID 1: "Ubuntu 22.04 LTS"
- ID 2: "Debian 12"
- ID 3: "Rocky Linux 9"

## MockClient

`MockClient` provides configurable function fields for precise control in unit tests. Each method can be configured with a custom function.

### Basic Usage

```go
func TestWithMock(t *testing.T) {
    ctx := context.Background()

    mock := &gonatesting.MockClient{
        GetServerFunc: func(ctx context.Context, id int) (gona.Server, error) {
            // Return exact response you need
            return gona.Server{
                ID:           id,
                Name:         "mocked-server",
                ServerStatus: "RUNNING",
            }, nil
        },
    }

    server, err := mock.GetServer(ctx, 999)
    if err != nil {
        t.Fatalf("GetServer failed: %v", err)
    }

    if server.ID != 999 {
        t.Errorf("expected ID 999, got %d", server.ID)
    }
}
```

### Testing Error Conditions

```go
func TestMockError(t *testing.T) {
    ctx := context.Background()

    mock := &gonatesting.MockClient{
        DeleteServerFunc: func(ctx context.Context, id int, cancelBilling bool) error {
            if id == 404 {
                return fmt.Errorf("server not found")
            }
            return nil
        },
    }

    // Test error case
    err := mock.DeleteServer(ctx, 404, false)
    if err == nil {
        t.Fatal("expected error, got nil")
    }

    // Test success case
    err = mock.DeleteServer(ctx, 123, false)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

### Call Tracking

```go
func TestMockCallTracking(t *testing.T) {
    ctx := context.Background()

    mock := &gonatesting.MockClient{
        CreateSSHKeyFunc: func(ctx context.Context, name, key string) (gona.SSHKey, error) {
            return gona.SSHKey{ID: 1, Name: name, Key: key}, nil
        },
    }

    mock.CreateSSHKey(ctx, "key1", "ssh-rsa AAA...")
    mock.CreateSSHKey(ctx, "key2", "ssh-rsa BBB...")

    calls := mock.GetCalls()
    if len(calls) != 2 {
        t.Errorf("expected 2 calls, got %d", len(calls))
    }
}
```

### Not Configured Error

If you call a method without configuring its function, MockClient returns a clear error:

```go
func TestNotConfigured(t *testing.T) {
    mock := &gonatesting.MockClient{}
    ctx := context.Background()

    _, err := mock.GetServer(ctx, 123)
    if err == nil {
        t.Fatal("expected error for unconfigured method")
    }

    // Error will be: "GetServerFunc not configured"
}
```

## Combining Approaches

You can combine builders with fake/mock clients for flexible testing:

```go
func TestServerCreationWithBuilder(t *testing.T) {
    ctx := context.Background()

    // Use builder to create expected server
    expectedServer := gonatesting.NewServerBuilder().
        WithID(123).
        WithName("test.example.com").
        WithStatus("RUNNING").
        Build()

    // Configure mock to return it
    mock := &gonatesting.MockClient{
        GetServerFunc: func(ctx context.Context, id int) (gona.Server, error) {
            if id == expectedServer.ID {
                return expectedServer, nil
            }
            return gona.Server{}, fmt.Errorf("server not found")
        },
    }

    // Test
    server, err := mock.GetServer(ctx, 123)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if server.Name != "test.example.com" {
        t.Errorf("expected name test.example.com, got %s", server.Name)
    }
}
```

## Testing with Interfaces

You can write tests that depend on only the operations you need:

```go
// Your code depends only on ServerOperations
func ScaleServers(client gona.ServerOperations, count int) error {
    // Implementation using only server methods
    return nil
}

// Test with mock implementing only ServerOperations
func TestScaleServers(t *testing.T) {
    mock := &gonatesting.MockClient{
        CreateServerFunc: func(ctx context.Context, r *gona.CreateServerRequest) (gona.ServerBuild, error) {
            return gona.ServerBuild{ServerID: 123}, nil
        },
    }

    // Mock satisfies ServerOperations interface
    err := ScaleServers(mock, 3)
    if err != nil {
        t.Fatalf("ScaleServers failed: %v", err)
    }
}
```

## Best Practices

### Use Builders for Test Data

✅ **Good** - Using builders with defaults:
```go
server := gonatesting.NewServerBuilder().WithID(123).Build()
```

❌ **Avoid** - Manual struct construction:
```go
server := gona.Server{
    ID: 123,
    Name: "",  // Easy to forget required fields
    OS: "",
    // ... many more fields
}
```

### Choose the Right Test Double

- **FakeClient** - Integration-style tests, stateful behavior, multiple operations
- **MockClient** - Unit tests, precise control, single operations
- **Builders** - Creating test data for any approach

### Interface Segregation

Depend on specific operation interfaces instead of the full ClientInterface:

✅ **Good**:
```go
func DeployServer(client gona.ServerOperations, req *gona.CreateServerRequest) error
```

❌ **Avoid**:
```go
func DeployServer(client gona.ClientInterface, req *gona.CreateServerRequest) error
```

### Error Injection

Always test error paths using error injection:

```go
fake := gonatesting.NewFakeClient()
fake.CreateServerError = fmt.Errorf("quota exceeded")
// Test error handling
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out

# Run specific test
go test -v -run TestServerLifecycle ./gona
```
