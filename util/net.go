package util

import (
	"github.com/misakacoder/logger"
	"haibara/consts"
	"net"
)

func GetLocalAddr() (addr []string) {
	addr = append(addr, consts.Localhost, "127.0.0.1")
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, value := range interfaces {
		if (value.Flags & net.FlagUp) != 0 {
			addresses, _ := value.Addrs()
			for _, address := range addresses {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						addr = append(addr, ipNet.IP.String())
					}
				}
			}
		}
	}
	return
}
