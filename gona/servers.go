package gona

import (
	"context"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

// Server struct defines what a VPS looks like
type Server struct {
	Name                     string `json:"fqdn"`
	ID                       int    `json:"mbpkgid"`
	OS                       string `json:"os"`
	OSID                     int    `json:"os_id"`
	PrimaryIPv4              string `json:"ip"`
	PrimaryIPv6              string `json:"ipv6"`
	PlanID                   int    `json:"plan_id"`
	Package                  string `json:"package"`
	PackageBilling           string `json:"package_billing"`
	PackageBillingContractId string `json:"package_billing_contract_id"`
	Location                 string `json:"city"`
	LocationID               int    `json:"location_id"`
	ServerStatus             string `json:"status"`
	PowerStatus              string `json:"state"`
	Installed                int    `json:"installed"`
	CloudPool                string `json:"cloud_pool_name"`
}

// GetServers external method on Client to list your instances
func (c *Client) GetServers(ctx context.Context) ([]Server, error) {
	var serverList []Server
	if err := c.get(ctx, "cloud/servers", &serverList); err != nil {
		return nil, err
	}
	return serverList, nil
}

// GetServer external method on Client to get an instance
func (c *Client) GetServer(ctx context.Context, id int) (server Server, err error) {
	if err := c.get(ctx, "cloud/server?mbpkgid="+strconv.Itoa(id), &server); err != nil {
		return server, err
	}
	return server, nil
}

// CreateServerRequest is as set of parameters for a server creation call.
type CreateServerRequest struct {
	Plan                     string `url:"plan,omitempty"`
	Location                 int    `url:"location,omitempty"`
	Image                    int    `url:"image,omitempty"`
	FQDN                     string `url:"fqdn,omitempty"`
	SSHKey                   string `url:"ssh_key,omitempty"`
	SSHKeyID                 int    `url:"ssh_key_id,omitempty"`
	Password                 string `url:"password,omitempty"`
	PackageBilling           string `url:"package_billing,omitempty"`
	PackageBillingContractId string `url:"package_billing_contract_id,omitempty"`
	CloudConfig              string // Assumed to be Base64 encoded if set
	ScriptContent            string `url:"script_content,omitempty"`
	Params                   string `url:"params,omitempty"`

	CloudPool CloudPool `url:"cloud_pool_id,omitempty"`
}

// ServerBuild is a server creation response message.
type ServerBuild struct {
	ServerID int    `json:"mbpkgid"`
	Status   string `json:"status"`
	Build    int    `json:"build"`
}

// CreateServer external method on Client to buy and build a new instance.
func (c *Client) CreateServer(ctx context.Context, r *CreateServerRequest) (b ServerBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if r.CloudConfig != "" {
		values.Add("script_type", "cloud_init")
		values.Set("script_content", r.CloudConfig)
	} else if values.Has("script_content") {
		values.Add("script_type", "user-data")
	}

	if err := c.post(ctx, "cloud/server/buy_build", []byte(values.Encode()), &b); err != nil {
		return b, err
	}

	return b, nil
}

// BuildServerRequest is a set of parameters for a server re-building call.
type BuildServerRequest struct {
	Plan                     string `url:"plan,omitempty"`
	Location                 int    `url:"location,omitempty"`
	Image                    int    `url:"image,omitempty"`
	FQDN                     string `url:"fqdn,omitempty"`
	SSHKey                   string `url:"ssh_key,omitempty"`
	SSHKeyID                 int    `url:"ssh_key_id,omitempty"`
	Password                 string `url:"password,omitempty"`
	PackageBilling           string `url:"package_billing,omitempty"`
	PackageBillingContractId string `url:"package_billing_contract_id,omitempty"`
	CloudConfig              string // Assumed to be Base64 encoded if set
	ScriptContent            string `url:"script_content,omitempty"`
	Params                   string `url:"params,omitempty"`
}

// BuildServer external method on Client to re-build an instance
func (c *Client) BuildServer(ctx context.Context, id int, r *BuildServerRequest) (b ServerBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}

	if r.CloudConfig != "" {
		values.Add("script_type", "cloud_init")
		values.Set("script_content", r.CloudConfig)
	} else if values.Has("script_content") {
		values.Add("script_type", "user-data")
	}

	// if r.Params != "" {
	//     values.Add("params", r.Params)
	// }

	if err := c.post(ctx, "cloud/server/build/"+strconv.Itoa(id), []byte(values.Encode()), &b); err != nil {
		return b, err
	}

	return b, nil
}

// DeleteServer external method on Client to destroy an instance.
func (c *Client) DeleteServer(ctx context.Context, id int, cancelBilling bool) error {
	values := url.Values{}
	if cancelBilling {
		values.Add("cancel_billing", "1")
	}
	return c.post(ctx, "cloud/server/delete?mbpkgid="+strconv.Itoa(id), []byte(values.Encode()), nil)
}

// UnlinkServer external method on Client to unlink a billing package from a location
func (c *Client) UnlinkServer(ctx context.Context, id int) error {
	return c.post(ctx, "cloud/server/unlink/"+strconv.Itoa(id), nil, nil)
}

// StartServer external method on Client to boot up an instance
func (c *Client) StartServer(ctx context.Context, id int) error {

	return c.post(ctx, "cloud/server/start/"+strconv.Itoa(id), nil, nil)
}

// StopServer external method on Client to shut down an instance
func (c *Client) StopServer(ctx context.Context, id int) error {

	return c.post(ctx, "cloud/server/shutdown/"+strconv.Itoa(id), nil, nil)
}
