package ip

import (
	"errors"
	"fmt"
	"net"
)

var errNetworkMembership = errors.New("you must be in the same network")

// inNetwork checks if the caller belongs to a network.
func inNetwork(network *net.IPNet) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if v, ok := addr.(*net.IPNet); ok && network.Contains(v.IP) {
				return nil
			}
		}
	}

	return fmt.Errorf("%w: %s", errNetworkMembership, network.IP)
}

// next IP address.
func next(addr net.IP) {
	for i := len(addr) - 1; i >= 0; i-- {
		addr[i]++
		if addr[i] > 0 {
			break
		}
	}
}

// Resolve the IP addresses of a given CIDR.
func Resolve(cidr string) ([]net.IP, error) {
	address, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	if err = inNetwork(network); err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := address.Mask(network.Mask); network.Contains(ip); next(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}

	return ips, nil
}
