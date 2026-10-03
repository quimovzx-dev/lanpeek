// Package network implements unprivileged TCP connect scans. No raw packets are used.
package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const MaxHosts = 4096

var CommonPorts = []int{22, 80, 443, 445, 3389}

// PortResult describes a TCP connect attempt; service is a conventional hint,
// not identification of the program actually listening.
type PortResult struct {
	Port    int    `json:"port"`
	Service string `json:"service_hint"`
	Status  string `json:"status"`
	Detail  string `json:"detail,omitempty"`
}

type HostResult struct {
	IP       string `json:"ip"`
	Status   string `json:"status"`
	Hostname string `json:"hostname,omitempty"`
}

func Service(port int) string {
	names := map[int]string{21: "ftp", 22: "ssh", 23: "telnet", 25: "smtp", 53: "dns", 80: "http", 110: "pop3", 143: "imap", 443: "https", 445: "smb", 465: "smtps", 587: "submission", 993: "imaps", 995: "pop3s", 3306: "mysql", 3389: "rdp", 5432: "postgresql", 6379: "redis", 8080: "http-alt", 8443: "https-alt"}
	if name, ok := names[port]; ok {
		return name
	}
	return "unknown"
}

func ParsePorts(value string) ([]int, error) {
	seen := make(map[int]bool)
	var ports []int
	for _, part := range strings.Split(value, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > 65535 {
			return nil, fmt.Errorf("invalid port %q: use comma-separated numbers from 1 to 65535", part)
		}
		if !seen[n] {
			seen[n] = true
			ports = append(ports, n)
		}
	}
	sort.Ints(ports)
	return ports, nil
}

// Hosts returns usable IPv4 addresses, excluding network and broadcast for
// prefixes through /30. Both /31 addresses and the single /32 address are usable.
func Hosts(cidr string) ([]string, error) {
	p, err := netip.ParsePrefix(cidr)
	if err != nil || !p.Addr().Is4() {
		return nil, fmt.Errorf("CIDR must be an IPv4 prefix, for example 192.168.1.0/24")
	}
	p = p.Masked()
	count := uint64(1) << uint(32-p.Bits())
	usable := count
	if p.Bits() <= 30 {
		usable -= 2
	}
	if usable > MaxHosts {
		return nil, fmt.Errorf("CIDR has %d hosts; limit is %d (use /20 or narrower)", usable, MaxHosts)
	}
	first := p.Addr()
	if p.Bits() <= 30 {
		first = first.Next()
	}
	hosts := make([]string, 0, int(usable))
	for a, i := first, uint64(0); i < usable; a, i = a.Next(), i+1 {
		hosts = append(hosts, a.String())
	}
	return hosts, nil
}

func Probe(ctx context.Context, host string, port int, timeout time.Duration) PortResult {
	r := PortResult{Port: port, Service: Service(port)}
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	switch {
	case err == nil:
		_ = conn.Close()
		r.Status = "open"
	case errors.Is(err, syscall.ECONNREFUSED):
		r.Status = "refused"
	case ctx.Err() != nil:
		r.Status = "canceled"
	case errors.Is(err, context.DeadlineExceeded):
		r.Status = "timeout"
	default:
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			r.Status = "timeout"
		} else {
			r.Status = "error"
			r.Detail = err.Error()
		}
	}
	return r
}

// parallel maps work using a fixed worker pool and returns results in input order.
func parallel[T any](ctx context.Context, count, concurrency int, work func(int) T) ([]T, error) {
	results := make([]T, count)
	jobs := make(chan int)
	done := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for n := range jobs {
				if ctx.Err() == nil {
					results[n] = work(n)
				}
			}
		}()
	}
feed:
	for n := 0; n < count; n++ {
		select {
		case jobs <- n:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	for i := 0; i < concurrency; i++ {
		<-done
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func validOptions(timeout time.Duration, concurrency int) error {
	if timeout <= 0 || timeout > 30*time.Second {
		return fmt.Errorf("timeout must be greater than 0 and at most 30s")
	}
	if concurrency < 1 || concurrency > 256 {
		return fmt.Errorf("concurrency must be between 1 and 256")
	}
	return nil
}

func ScanPorts(ctx context.Context, host string, ports []int, timeout time.Duration, concurrency int) ([]PortResult, error) {
	if err := validOptions(timeout, concurrency); err != nil {
		return nil, err
	}
	return parallel(ctx, len(ports), concurrency, func(i int) PortResult { return Probe(ctx, host, ports[i], timeout) })
}

func Discover(ctx context.Context, hosts []string, timeout time.Duration, concurrency int, dns bool) ([]HostResult, error) {
	if err := validOptions(timeout, concurrency); err != nil {
		return nil, err
	}
	return parallel(ctx, len(hosts), concurrency, func(i int) HostResult {
		r := HostResult{IP: hosts[i], Status: "unresponsive"}
		for _, port := range CommonPorts {
			p := Probe(ctx, r.IP, port, timeout)
			// A refusal is also positive evidence that this IP answered a TCP probe.
			if p.Status == "open" || p.Status == "refused" {
				r.Status = "reachable"
				break
			}
			if ctx.Err() != nil {
				break
			}
		}
		if dns && r.Status == "reachable" && ctx.Err() == nil {
			dnsCtx, cancel := context.WithTimeout(ctx, timeout)
			names, err := net.DefaultResolver.LookupAddr(dnsCtx, r.IP)
			cancel()
			if err == nil && len(names) > 0 {
				r.Hostname = strings.TrimSuffix(names[0], ".")
			}
		}
		return r
	})
}
