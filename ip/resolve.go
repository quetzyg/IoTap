package ip

import "net"

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

	var ips []net.IP
	for ip := address.Mask(network.Mask); network.Contains(ip); next(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}

	return ips, nil
}
