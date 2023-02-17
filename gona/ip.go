package gona

import (
	"strconv"

	"inet.af/netaddr"
)

type IPType string

const (
	IPv4 IPType = "ipv4"
	IPv6 IPType = "ipv6"
)

type IPs struct {
	IPv4 []IP `json:"IPv4"`
	IPv6 []IP `json:"IPv6"`
}

type IP struct {
	ID        int `json:"id"`
	Primary   int `json:"primary"`
	Reverse   string `json:"reverse"`
	IP        string `json:"ip"`
	Gateway   string `json:"gateway"`
	Netmask   string `json:"netmask"`
	Broadcast string `json:"broadcast"`
}

func (ips *IPs) GetIPsMap() *map[string]IPType {
	m := make(map[string]IPType)

	for _, ip := range ips.IPv4 {
		m[ip.IP] = IPv4
	}
	for _, ip := range ips.IPv6 {
		m[ip.IP] = IPv6

		ip6, err := netaddr.ParseIP(ip.IP)
		if err == nil {
			m[ip6.StringExpanded()] = IPv6
		}
	}

	return &m
}

// GetIPs returns a list of IPs for the selected mbPkgID from the API
func (c *Client) GetIPs(mbPkgID int) (ips IPs, err error) {
	if err := c.get("cloud/networkips/"+strconv.Itoa(mbPkgID), &ips); err != nil {
		return IPs{}, err
	}

	return ips, nil
}