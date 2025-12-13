// Package gonatesting provides test utilities for the gona library,
// including test data builders, a stateful FakeClient for integration
// tests, and a configurable MockClient for unit tests.
//
// # Test Data Builders
//
// Builders provide fluent APIs for creating test data with sensible defaults:
//
//	server := NewServerBuilder().
//	    WithID(123).
//	    WithName("test.example.com").
//	    Build()
//
// Quick helpers are also available for common scenarios:
//
//	server := NewRunningServer(123)
//	key := NewTestSSHKey(100, "my-key")
//
// # FakeClient
//
// FakeClient is an in-memory implementation of gona.ClientInterface
// for integration-style tests:
//
//	fake := NewFakeClient()
//	fake.AddServer(NewRunningServer(123))
//	server, _ := fake.GetServer(ctx, 123)
//
// FakeClient supports:
//   - State management (servers, SSH keys, BGP sessions, IPs)
//   - Error injection via exported error fields
//   - Call tracking via GetCalls() and WasCalled()
//
// # MockClient
//
// MockClient provides function-field based mocking for unit tests:
//
//	mock := &MockClient{
//	    GetServerFunc: func(ctx context.Context, id int) (gona.Server, error) {
//	        return gona.Server{ID: id, Name: "mocked"}, nil
//	    },
//	}
package gonatesting
