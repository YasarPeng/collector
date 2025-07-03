package collector

import (
	"net"
	"strings"

	psnet "github.com/shirou/gopsutil/v3/net"
)

type NetInterfaceStats struct {
	Name string `json:"name"`
	// MTU          int      `json:"mtu"`
	// HardwareAddr string   `json:"mac_address"`
	IPAddrs   []string `json:"ip_addresses"`
	Flags     []string `json:"flags"`
	BytesSent uint64   `json:"bytes_sent"`
	BytesRecv uint64   `json:"bytes_recv"`
}

var virtualPrefixes = []string{
	"lo",
	"br-",
	"veth",
	"virbr",
	"vmnet",
	"tun",
	"utun",
	"tap",
	"flannel",
	"wg",
	"kube",
	"zt",
	"tailscale",
}

func isVirtualInterface(name string) bool {
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func flagsToStrings(f net.Flags) []string {
	var flags []string
	if f&net.FlagUp != 0 {
		flags = append(flags, "up")
	}
	if f&net.FlagLoopback != 0 {
		flags = append(flags, "loopback")
	}
	if f&net.FlagBroadcast != 0 {
		flags = append(flags, "broadcast")
	}
	if f&net.FlagMulticast != 0 {
		flags = append(flags, "multicast")
	}
	return flags
}

func CollectNetworkInfo() ([]NetInterfaceStats, error) {
	var result []NetInterfaceStats

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ioCounters, err := psnet.IOCounters(true)
	if err != nil {
		return nil, err
	}
	ioMap := make(map[string]psnet.IOCountersStat)
	for _, counter := range ioCounters {
		ioMap[counter.Name] = counter
	}

	for _, iface := range interfaces {
		// 仅启用状态且非虚拟网卡
		if iface.Flags&net.FlagUp == 0 || isVirtualInterface(iface.Name) {
			continue
		}

		stats := NetInterfaceStats{
			Name: iface.Name,
			// MTU:          iface.MTU,
			// HardwareAddr: iface.HardwareAddr.String(),
			Flags: flagsToStrings(iface.Flags),
		}

		// 获取 IP 地址
		if byName, err := net.InterfaceByName(iface.Name); err == nil {
			if addrs, err := byName.Addrs(); err == nil {
				for _, addr := range addrs {
					stats.IPAddrs = append(stats.IPAddrs, addr.String())
				}
			}
		}

		// 加载流量数据
		if counter, ok := ioMap[iface.Name]; ok {
			stats.BytesSent = counter.BytesSent
			stats.BytesRecv = counter.BytesRecv
		}

		result = append(result, stats)
	}
	return result, nil
}
