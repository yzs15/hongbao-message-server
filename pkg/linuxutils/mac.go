package linuxutils

import (
	"net"
)

func GetMac() (string, error) {
	netInterface, err := net.InterfaceByName("eth0")
	if err != nil {
		return "", err
	}
	macAddress := netInterface.HardwareAddr
	return macAddress.String(), nil
}
