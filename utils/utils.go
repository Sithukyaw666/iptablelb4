package utils

import (
	"fmt"
	"math"
	"net"
)

func GenerateIptablerules(index, length int, localIp, dip, dport, algorithm string) ([]string, []string) {

	destination := fmt.Sprintf("%s:%s", dip, dport)

	traffic := fmt.Sprintf("%v", length-index)

	probability := fmt.Sprintf("%v", (math.Floor(float64(100)/float64(length-index)))/100)
	ingress := []string{}
	egress := []string{
		"-d", dip,
		"-p", "tcp",
		"-m", "tcp",
		"--dport", dport,
		"-j", "SNAT",
		"--to-source", localIp,
	}

	if algorithm == "round-robin" {
		ingress = []string{
			"-p", "tcp",
			"--match", "statistic",
			"--mode", "nth",
			"--every", traffic,
			"--dport", dport,
			"--packet", "0",
			"-j", "DNAT",
			"--to-destination", destination,
		}

	} else {
		ingress = []string{
			"-p", "tcp",
			"--match", "statistic",
			"--mode", "random",
			"--probability", probability,
			"--dport", "80",
			"-j", "DNAT",
			"--to-destination", destination,
		}
	}

	return ingress, egress

}

func GetLocalIPs() ([]string, error) {
	var ips []string
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// Convert net.IP to string and append
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func IsPredefinedChain(chain string) bool {
	predefinedChains := []string{"INPUT", "OUTPUT", "PREROUTING", "POSTROUTING", "DOCKER"}
	for _, pChain := range predefinedChains {
		if chain == pChain {
			return true
		}
	}
	return false
}
