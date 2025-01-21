package utils

import (
	"fmt"
	"math"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

func GenerateIptablerules(index, length int, dip, dport, algorithm, port string) ([]string, []string) {

	destination := fmt.Sprintf("%s:%s", dip, dport)

	traffic := fmt.Sprintf("%v", length-index)

	probability := fmt.Sprintf("%v", (math.Floor(float64(100)/float64(length-index)))/100)
	ingress := []string{}
	egress := []string{
		"-d", dip,
		"-p", "tcp",
		"-m", "tcp",
		"--dport", dport,
		"-j", "MASQUERADE",
	}

	if algorithm == "round-robin" {
		ingress = []string{
			"-p", "tcp",
			"--match", "statistic",
			"--mode", "nth",
			"--every", traffic,
			"--dport", port,
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
			"--dport", port,
			"-j", "DNAT",
			"--to-destination", destination,
		}
	}

	return ingress, egress

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

func ExtractModeAndDestination(input string) (string, string, string, string, error) {
	// Regular expression to match --mode
	modeRegex := regexp.MustCompile(`--mode\s+(\w+)`)
	modeMatch := modeRegex.FindStringSubmatch(input)
	if len(modeMatch) < 2 {
		return "", "", "", "", fmt.Errorf("mode not found")
	}

	// Regular expression to match --to-destination IP and port
	destRegex := regexp.MustCompile(`--to-destination\s+(\d+\.\d+\.\d+\.\d+):(\d+)`)
	destMatch := destRegex.FindStringSubmatch(input)
	if len(destMatch) < 3 {
		return "", "", "", "", fmt.Errorf("destination not found")
	}

	listenRegex := regexp.MustCompile(`--dport\s+(\d+)`)
	listenMatch := listenRegex.FindStringSubmatch(input)

	if len(listenMatch) < 2 {
		return "", "", "", "", fmt.Errorf("destination not found")
	}

	// Return the extracted values
	return modeMatch[1], listenMatch[1], destMatch[1], destMatch[2], nil
}
func EnableIPForwarding() error {
	filePath := "/proc/sys/net/ipv4/ip_forward"

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open %s for writing: %v", filePath, err)
	}
	defer file.Close()

	_, err = file.WriteString("1\n")
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", filePath, err)
	}

	return nil
}
func StandardResponse(c *gin.Context, statusCode int, status string, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
