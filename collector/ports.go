package collector

import (
	"regexp"
	"strings"
)

type PortUsage struct {
	Port    string `json:"port"`
	Process string `json:"process"` // 完整进程信息
	// pid     int    `json:"pid"`
	// Protocol string `json:"protocol"` // tcp/udp
	// IP       string `json:"ip"`       // 监听IP
}

type PortUsageGroup struct {
	TCP []PortUsage `json:"tcp"`
	UDP []PortUsage `json:"udp"`
}

func CollectPortsInfo() (*PortUsageGroup, error) {
	cmd, err := RunCmd("ss", "-tunlp")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(cmd, "\n")
	portRegex := regexp.MustCompile(`(?:[\d\.]+|\[.+\])(?:%\w+)?:(\d+)`)
	result := &PortUsageGroup{
		TCP: []PortUsage{},
		UDP: []PortUsage{},
	}

	// 用于去重
	seenPorts := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		protocol := strings.ToLower(fields[0])
		localAddr := fields[4]

		portMatch := portRegex.FindStringSubmatch(localAddr)
		if len(portMatch) < 2 {
			continue
		}
		port := portMatch[1]

		// 去重检查
		portKey := protocol + ":" + port
		if seenPorts[portKey] {
			continue
		}
		seenPorts[portKey] = true

		// 提取完整进程信息
		var processInfo string
		if usersIndex := strings.Index(line, "users:((\""); usersIndex > 0 {
			processInfo = line[usersIndex+len("users:(("):]
			processInfo = strings.TrimSuffix(processInfo, "))")
		}

		pu := PortUsage{
			Port:    port,
			Process: processInfo,
		}

		switch protocol {
		case "tcp", "tcp6":
			result.TCP = append(result.TCP, pu)
		case "udp", "udp6":
			result.UDP = append(result.UDP, pu)
		}
	}
	return result, nil
}
