package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

type MemoryInfo struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	SwapTotal uint64 `json:"swap_total"`
}

func CollectMemoryInfo() (*MemoryInfo, error) {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to collect memory stats: %w", err)
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to collect swap memory stats: %w", err)
	}

	stats := &MemoryInfo{
		Total:     BytesToMb(memoryInfo.Total),
		Available: BytesToMb(memoryInfo.Available),
		SwapTotal: BytesToMb(swapInfo.Total),
	}

	return stats, nil
}
