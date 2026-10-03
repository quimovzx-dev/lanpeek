# lanpeek

> Fast LAN network toolkit CLI written in Go - discover devices, scan ports, inspect your network.

![Go version](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go)
![License: MIT](https://img.shields.io/badge/license-MIT-green)
![Release](https://img.shields.io/badge/release-coming%20soon-lightgrey)

> **Status: v1 source available.** Linux-first, unprivileged TCP connect probes. No tagged release or prebuilt binaries yet.

## Features

- **Device discovery** - find hosts on your local network from a CIDR range
- **Port scanning** - fast TCP connect probes against a host
- **Network info** - show your interface and network details at a glance
- **Nice terminal output** - clean, colorful tables built with lipgloss
- **Single binary** - no runtime dependencies

v1 targets Linux first. Other platforms may follow.

## Install

Requires Go 1.22 or later. Build from source:

```sh
git clone https://github.com/quimovzx-dev/lanpeek.git
cd lanpeek
go build -o lanpeek .
./lanpeek --help
```

There are no prebuilt releases yet.

## Usage

```sh
# Show info about your network
./lanpeek info

# Discover devices on a subnet
./lanpeek discover --cidr 192.168.1.0/24

# Scan specific ports on a host
./lanpeek ports 192.168.1.10 --ports 22,80,443
```

Use `--json` for machine-readable output, `--timeout 500ms` for the per-probe
limit and `--concurrency 32` for the worker count. Add `--no-dns` to discovery
to skip reverse DNS. Progress is printed to stderr, never mixed into JSON.

Discovery checks ports 22, 80, 443, 445 and 3389. A TCP refusal counts as evidence
that a host is reachable; no response does not prove a host is offline. IPv4
prefixes are limited to 4096 usable hosts. Port service names are conventional
hints, not service identification. Gateways are best-effort Linux IPv4 only.
See [DEVELOPMENT.md](DEVELOPMENT.md) for details and testing.

## Roadmap

- [x] `info` command
- [x] `discover` command
- [x] `ports` command (TCP connect probes)
- [x] Styled output with lipgloss
- [ ] First tagged release with goreleaser
- [x] Linux support (v1)
- [ ] macOS and Windows support
- [x] JSON output for scripting

## Built with

- [Go](https://go.dev)
- [cobra](https://github.com/spf13/cobra)
- [lipgloss](https://github.com/charmbracelet/lipgloss)

## Contributing

Ideas, bug reports and pull requests are welcome. Open an issue to discuss a change before sending a large PR.

Only scan networks and hosts you own or have permission to test.

## License

MIT. See `LICENSE`.
