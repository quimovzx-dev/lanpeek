package network

import (
	"context"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestParsePorts(t *testing.T) {
	got, err := ParsePorts("443, 22,80,22")
	if err != nil || !reflect.DeepEqual(got, []int{22, 80, 443}) {
		t.Fatalf("got %v, %v", got, err)
	}
	for _, s := range []string{"", "0", "65536", "80,", "abc", "22-80", "-1"} {
		if _, err := ParsePorts(s); err == nil {
			t.Errorf("accepted %q", s)
		}
	}
}

func TestHosts(t *testing.T) {
	for _, tc := range []struct {
		cidr string
		want []string
	}{
		{"192.168.1.2/30", []string{"192.168.1.1", "192.168.1.2"}},
		{"192.168.1.0/31", []string{"192.168.1.0", "192.168.1.1"}},
		{"127.0.0.1/32", []string{"127.0.0.1"}},
	} {
		got, err := Hosts(tc.cidr)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: %v, %v", tc.cidr, got, err)
		}
	}
	for _, cidr := range []string{"::1/128", "bad", "0.0.0.0/0", "10.0.0.0/19"} {
		if _, err := Hosts(cidr); err == nil {
			t.Errorf("accepted %s", cidr)
		}
	}
	got, err := Hosts("10.0.0.0/20")
	if err != nil || len(got) != 4094 {
		t.Fatalf("/20: %d, %v", len(got), err)
	}
}

func TestLocalProbes(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	defer listener.Close()
	got, err := ScanPorts(context.Background(), "127.0.0.1", []int{port}, time.Second, 2)
	if err != nil || len(got) != 1 || got[0].Status != "open" {
		t.Fatalf("open: %+v %v", got, err)
	}
	// Close the bound port before testing refusal. No remote networks are probed.
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	r := Probe(context.Background(), "127.0.0.1", port, time.Second)
	if r.Status != "refused" {
		t.Fatalf("port %s: %+v", strconv.Itoa(port), r)
	}
	discovered, err := Discover(context.Background(), []string{"127.0.0.1"}, time.Second, 2, false)
	if err != nil || discovered[0].Status != "reachable" || discovered[0].Hostname != "" {
		t.Fatalf("discover: %+v %v", discovered, err)
	}
}

func TestCanceledAndInvalidOptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ScanPorts(ctx, "127.0.0.1", []int{22}, time.Second, 1); err == nil {
		t.Fatal("missing canceled error")
	}
	for _, concurrency := range []int{0, 257} {
		if _, err := ScanPorts(context.Background(), "localhost", []int{22}, time.Second, concurrency); err == nil {
			t.Fatal("invalid concurrency accepted")
		}
	}
	if _, err := Discover(context.Background(), nil, 0, 1, false); err == nil {
		t.Fatal("zero timeout accepted")
	}
}

func TestWorkerPoolOrder(t *testing.T) {
	got, err := parallel(context.Background(), 100, 32, func(n int) int { return n * n })
	if err != nil {
		t.Fatal(err)
	}
	for n, v := range got {
		if v != n*n {
			t.Fatalf("at %d got %d", n, v)
		}
	}
}
