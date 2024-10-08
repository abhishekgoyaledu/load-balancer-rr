package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	MinPort = 0
	MaxPort = 65535
)

// ValidatePort checks if the provided port number is within the valid range (1-65535) and not currently in use.
var ValidatePort = func(port int) error {
	// Check if the port is within the valid range
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	occupiedPorts, err := getOccupiedPorts()
	if err != nil {
		return fmt.Errorf("failed to get occupied ports: %v", err)
	}

	for _, occupiedPort := range occupiedPorts {
		if port == occupiedPort {
			return errors.New("port is already in use")
		}
	}

	return nil
}

// getOccupiedPorts returns a slice of occupied port numbers.
func getOccupiedPorts() ([]int, error) {
	// Execute the lsof command to list all listening TCP ports
	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Use a regular expression to find all port numbers in the command output
	re := regexp.MustCompile(`:(\d+)\s`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	ports := []int{}
	for _, match := range matches {
		// Convert the matched port number from string to int and add to the list
		if len(match) > 1 {
			port, err := strconv.Atoi(match[1])
			if err == nil {
				ports = append(ports, port)
			}
		}
	}

	return ports, nil
}
