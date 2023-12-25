package ip

import (
	"errors"
	"fmt"
	"net"
)

var (
	errNetworkCannotBeNil = errors.New("the network cannot be nil")
	errNetworkMembership  = errors.New("you must be in the same network")
)

// validateNetworkMembership of the caller.
func validateNetworkMembership(network *net.IPNet) error {
	if network == nil {
		return errNetworkCannotBeNil
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}

		var ip net.IP

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if network.Contains(ip) {
				return nil
			}
		}
	}

	return fmt.Errorf("%w: %s", errNetworkMembership, network.IP)
}

// Resolve the IP addresses of a given CIDR.
func Resolve(cidr string) ([]net.IP, error) {
	address, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	err = validateNetworkMembership(network)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := address.Mask(network.Mask); network.Contains(ip); next(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}

	return ips, nil
}

// next IP address.
func next(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
