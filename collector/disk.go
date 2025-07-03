package collector

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// 排除掉k8s、docker等容器相关的挂载点，不然显示居多
func isIgnoredMount(fstype string, mount string) bool {
	ignoredFSTypes := map[string]bool{
		"tmpfs":    true,
		"devtmpfs": true,
		"overlay":  true,
		"shm":      true,
		"squashfs": true,
		"proc":     true,
		"sysfs":    true,
		"cgroup":   true,
		"devfs":    true,
	}

	ignoredMountPrefixes := []string{
		"/dev",
		"/run",
		"/var/run",
		"/boot",
		"/var/lib/docker",
		"/var/lib/kubelet",
		"/var/lib/containerd",
	}

	if ignoredFSTypes[strings.ToLower(fstype)] {
		return true
	}
	for _, prefix := range ignoredMountPrefixes {
		if strings.HasPrefix(mount, prefix) {
			return true
		}
	}
	return false
}

func getDiskType(device string) string {
	base := filepath.Base(device)
	if strings.HasPrefix(base, "sd") || strings.HasPrefix(base, "hd") {
		// 去掉分区号，如 sda1 -> sda
		for i := len(base) - 1; i >= 0; i-- {
			if base[i] < '0' || base[i] > '9' {
				base = base[:i+1]
				break
			}
		}
	}
	rotationalPath := "/sys/block/" + base + "/queue/rotational"
	data, err := os.ReadFile(rotationalPath)
	if err != nil {
		return "unknown"
	}
	if strings.TrimSpace(string(data)) == "0" {
		return "ssd"
	}
	return "hdd"
}

// 获取指定路径的磁盘使用情况,主要检查部署目录
func GetDiskUsage(path string) (*disk.UsageStat, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		log.Println("failed to get usage for", path, ":", err)
		return nil, err
	}
	return usage, nil
}

func CollectDiskInfo() ([]map[string]any, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskInfo []map[string]any
	for _, part := range partitions {
		if isIgnoredMount(part.Fstype, part.Mountpoint) {
			continue
		}

		usage, err := GetDiskUsage(part.Mountpoint)
		if err != nil {
			continue
		}

		diskType := getDiskType(part.Device)

		diskInfo = append(diskInfo, map[string]any{
			"device":       part.Device,
			"mount_path":   usage.Path,
			"disk_type":    diskType,
			"total":        BytesToMb(usage.Total),
			"free":         BytesToMb(usage.Free),
			"used":         BytesToMb(usage.Used),
			"fstype":       part.Fstype,
			"inodes_total": usage.InodesTotal,
			"inodes_used":  usage.InodesUsed,
			"inodes_free":  usage.InodesFree,
			"free_percent": strconv.FormatFloat(100-usage.UsedPercent, 'f', 0, 32),
		})
	}
	return diskInfo, nil
}
