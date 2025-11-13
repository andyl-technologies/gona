package gona

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateServer_WithCloudConfig(t *testing.T) {
	// Sample cloud-init configuration
	const cloudInitYAML = `#cloud-config
users:
  - name: demo
    sudo: ALL=(ALL) NOPASSWD:ALL
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0g+ZTxC5wxi2sspF9cY demo@example.com
packages:
  - curl
  - vim
`
	// Encode to Base64 as expected by the API
	cloudConfigBase64 := base64.StdEncoding.EncodeToString([]byte(cloudInitYAML))

	tests := map[string]struct {
		request        *CreateServerRequest
		expectedParams map[string]string
		mockResponse   ServerBuild
		expectError    bool
	}{
		"with cloud config": {
			request: &CreateServerRequest{
				Plan:        "plan-test",
				Location:    1,
				Image:       100,
				FQDN:        "test.example.com",
				CloudConfig: cloudConfigBase64,
			},
			expectedParams: map[string]string{
				"plan":           "plan-test",
				"location":       "1",
				"image":          "100",
				"fqdn":           "test.example.com",
				"script_type":    "cloud_init",
				"script_content": cloudConfigBase64,
			},
			mockResponse: ServerBuild{
				ServerID: 12345,
				Status:   "building",
				Build:    1,
			},
			expectError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify HTTP method
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify endpoint path
				if r.URL.Path != "/cloud/server/buy_build" {
					t.Errorf("Expected path /cloud/server/buy_build, got %s", r.URL.Path)
				}

				// Parse request body as form data
				if err := r.ParseForm(); err != nil {
					t.Fatalf("Failed to parse form: %v", err)
				}

				// Verify expected query parameters
				for key, expectedValue := range tt.expectedParams {
					actualValue := r.Form.Get(key)
					if actualValue != expectedValue {
						t.Errorf("Parameter %s = %q, want %q", key, actualValue, expectedValue)
					}
				}

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockAPIResponse(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server
			client := NewClientCustom("test-api-key", server.URL+"/")
			ctx := context.Background()

			// Call CreateServer
			result, err := client.CreateServer(ctx, tt.request)

			// Verify error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("CreateServer() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				// Verify response data
				if result.ServerID != tt.mockResponse.ServerID {
					t.Errorf("ServerID = %d, want %d", result.ServerID, tt.mockResponse.ServerID)
				}
				if result.Status != tt.mockResponse.Status {
					t.Errorf("Status = %q, want %q", result.Status, tt.mockResponse.Status)
				}
				if result.Build != tt.mockResponse.Build {
					t.Errorf("Build = %d, want %d", result.Build, tt.mockResponse.Build)
				}
			}
		})
	}
}

func TestBuildServer_WithCloudConfig(t *testing.T) {
	// Sample cloud-init configuration
	const cloudInitYAML = `#cloud-config
runcmd:
  - echo "Server rebuilt with cloud-init" > /etc/motd
write_files:
  - path: /etc/test.conf
    content: |
      # Test configuration
      option = value
`
	// Encode to Base64 as expected by the API
	cloudConfigBase64 := base64.StdEncoding.EncodeToString([]byte(cloudInitYAML))

	tests := map[string]struct {
		serverID       int
		request        *BuildServerRequest
		expectedParams map[string]string
		mockResponse   ServerBuild
		expectError    bool
	}{
		"with cloud config": {
			serverID: 12345,
			request: &BuildServerRequest{
				Image:       200,
				FQDN:        "rebuilt.example.com",
				CloudConfig: cloudConfigBase64,
			},
			expectedParams: map[string]string{
				"image":          "200",
				"fqdn":           "rebuilt.example.com",
				"script_type":    "cloud_init",
				"script_content": cloudConfigBase64,
			},
			mockResponse: ServerBuild{
				ServerID: 12345,
				Status:   "rebuilding",
				Build:    2,
			},
			expectError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify HTTP method
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify endpoint path contains server ID
				expectedPath := "/cloud/server/build/12345"
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Parse request body as form data
				if err := r.ParseForm(); err != nil {
					t.Fatalf("Failed to parse form: %v", err)
				}

				// Verify expected query parameters
				for key, expectedValue := range tt.expectedParams {
					actualValue := r.Form.Get(key)
					if actualValue != expectedValue {
						t.Errorf("Parameter %s = %q, want %q", key, actualValue, expectedValue)
					}
				}

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockAPIResponse(tt.mockResponse))
			}))
			defer server.Close()

			// Create client with mock server
			client := NewClientCustom("test-api-key", server.URL+"/")
			ctx := context.Background()

			// Call BuildServer
			result, err := client.BuildServer(ctx, tt.serverID, tt.request)

			// Verify error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("BuildServer() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				// Verify response data
				if result.ServerID != tt.mockResponse.ServerID {
					t.Errorf("ServerID = %d, want %d", result.ServerID, tt.mockResponse.ServerID)
				}
				if result.Status != tt.mockResponse.Status {
					t.Errorf("Status = %q, want %q", result.Status, tt.mockResponse.Status)
				}
				if result.Build != tt.mockResponse.Build {
					t.Errorf("Build = %d, want %d", result.Build, tt.mockResponse.Build)
				}
			}
		})
	}
}
