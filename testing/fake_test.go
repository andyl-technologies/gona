package gonatesting_test

import (
	"context"
	"errors"
	"testing"

	"github.com/netactuate/gona/gona"
	gonatesting "github.com/netactuate/gona/testing"
)

func TestFakeClient_CreateServer(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

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

	if build.ServerID < 1000 {
		t.Errorf("expected ServerID >= 1000, got %d", build.ServerID)
	}
	if build.Status != "building" {
		t.Errorf("expected status 'building', got %q", build.Status)
	}

	// Verify server can be retrieved
	server, err := fake.GetServer(ctx, build.ServerID)
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}

	if server.Name != "test.example.com" {
		t.Errorf("expected name test.example.com, got %s", server.Name)
	}
	if server.ServerStatus != "RUNNING" {
		t.Errorf("expected status RUNNING, got %s", server.ServerStatus)
	}
}

func TestFakeClient_CreateServer_Validation(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *gona.CreateServerRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "cannot be nil",
		},
		{
			name:    "empty plan",
			req:     &gona.CreateServerRequest{Location: 1},
			wantErr: "Plan is required",
		},
		{
			name:    "invalid location",
			req:     &gona.CreateServerRequest{Plan: "test", Location: 0},
			wantErr: "invalid Location ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fake.CreateServer(ctx, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != "" && !contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFakeClient_ErrorInjection(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	expectedErr := errors.New("quota exceeded")
	fake.CreateServerError = expectedErr

	_, err := fake.CreateServer(ctx, &gona.CreateServerRequest{Plan: "test", Location: 1})
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestFakeClient_CallTracking(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	// Perform some operations
	fake.GetLocations(ctx)
	fake.GetOSs(ctx)

	calls := fake.GetCalls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}

	if !wasCalled(calls, "GetLocations") {
		t.Error("expected GetLocations to be tracked")
	}
	if !wasCalled(calls, "GetOSs") {
		t.Error("expected GetOSs to be tracked")
	}

	// Test ResetCalls
	fake.ResetCalls()
	if len(fake.GetCalls()) != 0 {
		t.Error("expected calls to be cleared after ResetCalls")
	}
}

func TestFakeClient_Reset(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	// Add some data
	fake.AddServer(gona.Server{ID: 999, Name: "test"})
	fake.CreateServerError = errors.New("test error")
	fake.GetLocations(ctx) // Add a call

	// Reset
	fake.Reset()

	// Verify server is gone
	_, err := fake.GetServer(ctx, 999)
	if err == nil {
		t.Error("expected error after reset, server should not exist")
	}

	// Verify error is cleared
	if fake.CreateServerError != nil {
		t.Error("expected error to be nil after reset")
	}

	// GetCalls() will have 1 call (GetServer from checking if server is gone)
	// So we verify Reset clears the previous calls by checking only GetServer is tracked
	calls := fake.GetCalls()
	if len(calls) != 1 || !wasCalled(calls, "GetServer") {
		t.Errorf("expected only GetServer call after reset, got %v", calls)
	}

	// Verify default locations still exist
	locs, _ := fake.GetLocations(ctx)
	if len(locs) != 3 {
		t.Errorf("expected 3 default locations, got %d", len(locs))
	}
}

func TestFakeClient_SSHKeyOperations(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	// Create
	key, err := fake.CreateSSHKey(ctx, "my-key", "ssh-rsa AAAA...")
	if err != nil {
		t.Fatalf("CreateSSHKey failed: %v", err)
	}
	if key.Name != "my-key" {
		t.Errorf("expected name 'my-key', got %q", key.Name)
	}

	// Get
	retrieved, err := fake.GetSSHKey(ctx, key.ID)
	if err != nil {
		t.Fatalf("GetSSHKey failed: %v", err)
	}
	if retrieved.ID != key.ID {
		t.Errorf("expected ID %d, got %d", key.ID, retrieved.ID)
	}

	// Delete
	err = fake.DeleteSSHKey(ctx, key.ID)
	if err != nil {
		t.Fatalf("DeleteSSHKey failed: %v", err)
	}

	// Verify deleted
	_, err = fake.GetSSHKey(ctx, key.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestFakeClient_CreateSSHKey_Validation(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	_, err := fake.CreateSSHKey(ctx, "", "ssh-rsa AAAA...")
	if err == nil {
		t.Error("expected error for empty name")
	}
	if !contains(err.Error(), "cannot be empty") {
		t.Errorf("expected 'cannot be empty' in error, got %q", err.Error())
	}
}

func TestFakeClient_BGPSessionOperations(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	// Create
	session, err := fake.CreateBGPSessions(ctx, 100, 1, false, false)
	if err != nil {
		t.Fatalf("CreateBGPSessions failed: %v", err)
	}
	if session.ProviderIPType != "IPv4" {
		t.Errorf("expected IPv4, got %s", session.ProviderIPType)
	}

	// Get
	sessions, err := fake.GetBGPSessions(ctx, 100)
	if err != nil {
		t.Fatalf("GetBGPSessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
}

func TestFakeClient_CreateBGPSessions_Validation(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	tests := []struct {
		name    string
		mbPkgID int
		groupID int
		wantErr string
	}{
		{"invalid mbPkgID", 0, 1, "invalid mbPkgID"},
		{"invalid groupID", 1, 0, "invalid groupID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fake.CreateBGPSessions(ctx, tt.mbPkgID, tt.groupID, false, false)
			if err == nil {
				t.Fatal("expected error")
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFakeClient_GetLocations(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	locs, err := fake.GetLocations(ctx)
	if err != nil {
		t.Fatalf("GetLocations failed: %v", err)
	}

	if len(locs) != 3 {
		t.Errorf("expected 3 default locations, got %d", len(locs))
	}

	// Verify defensive copy - modifying returned slice shouldn't affect internal state
	locs[0].Name = "MODIFIED"
	locs2, _ := fake.GetLocations(ctx)
	if locs2[0].Name == "MODIFIED" {
		t.Error("defensive copy failed - internal state was modified")
	}
}

func TestFakeClient_GetOSs(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	oses, err := fake.GetOSs(ctx)
	if err != nil {
		t.Fatalf("GetOSs failed: %v", err)
	}

	if len(oses) != 3 {
		t.Errorf("expected 3 default OSes, got %d", len(oses))
	}
}

func TestFakeClient_BuildServer_Validation(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int
		req     *gona.BuildServerRequest
		wantErr string
	}{
		{
			name:    "invalid id",
			id:      0,
			req:     &gona.BuildServerRequest{},
			wantErr: "invalid server ID",
		},
		{
			name:    "nil request",
			id:      1,
			req:     nil,
			wantErr: "cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fake.BuildServer(ctx, tt.id, tt.req)
			if err == nil {
				t.Fatal("expected error")
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFakeClient_AddServer(t *testing.T) {
	fake := gonatesting.NewFakeClient()
	ctx := context.Background()

	server := gona.Server{
		ID:           12345,
		Name:         "added-server",
		ServerStatus: "RUNNING",
	}

	fake.AddServer(server)

	retrieved, err := fake.GetServer(ctx, 12345)
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}
	if retrieved.Name != "added-server" {
		t.Errorf("expected name 'added-server', got %q", retrieved.Name)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func wasCalled(calls []string, method string) bool {
	for _, c := range calls {
		if c == method || contains(c, method) {
			return true
		}
	}
	return false
}
