package collector

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
)

type CPUInfo struct {
	Cores int      `json:"cores"`
	Model string   `json:"model"`
	Hz    uint16   `json:"hz"`
	Flags []string `json:"flags"`
}

func CollectCPUInfo() (*CPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	// 使用第一个 CPU 信息作为代表
	info := cpuInfo[0]

	stats := &CPUInfo{
		Cores: runtime.NumCPU(),
		Model: info.ModelName,
		Hz:    uint16(info.Mhz),
		Flags: info.Flags,
	}

	return stats, nil
}
