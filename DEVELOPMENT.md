# Building v1

Requires Go 1.22 or later. Linux is the primary target. No root permissions,
raw packets, packet capture, external executables or runtime dependencies.

```sh
go build -o lanpeek .
go test -race ./...
go vet ./...
./lanpeek info
./lanpeek ports 127.0.0.1 --ports 22,80,443
./lanpeek discover --cidr 192.168.1.0/24 --no-dns
```

Only scan networks you own or have permission to test.

Global flags:

- `--json`: machine-readable arrays on stdout; progress stays on stderr.
- `--timeout 500ms`: maximum duration per TCP connection or reverse DNS lookup.
  Must be greater than zero and at most 30 seconds.
- `--concurrency 32`: worker count, from 1 to 256.

Discovery supports IPv4 prefixes with at most 4096 usable hosts. It excludes
network and broadcast addresses through /30; /31 and /32 follow point-to-point
and single-host semantics. It checks TCP ports 22, 80, 443, 445 and 3389 in that
order, stopping at the first response. A refused connection counts as reachable.
Unresponsive hosts can still be online behind a firewall. Discovery is not an
ARP, ICMP or complete device inventory. Reverse DNS is best-effort, bounded by
`--timeout`, and only attempted for reachable hosts; `--no-dns` disables it.

Port scanning accepts an IP or hostname and comma-separated ports (no ranges).
Results are ordered numerically, duplicate ports are removed, and statuses are
`open`, `refused`, `timeout`, or `error`. Errors retain their details in table and
JSON output. Service names are hints from standard port conventions, not actual
service detection. Hostname scans use the Go dialer's resolution/address choice,
not an exhaustive scan of every address associated with the hostname.

`info` reports each interface's up/down flag, addresses, masks, local MAC and MTU
(the MTU is included in JSON). Linux IPv4 gateways are read from `/proc/net/route`
on a best-effort basis; IPv6 gateways are not currently reported. Blank values
are omitted or empty in JSON and shown as `-` in tables. IPv6 masks are hex.

Ctrl-C cancels scans and exits nonzero, without writing incomplete JSON results.
There is no release automation or tagged release in this source drop.
