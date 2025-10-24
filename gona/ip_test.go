package gona

import (
	"reflect"
	"testing"
)

func TestGetIPsMap(t *testing.T) {
	ips := IPs{
		IPv4: []IP{
			{IP: "127.0.0.1"},
		},
		IPv6: []IP{
			{IP: "::1"},
		},
	}

	expected := map[string]IPType{
		"127.0.0.1": IPv4,
		"::1":       IPv6,
		"0000:0000:0000:0000:0000:0000:0000:0001": IPv6,
	}

	result := ips.GetIPsMap()
	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("GetIPsMap() = %v; want %v", *result, expected)
	}
}
