package network

import (
	"log"
	"net"
)

const cloudflare = "1.1.1.1:80"

// Address of the running system.
func Address() net.IP {
	con, err := net.Dial("udp", cloudflare)
	if err != nil {
		log.Fatalf("unable to dial: %v", err)
	}

	defer func(con net.Conn) {
		err = con.Close()
		if err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}(con)

	return con.LocalAddr().(*net.UDPAddr).IP
}
