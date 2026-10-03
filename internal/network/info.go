package network

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Address struct {
	IP     string `json:"ip"`
	Prefix int    `json:"prefix"`
	Mask   string `json:"mask"`
}
type Interface struct {
	Name      string    `json:"name"`
	Index     int       `json:"index"`
	State     string    `json:"state"`
	MAC       string    `json:"mac"`
	MTU       int       `json:"mtu"`
	Addresses []Address `json:"addresses"`
	Gateway   string    `json:"gateway,omitempty"`
}

func Info() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	gateways := linuxGateways()
	result := make([]Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		row := Interface{Name: iface.Name, Index: iface.Index, State: "down", MAC: iface.HardwareAddr.String(), MTU: iface.MTU, Addresses: []Address{}, Gateway: gateways[iface.Name]}
		if iface.Flags&net.FlagUp != 0 {
			row.State = "up"
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("addresses for %s: %w", iface.Name, err)
		}
		for _, addr := range addrs {
			ip, prefix, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			ones, bits := prefix.Mask.Size()
			mask := prefix.Mask.String()
			if bits == 32 {
				mask = net.IP(prefix.Mask).String()
			}
			row.Addresses = append(row.Addresses, Address{IP: ip.String(), Prefix: ones, Mask: mask})
		}
		result = append(result, row)
	}
	return result, nil
}

// Best-effort Linux IPv4 default gateways. Absence is not an error (other OS,
// restricted /proc, point-to-point routes, or IPv6-only default routing).
func linuxGateways() map[string]string {
	result := map[string]string{}
	if runtime.GOOS != "linux" {
		return result
	}
	f, err := os.Open("/proc/net/route")
	if err != nil {
		return result
	}
	defer f.Close()
	return parseGateways(bufio.NewScanner(f))
}

func parseGateways(scanner *bufio.Scanner) map[string]string {
	result := map[string]string{}
	metrics := map[string]uint64{}
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 8 || fields[1] != "00000000" {
			continue
		}
		flags, err := strconv.ParseUint(fields[3], 16, 32)
		if err != nil || flags&3 != 3 {
			continue
		}
		metric, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			continue
		}
		bytes, err := hex.DecodeString(fields[2])
		if err != nil || len(bytes) != 4 {
			continue
		}
		if old, ok := metrics[fields[0]]; ok && old <= metric {
			continue
		}
		metrics[fields[0]] = metric
		result[fields[0]] = net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
	}
	return result
}
