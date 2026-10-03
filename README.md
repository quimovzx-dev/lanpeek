# lanpeek

> Fast LAN network toolkit CLI written in Go - discover devices, scan ports, inspect your network.

![Go version](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go)
![License: MIT](https://img.shields.io/badge/license-MIT-green)
![Release](https://img.shields.io/badge/release-coming%20soon-lightgrey)

> **Status: early development.** The code is still being written. Everything below describes the planned v1.

## Features

- **Device discovery** - find hosts on your local network from a CIDR range (planned)
- **Port scanning** - fast TCP connect probes against a host (planned)
- **Network info** - show your interface and network details at a glance (planned)
- **Nice terminal output** - clean, colorful tables built with lipgloss (planned)
- **Single binary** - no runtime dependencies, released with goreleaser (planned)

v1 targets Linux first. Other platforms may follow.

## Install

Coming soon. Once the first release is out, you will be able to:

- Download a prebuilt binary from the [Releases](https://github.com/quimovzx-dev/lanpeek/releases) page
- Or build from source with Go:

```sh
go install github.com/quimovzx-dev/lanpeek@latest
```

Neither works yet, since there is no release.

## Usage (planned)

```sh
# Show info about your network
lanpeek info

# Discover devices on a subnet
lanpeek discover --cidr 192.168.1.0/24

# Scan specific ports on a host
lanpeek ports 192.168.1.10 --ports 22,80,443
```

Command names and flags may change before the first release.

## Roadmap

- [ ] `info` command
- [ ] `discover` command
- [ ] `ports` command (TCP connect probes)
- [ ] Styled output with lipgloss
- [ ] First tagged release with goreleaser
- [ ] Linux support (v1)
- [ ] macOS and Windows support
- [ ] JSON output for scripting

## Built with

- [Go](https://go.dev)
- [cobra](https://github.com/spf13/cobra)
- [lipgloss](https://github.com/charmbracelet/lipgloss)

## Contributing

Ideas, bug reports and pull requests are welcome. Open an issue to discuss a change before sending a large PR.

Only scan networks and hosts you own or have permission to test.

## License

MIT. See `LICENSE`.
