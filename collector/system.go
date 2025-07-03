package collector

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
)

type SystemInfo struct {
	Firewall    string   `json:"firewall"`
	SELinux     string   `json:"selinux"`
	OSName      string   `json:"os_name"`
	OSVersion   string   `json:"os_version"`
	Kernel      string   `json:"kernel"`
	CurrentTime string   `json:"current_time"`
	Hostname    string   `json:"hostname"`
	Arch        string   `json:"arch"`
	IPAddrs     []string `json:"ip_addrs"`
}

func CollectOsInfo() (*SystemInfo, error) {
	info, _ := host.Info()

	// 1. firewall状态
	firewallStatus := "inactive"
	switch {
	case strings.Contains(strings.ToLower(info.Platform), "centos"),
		strings.Contains(strings.ToLower(info.Platform), "rhel"):
		if out, err := RunCmd("systemctl", "is-active", "firewalld"); err == nil && strings.TrimSpace(out) == "active" {
			firewallStatus = "active"
		}

	case strings.Contains(strings.ToLower(info.Platform), "ubuntu"),
		strings.Contains(strings.ToLower(info.Platform), "debian"):
		if out, err := RunCmd("systemctl", "is-active", "ufw"); err == nil && strings.TrimSpace(out) == "active" {
			firewallStatus = "active"
		}
	}
	// 2. selinux 状态
	selinuxStatus := "disabled"
	if out, err := RunCmd("getenforce"); err == nil && strings.TrimSpace(out) != "Disabled" {
		selinuxStatus = "enabled"
	} else if out, err := RunCmd("cat", "/etc/selinux/config"); err == nil {
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(strings.ToLower(line), "selinux=") &&
				!strings.Contains(strings.ToLower(line), "disabled") {
				selinuxStatus = "enabled"
				break
			}
		}
	}

	// 3. 系统版本
	osVersion, osName := info.PlatformVersion, info.Platform

	// 5. 内核版本
	kernelRegex, _ := regexp.Compile(`v?(\d+\.\d+\.\d+)`)
	kernel := kernelRegex.FindString(info.KernelVersion)

	// 6. 主机名
	hostname := info.Hostname

	// 7. CPU架构
	arch := info.KernelArch

	// 8. IP 地址
	var ips []string
	excludePrefixes := []string{"docker", "br-", "veth", "lo", "tun", "virbr", "flannel", "cni", "nerdctl", "kube-"}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// 排除虚拟/容器接口
		skip := false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(iface.Name, prefix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 只保留IPv4内网地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			// 检查是否为私有地址
			if ip.IsPrivate() {
				ips = append(ips, ip.String())
			}
		}
	}

	return &SystemInfo{
		Firewall:    firewallStatus,
		SELinux:     selinuxStatus,
		OSName:      osName,
		OSVersion:   osVersion,
		Kernel:      kernel,
		CurrentTime: fmt.Sprintf("%d", time.Now().Unix()),
		Hostname:    hostname,
		Arch:        arch,
		IPAddrs:     ips,
	}, nil
}
