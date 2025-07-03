package collector

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type GPUInfo struct {
	Vendor    string        `json:"vendor"`         // GPU 厂商
	Model     string        `json:"model"`          // GPU 型号
	DriverVer string        `json:"driver_version"` // 驱动版本
	CUDAVer   string        `json:"cuda_version"`   // CUDA 版本（仅 NVIDIA）
	Cards     []GPUCardInfo `json:"cards"`          // GPU 卡信息
	Count     int           `json:"count"`          // GPU 卡数量

}

type GPUCardInfo struct {
	ID          int16  `json:"id"`           // GPU 卡 ID
	UUID        string `json:"uuid"`         // GPU 卡唯一标识符
	MemoryTotal string `json:"memory_total"` // GPU 总显存
	MemoryFree  string `json:"memory_free"`  // GPU 可用显存
}

func getNvidiaGPUInfo() (*GPUInfo, error) {
	// 1. 获取 NVIDIA GPU 卡信息
	smiOut, err := RunCmd("nvidia-smi",
		"--query-gpu=index,uuid,memory.total,memory.free",
		"--format=csv,noheader,nounits")
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader([]byte(smiOut)))
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var cards []GPUCardInfo
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		index, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid GPU ID '%s': %w", line[0], err)
		}
		cards = append(cards, GPUCardInfo{
			ID:          int16(index),
			UUID:        strings.TrimSpace(line[1]),
			MemoryTotal: strings.TrimSpace(line[2]),
			MemoryFree:  strings.TrimSpace(line[3]),
		})
	}
	// 2. 获取 NVIDIA 驱动版本和显卡型号
	driverOut, err := RunCmd("nvidia-smi",
		"--query-gpu=name,driver_version",
		"--format=csv,noheader,nounits")
	if err != nil {
		return nil, err
	}

	driverLines := strings.Split(strings.TrimSpace(driverOut), "\n")
	if len(driverLines) == 0 {
		return nil, fmt.Errorf("empty driver output")
	}
	fields := strings.Split(driverLines[0], ",")
	if len(fields) < 2 {
		return nil, fmt.Errorf("unexpected driver output: %s", driverLines[0])
	}

	// 3. 获取 CUDA 版本
	fullSmiOut, err := RunCmd("nvidia-smi", "--version")
	if err != nil {
		return nil, err
	}
	cudaVer := "unknown"
	for _, line := range strings.Split(fullSmiOut, "\n") {
		if strings.Contains(line, "CUDA Version") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				cudaVer = strings.TrimSpace(parts[1])
			}
			break
		}
	}

	return &GPUInfo{
		Vendor:    "NVIDIA",
		Model:     strings.TrimSpace(fields[0]),
		DriverVer: strings.TrimSpace(fields[1]),
		CUDAVer:   cudaVer,
		Cards:     cards,
		Count:     len(cards),
	}, nil
}

func getHygonGPUInfo() (*GPUInfo, error) {
	return &GPUInfo{
		Vendor:    "Hygon",
		Model:     "Hygon GPU",
		DriverVer: "unknown",
		CUDAVer:   "unknown",
		Cards:     []GPUCardInfo{},
	}, nil
}

func CollectGPUInfo() (*GPUInfo, error) {
	lspciOut, err := RunCmd("lspci")
	if err != nil {
		return nil, err
	}
	lspciStr := strings.ToLower(lspciOut)

	switch {
	case strings.Contains(lspciStr, "nvidia"):
		return getNvidiaGPUInfo()
	case strings.Contains(lspciStr, "hygon"):
		return getHygonGPUInfo()
	default:
		return nil, fmt.Errorf("no supported GPU (NVIDIA/Hygon) found")
	}
}
