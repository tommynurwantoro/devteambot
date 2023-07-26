package ipaddress

import "net"

// ServerIP gets server ip
// TODO: rewrite to singleton
func ServerIP() string {
	var ip string
	address, err := net.InterfaceAddrs()
	if err != nil {
		address = []net.Addr{}
	}

	for _, a := range address {
		if ipNET, ok := a.(*net.IPNet); ok && !ipNET.IP.IsLoopback() {
			if ipNET.IP.To4() != nil {
				ip = ipNET.IP.String()
			}
		}
	}
	return ip
}
