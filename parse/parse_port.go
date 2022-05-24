package parse

import (
	"nacs/common"
	"nacs/utils"
	"strconv"
	"strings"
)

func OnlyPorts(onlyPorts string) (scanPorts []int) {
	if onlyPorts == "" {
		return
	}
	slices := strings.Split(onlyPorts, ",")
	for _, port := range slices {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			startPort, _ := strconv.Atoi(ranges[0])
			endPort, _ := strconv.Atoi(ranges[1])
			if startPort < endPort {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}
		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			scanPorts = append(scanPorts, i)
		}
	}
	scanPorts = utils.RemoveDuplicateInt(scanPorts)
	return scanPorts
}

func AddPorts(addPorts string) (scanPorts []int) {
	if addPorts == "" {
		return
	}
	slices := strings.Split(addPorts, ",")
	for _, port := range slices {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			startPort, _ := strconv.Atoi(ranges[0])
			endPort, _ := strconv.Atoi(ranges[1])
			if startPort < endPort {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}
		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			scanPorts = append(scanPorts, i)
		}
	}
	for port := range common.DefaultPorts {
		scanPorts = append(scanPorts, port)
	}
	scanPorts = utils.RemoveDuplicateInt(scanPorts)
	return scanPorts
}
