package collector

// 获取sysctl 所有字段

import (
	"maps"
	"strings"
)

func sysCtl() (map[string]string, error) {
	out, err := RunCmd("sysctl", "-a")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	result := make(map[string]string)
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result, nil
}

func CollectSysctlInfo() (map[string]string, error) {
	sysctlMap := make(map[string]string)
	sysctl, err := sysCtl()
	if err != nil {
		return nil, err
	}
	maps.Copy(sysctlMap, sysctl)
	return sysctlMap, nil
}
