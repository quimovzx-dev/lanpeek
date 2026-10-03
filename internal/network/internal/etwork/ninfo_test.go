package network

import (
	"bufio"
	"strings"
	"testing"
)

func TestGatewayParsing(t *testing.T) {
	routes := `Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT
eth0 00000000 0101A8C0 0003 0 0 100 00000000 0 0 0
eth0 00000000 0201A8C0 0003 0 0 200 00000000 0 0 0
eth1 00000000 0100000A 0003 0 0 10 00000000 0 0 0
eth2 00000000 0100000A 0001 0 0 10 00000000 0 0 0
bad
`
	got := parseGateways(bufio.NewScanner(strings.NewReader(routes)))
	if got["eth0"] != "192.168.1.1" || got["eth1"] != "10.0.0.1" || len(got) != 2 {
		t.Fatalf("%v", got)
	}
}

func TestInfo(t *testing.T) {
	rows, err := Info()
	if err != nil || len(rows) == 0 {
		t.Fatalf("%v %v", rows, err)
	}
	for _, row := range rows {
		if row.Name == "" || row.Addresses == nil {
			t.Fatalf("%+v", row)
		}
	}
}
