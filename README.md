

# gona
`import "."`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package gona provides a simple golang interface to the NetActuate
Rest API at <a href="https://vapi2.netactuate.com/">https://vapi2.netactuate.com/</a>




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func GetKeyFromEnv() string](#GetKeyFromEnv)
* [type BGPCreateSessionsInput](#BGPCreateSessionsInput)
* [type BGPSession](#BGPSession)
  * [func (s *BGPSession) IsLocked() bool](#BGPSession.IsLocked)
  * [func (s *BGPSession) IsProviderIPTypeV4() bool](#BGPSession.IsProviderIPTypeV4)
* [type BuildServerRequest](#BuildServerRequest)
* [type Client](#Client)
  * [func NewClient(apikey string) *Client](#NewClient)
  * [func NewClientCustom(apikey string, apiurl string) *Client](#NewClientCustom)
  * [func (c *Client) BuildServer(id int, r *BuildServerRequest) (b ServerBuild, err error)](#Client.BuildServer)
  * [func (c *Client) CreateBGPSessions(mbPkgID int, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error)](#Client.CreateBGPSessions)
  * [func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error)](#Client.CreateSSHKey)
  * [func (c *Client) CreateServer(r *CreateServerRequest) (b ServerBuild, err error)](#Client.CreateServer)
  * [func (c *Client) DeleteSSHKey(id int) error](#Client.DeleteSSHKey)
  * [func (c *Client) DeleteServer(id int, cancelBilling bool) error](#Client.DeleteServer)
  * [func (c *Client) GetBGPSession(id int) (*BGPSession, error)](#Client.GetBGPSession)
  * [func (c *Client) GetBGPSessions(mbPkgID int) ([]*BGPSession, error)](#Client.GetBGPSessions)
  * [func (c *Client) GetIPs(mbPkgID int) (ips IPs, err error)](#Client.GetIPs)
  * [func (c *Client) GetLocations() ([]Location, error)](#Client.GetLocations)
  * [func (c *Client) GetOSs() ([]OS, error)](#Client.GetOSs)
  * [func (c *Client) GetPackage(id int) (pkg Package, err error)](#Client.GetPackage)
  * [func (c *Client) GetPackages() ([]Package, error)](#Client.GetPackages)
  * [func (c *Client) GetPlans() ([]Plan, error)](#Client.GetPlans)
  * [func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error)](#Client.GetSSHKey)
  * [func (c *Client) GetSSHKeys() (keys []SSHKey, err error)](#Client.GetSSHKeys)
  * [func (c *Client) GetServer(id int) (server Server, err error)](#Client.GetServer)
  * [func (c *Client) GetServers() ([]Server, error)](#Client.GetServers)
  * [func (c *Client) StartServer(id int) error](#Client.StartServer)
  * [func (c *Client) StopServer(id int) error](#Client.StopServer)
  * [func (c *Client) UnlinkServer(id int) error](#Client.UnlinkServer)
* [type CreateServerRequest](#CreateServerRequest)
* [type IP](#IP)
* [type IPType](#IPType)
* [type IPs](#IPs)
  * [func (ips *IPs) GetIPsMap() *map[string]IPType](#IPs.GetIPsMap)
* [type Location](#Location)
* [type OS](#OS)
* [type Package](#Package)
* [type Plan](#Plan)
* [type Prefix](#Prefix)
* [type SSHKey](#SSHKey)
* [type Server](#Server)
* [type ServerBuild](#ServerBuild)


#### <a name="pkg-files">Package files</a>
[bgp.go](/src/target/bgp.go) [client.go](/src/target/client.go) [ip.go](/src/target/ip.go) [locations.go](/src/target/locations.go) [os.go](/src/target/os.go) [packages.go](/src/target/packages.go) [plans.go](/src/target/plans.go) [servers.go](/src/target/servers.go) [sshkeys.go](/src/target/sshkeys.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    Version      = "0.2.0"
    BaseEndpoint = "https://vapi2.netactuate.com/api/"
    ContentType  = "application/json"
)
```
Version, BaseEndpoint, ContentType constants




## <a name="GetKeyFromEnv">func</a> [GetKeyFromEnv](/src/target/client.go?s=768:795#L38)
``` go
func GetKeyFromEnv() string
```
GetKeyFromEnv is a simple function to grab the value for
"NA_API_KEY" from the environment




## <a name="BGPCreateSessionsInput">type</a> [BGPCreateSessionsInput](/src/target/bgp.go?s=2787:3020#L104)
``` go
type BGPCreateSessionsInput struct {
    MbPkgID   int `json:"mbpkgid"` // Contract BGP ID
    GroupID   int `json:"group_id"`
    Redundant int `json:"redundant"` //Force session redundancy
    IPV6      int `json:"ipv6"`      // IPv6 Session
}

```









## <a name="BGPSession">type</a> [BGPSession](/src/target/bgp.go?s=55:1107#L9)
``` go
type BGPSession struct {
    ID             int         `json:"id"`
    CustomerIP     string      `json:"customer_peer_ip"`
    GroupID        int         `json:"group_id"`
    Locked         int         `json:"locked"`
    Description    string      `json:"description"`
    State          interface{} `json:"state"`
    RoutesReceived interface{} `json:"routes_received"`
    LastUpdate     interface{} `json:"last_update"`
    ConfigStatus   int         `json:"config_status"`
    Password       interface{} `json:"password"`
    Prefixes       []Prefix    `json:"prefixes"`
    ExportList     string      `json:"export_list"`
    Community      interface{} `json:"community"`
    ProviderPeerIP string      `json:"provider_peer_ip"`
    Location       string      `json:"location"`
    Latitude       string      `json:"latitude"`
    Longitude      string      `json:"longitude"`
    GroupName      string      `json:"group_name"`
    ProviderIPType string      `json:"provider_ip_type"`
    ProviderAsn    int         `json:"provider_asn,string"`
    CustomerAsn    int         `json:"customer_asn,string"`
}

```









### <a name="BGPSession.IsLocked">func</a> (\*BGPSession) [IsLocked](/src/target/bgp.go?s=1604:1640#L47)
``` go
func (s *BGPSession) IsLocked() bool
```



### <a name="BGPSession.IsProviderIPTypeV4">func</a> (\*BGPSession) [IsProviderIPTypeV4](/src/target/bgp.go?s=1668:1714#L51)
``` go
func (s *BGPSession) IsProviderIPTypeV4() bool
```



## <a name="BuildServerRequest">type</a> [BuildServerRequest](/src/target/servers.go?s=2907:3337#L86)
``` go
type BuildServerRequest struct {
    Location      int    `url:"location,omitempty"`
    Image         int    `url:"image,omitempty"`
    FQDN          string `url:"fqdn,omitempty"`
    SSHKey        string `url:"ssh_key,omitempty"`
    SSHKeyID      int    `url:"ssh_key_id,omitempty"`
    Password      string `url:"password,omitempty"`
    CloudConfig   string `url:"cloud_config,omitempty"`
    ScriptContent string `url:"script_content,omitempty"`
}

```
BuildServerRequest is a set of parameters for a server re-building call.










## <a name="Client">type</a> [Client](/src/target/client.go?s=567:669#L29)
``` go
type Client struct {
    // contains filtered or unexported fields
}

```
Client is the main object (struct) to which we attach most
methods/functions.
It has the following fields:
(client, userAgent, endPoint, apiKey)







### <a name="NewClient">func</a> [NewClient](/src/target/client.go?s=1540:1577#L66)
``` go
func NewClient(apikey string) *Client
```
NewClient takes an apikey and calls NewClientCustom with the hardcoded
BaseEndpoint constant API URL


### <a name="NewClientCustom">func</a> [NewClientCustom](/src/target/client.go?s=1015:1073#L45)
``` go
func NewClientCustom(apikey string, apiurl string) *Client
```
NewClientCustom is the main entrypoint for instantiating a Client struct.
It takes your API Key as it's sole argument
and returns the Client struct ready to talk to the API





### <a name="Client.BuildServer">func</a> (\*Client) [BuildServer](/src/target/servers.go?s=3404:3490#L98)
``` go
func (c *Client) BuildServer(id int, r *BuildServerRequest) (b ServerBuild, err error)
```
BuildServer external method on Client to re-build an instance




### <a name="Client.CreateBGPSessions">func</a> (\*Client) [CreateBGPSessions](/src/target/bgp.go?s=3022:3132#L111)
``` go
func (c *Client) CreateBGPSessions(mbPkgID int, groupID int, isIPV6 bool, redundant bool) (*BGPSession, error)
```



### <a name="Client.CreateSSHKey">func</a> (\*Client) [CreateSSHKey](/src/target/sshkeys.go?s=770:844#L34)
``` go
func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error)
```
CreateSSHKey creates a key




### <a name="Client.CreateServer">func</a> (\*Client) [CreateServer](/src/target/servers.go?s=2475:2555#L69)
``` go
func (c *Client) CreateServer(r *CreateServerRequest) (b ServerBuild, err error)
```
CreateServer external method on Client to buy and build a new instance.




### <a name="Client.DeleteSSHKey">func</a> (\*Client) [DeleteSSHKey](/src/target/sshkeys.go?s=1091:1134#L47)
``` go
func (c *Client) DeleteSSHKey(id int) error
```
DeleteSSHKey deletes a key




### <a name="Client.DeleteServer">func</a> (\*Client) [DeleteServer](/src/target/servers.go?s=3846:3909#L115)
``` go
func (c *Client) DeleteServer(id int, cancelBilling bool) error
```
DeleteServer external method on Client to destroy an instance.




### <a name="Client.GetBGPSession">func</a> (\*Client) [GetBGPSession](/src/target/bgp.go?s=1828:1887#L56)
``` go
func (c *Client) GetBGPSession(id int) (*BGPSession, error)
```
GetBGPSession external method on Client to get your BGP session




### <a name="Client.GetBGPSessions">func</a> (\*Client) [GetBGPSessions](/src/target/bgp.go?s=2105:2172#L67)
``` go
func (c *Client) GetBGPSessions(mbPkgID int) ([]*BGPSession, error)
```
GetBGPSessions external method on Client to get BGP sessions




### <a name="Client.GetIPs">func</a> (\*Client) [GetIPs](/src/target/ip.go?s=817:874#L50)
``` go
func (c *Client) GetIPs(mbPkgID int) (ips IPs, err error)
```
GetIPs returns a list of IPs for the selected mbPkgID from the API




### <a name="Client.GetLocations">func</a> (\*Client) [GetLocations](/src/target/locations.go?s=384:435#L14)
``` go
func (c *Client) GetLocations() ([]Location, error)
```
GetLocations public method on Client to get a list of locations




### <a name="Client.GetOSs">func</a> (\*Client) [GetOSs](/src/target/os.go?s=349:388#L15)
``` go
func (c *Client) GetOSs() ([]OS, error)
```
GetOSs returns a list of OS objects from the api




### <a name="Client.GetPackage">func</a> (\*Client) [GetPackage](/src/target/packages.go?s=718:778#L29)
``` go
func (c *Client) GetPackage(id int) (pkg Package, err error)
```
GetPackage external method on Client that takes an id (int) as it's sole
argument and returns a single Package object




### <a name="Client.GetPackages">func</a> (\*Client) [GetPackages](/src/target/packages.go?s=400:449#L16)
``` go
func (c *Client) GetPackages() ([]Package, error)
```
GetPackages external method on Client that returns a
list of Package object from the API




### <a name="Client.GetPlans">func</a> (\*Client) [GetPlans](/src/target/plans.go?s=396:439#L15)
``` go
func (c *Client) GetPlans() ([]Plan, error)
```
GetPlans external method on Client to list available Plans




### <a name="Client.GetSSHKey">func</a> (\*Client) [GetSSHKey](/src/target/sshkeys.go?s=550:611#L26)
``` go
func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error)
```
GetSSHKey will list the information on a specific key




### <a name="Client.GetSSHKeys">func</a> (\*Client) [GetSSHKeys](/src/target/sshkeys.go?s=297:353#L17)
``` go
func (c *Client) GetSSHKeys() (keys []SSHKey, err error)
```
GetSSHKeys will list all SSH Keys installed for the account




### <a name="Client.GetServer">func</a> (\*Client) [GetServer](/src/target/servers.go?s=1227:1288#L39)
``` go
func (c *Client) GetServer(id int) (server Server, err error)
```
GetServer external method on Client to get an instance




### <a name="Client.GetServers">func</a> (\*Client) [GetServers](/src/target/servers.go?s=985:1032#L30)
``` go
func (c *Client) GetServers() ([]Server, error)
```
GetServers external method on Client to list your instances




### <a name="Client.StartServer">func</a> (\*Client) [StartServer](/src/target/servers.go?s=4357:4399#L129)
``` go
func (c *Client) StartServer(id int) error
```
StartServer external method on Client to boot up an instance




### <a name="Client.StopServer">func</a> (\*Client) [StopServer](/src/target/servers.go?s=4582:4623#L139)
``` go
func (c *Client) StopServer(id int) error
```
StopServer external method on Client to shut down an instance




### <a name="Client.UnlinkServer">func</a> (\*Client) [UnlinkServer](/src/target/servers.go?s=4178:4221#L124)
``` go
func (c *Client) UnlinkServer(id int) error
```
UnlinkServer external method on Client to unlink a billing package from a location




## <a name="CreateServerRequest">type</a> [CreateServerRequest](/src/target/servers.go?s=1495:2216#L47)
``` go
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
    CloudConfig              string `url:"cloud_config,omitempty"`
    ScriptContent            string `url:"script_content,omitempty"`
}

```
CreateServerRequest is as set of parameters for a server creation call.










## <a name="IP">type</a> [IP](/src/target/ip.go?s=203:452#L21)
``` go
type IP struct {
    ID        int    `json:"id"`
    Primary   int    `json:"primary"`
    Reverse   string `json:"reverse"`
    IP        string `json:"ip"`
    Gateway   string `json:"gateway"`
    Netmask   string `json:"netmask"`
    Broadcast string `json:"broadcast"`
}

```









## <a name="IPType">type</a> [IPType](/src/target/ip.go?s=57:75#L9)
``` go
type IPType string
```

``` go
const (
    IPv4 IPType = "ipv4"
    IPv6 IPType = "ipv6"
)
```









## <a name="IPs">type</a> [IPs](/src/target/ip.go?s=132:201#L16)
``` go
type IPs struct {
    IPv4 []IP `json:"IPv4"`
    IPv6 []IP `json:"IPv6"`
}

```









### <a name="IPs.GetIPsMap">func</a> (\*IPs) [GetIPsMap](/src/target/ip.go?s=454:500#L31)
``` go
func (ips *IPs) GetIPsMap() *map[string]IPType
```



## <a name="Location">type</a> [Location](/src/target/locations.go?s=87:315#L4)
``` go
type Location struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    IATACode  string `json:"iata_code"`
    Continent string `json:"continent"`
    Flag      string `json:"flat"`
    Disabled  int    `json:"disabled"`
}

```
Location is an API response message of available deployment locations










## <a name="OS">type</a> [OS](/src/target/os.go?s=68:295#L4)
``` go
type OS struct {
    ID      int    `json:"id"`
    Os      string `json:"os"`
    Type    string `json:"type"`
    Subtype string `json:"subtype"`
    Size    string `json:"size"`
    Bits    string `json:"bits"`
    Tech    string `json:"tech"`
}

```
OS is a struct for storing the attributes of an OS










## <a name="Package">type</a> [Package](/src/target/packages.go?s=86:303#L6)
``` go
type Package struct {
    ID        int    `json:"mbpkgid,string"`
    Status    string `json:"package_status"`
    Locked    string `json:"locked"`
    PlanName  string `json:"name"`
    Installed int    `json:"installed,string"`
}

```
Package struct stores the purchaced package values










## <a name="Plan">type</a> [Plan](/src/target/plans.go?s=69:332#L4)
``` go
type Plan struct {
    ID        int    `json:"plan_id,string"`
    Name      string `json:"plan"`
    RAM       string `json:"ram"`
    Disk      string `json:"disk"`
    Transfer  string `json:"transfer"`
    Price     string `json:"price"`
    Available string `json:"available"`
}

```
Plan struct defines the purchaceable plans/packages










## <a name="Prefix">type</a> [Prefix](/src/target/bgp.go?s=1109:1602#L33)
``` go
type Prefix struct {
    ID          int         `json:"id"`
    MbID        int         `json:"mb_id"`
    Prefix      string      `json:"prefix"`
    Append      interface{} `json:"append"`
    RuleType    string      `json:"rule_type"`
    PrefixType  string      `json:"prefix_type"`
    Description string      `json:"description"`
    Date        string      `json:"date"`
    AllowedPps  int         `json:"allowed_pps"`
    BgpGroupID  int         `json:"bgp_group_id"`
    PrefixID    int         `json:"prefix_id"`
}

```









## <a name="SSHKey">type</a> [SSHKey](/src/target/sshkeys.go?s=66:232#L9)
``` go
type SSHKey struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Key         string `json:"ssh_key"`
    Fingerprint string `json:"fingerprint"`
}

```
SSHKey Struct










## <a name="Server">type</a> [Server](/src/target/servers.go?s=138:920#L11)
``` go
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
}

```
Server struct defines what a VPS looks like










## <a name="ServerBuild">type</a> [ServerBuild](/src/target/servers.go?s=2272:2398#L62)
``` go
type ServerBuild struct {
    ServerID int    `json:"mbpkgid"`
    Status   string `json:"status"`
    Build    int    `json:"build"`
}

```
ServerBuild is a server creation response message.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
