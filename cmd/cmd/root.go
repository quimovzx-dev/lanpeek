// Package cmd defines lanpeek's command-line interface.
package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/quimovzx-dev/lanpeek/internal/network"
	"github.com/quimovzx-dev/lanpeek/internal/output"
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	var jsonOutput bool
	var timeout time.Duration
	var concurrency int
	root := &cobra.Command{
		Use: "lanpeek", Short: "An unprivileged LAN network toolkit",
		Long:         "Inspect interfaces, discover LAN hosts and check TCP ports without elevated privileges.\nOnly scan networks you own or have permission to test.\nTCP connect probes only: no raw packets, ICMP, ARP scanning or packet capture.",
		SilenceUsage: true, SilenceErrors: true,
	}
	root.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Write JSON only to stdout (progress goes to stderr)")
	root.PersistentFlags().DurationVar(&timeout, "timeout", 500*time.Millisecond, "Per-connect and reverse-DNS timeout (maximum 30s)")
	root.PersistentFlags().IntVar(&concurrency, "concurrency", 32, "Maximum parallel workers (1-256)")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if timeout <= 0 || timeout > 30*time.Second {
			return fmt.Errorf("timeout must be greater than 0 and at most 30s")
		}
		if concurrency < 1 || concurrency > 256 {
			return fmt.Errorf("concurrency must be between 1 and 256")
		}
		return nil
	}
	info := &cobra.Command{Use: "info", Short: "Show interfaces, addresses, masks and local MAC addresses", Args: cobra.NoArgs,
		Long: "Show interface up/down flags, IPs, masks and local MAC addresses.\nIPv4 default gateways are best-effort on Linux; omitted if unavailable.",
		RunE: func(cmd *cobra.Command, args []string) error {
			rows, err := network.Info()
			if err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(cmd.OutOrStdout(), rows)
			}
			var cells [][]string
			for _, row := range rows {
				var addrs []string
				for _, a := range row.Addresses {
					addrs = append(addrs, fmt.Sprintf("%s/%d (mask %s)", a.IP, a.Prefix, a.Mask))
				}
				cells = append(cells, []string{row.Name, row.State, blank(row.MAC), strings.Join(addrs, "\n"), blank(row.Gateway)})
			}
			return output.Table(cmd.OutOrStdout(), []string{"INTERFACE", "STATE", "LOCAL MAC", "ADDRESS / MASK", "GATEWAY"}, cells)
		},
	}
	var cidr string
	var noDNS bool
	discover := &cobra.Command{Use: "discover --cidr <IPv4-prefix>", Short: "Find responsive hosts with common TCP connect probes", Args: cobra.NoArgs,
		Long: "Probe TCP ports 22, 80, 443, 445 and 3389, stopping at the first open or refused reply.\nOpen and refused replies both mark a host reachable. Unresponsive does not mean offline: firewalls can hide hosts.\nOnly reachable hosts get best-effort reverse DNS. IPv4 only, at most 4096 hosts.\nOnly scan networks you own or have permission to test.",
		RunE: func(cmd *cobra.Command, args []string) error {
			hosts, err := network.Hosts(cidr)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Probing %d hosts with %d workers...\n", len(hosts), concurrency)
			rows, err := network.Discover(cmd.Context(), hosts, timeout, concurrency, !noDNS)
			if err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(cmd.OutOrStdout(), rows)
			}
			var cells [][]string
			for _, row := range rows {
				cells = append(cells, []string{row.IP, row.Status, blank(row.Hostname)})
			}
			return output.Table(cmd.OutOrStdout(), []string{"IP", "STATUS", "HOSTNAME"}, cells)
		},
	}
	discover.Flags().StringVar(&cidr, "cidr", "", "IPv4 network to scan (required)")
	discover.Flags().BoolVar(&noDNS, "no-dns", false, "Skip reverse DNS lookups")
	_ = discover.MarkFlagRequired("cidr")
	var portList string
	ports := &cobra.Command{Use: "ports <host>", Short: "Check TCP ports on one IP or hostname", Args: cobra.ExactArgs(1),
		Long: "Show TCP open, refused, timeout or error status for each port. Service names are conventional hints, not service detection.\nOnly scan networks you own or have permission to test.",
		RunE: func(cmd *cobra.Command, args []string) error {
			host := strings.TrimSpace(args[0])
			if host == "" || strings.ContainsAny(host, "/ \t\r\n") {
				return fmt.Errorf("host must be an IP or hostname, not a URL")
			}
			// Accept bracketed IPv6 as well as bare IPv6, but not host:port.
			if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
				host = strings.Trim(host, "[]")
			}
			if strings.Contains(host, ":") && net.ParseIP(host) == nil && !strings.Contains(host, "%") {
				return fmt.Errorf("supply only a host, with ports in --ports")
			}
			numbers, err := network.ParsePorts(portList)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Checking %d TCP ports on %s...\n", len(numbers), host)
			rows, err := network.ScanPorts(cmd.Context(), host, numbers, timeout, concurrency)
			if err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(cmd.OutOrStdout(), rows)
			}
			var cells [][]string
			for _, row := range rows {
				cells = append(cells, []string{strconv.Itoa(row.Port), row.Service, row.Status, blank(row.Detail)})
			}
			return output.Table(cmd.OutOrStdout(), []string{"PORT", "SERVICE HINT", "STATUS", "DETAIL"}, cells)
		},
	}
	ports.Flags().StringVar(&portList, "ports", "22,80,443", "Comma-separated TCP ports (1-65535)")
	root.AddCommand(info, discover, ports)
	return root
}

func blank(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
