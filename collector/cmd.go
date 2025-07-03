package collector

import (
	"os/exec"
	"regexp"
	"strings"
)

// 判断命令是否存在，且输出版本

type CommonCmdInfo struct {
	Name    string `json:"name"`
	Output  string `json:"output"`
	Exist   bool   `json:"exist"`
	Version string `json:"version"`
}

func CollectCommonCmdInfo() ([]CommonCmdInfo, error) {
	commandMap := map[string][]string{
		"docker":  {"docker", "--version"},
		"nerdctl": {"nerdctl", "version", "-f", "{{.Client.Version}}"},
		"helm":    {"helm", "version", "--short"},
		"kubelet": {"kubelet", "--version"},
		"kubeadm": {"kubeadm", "version", "-o", "short"},
	}

	compileMap := map[string]string{
		"docker":  `v?(\d+\.\d+\.\d+)`,
		"nerdctl": `v?(\d+\.\d+\.\d+)`,
		"helm":    `v?(\d+\.\d+\.\d+)`,
		"kubelet": `v?(\d+\.\d+\.\d+)`,
		"kubeadm": `v?(\d+\.\d+\.\d+)`,
	}

	var results []CommonCmdInfo

	for name, args := range commandMap {
		_, err := exec.LookPath(args[0])

		info := CommonCmdInfo{
			Name:  name,
			Exist: err == nil,
		}
		// 1. 先判断命令是否存在
		if err != nil {
			info.Output = "not found"
			info.Version = ""
			results = append(results, info)
			continue
		}

		// 2. 如果命令存在，执行命令
		cmd, err := RunCmd(args[0], args[1:]...)
		// var out bytes.Buffer
		if err != nil {
			info.Output = strings.TrimSpace(cmd)
			info.Version = ""
			results = append(results, info)
			continue
		}
		versionRegex := regexp.MustCompile(compileMap[name])
		match := versionRegex.FindString(cmd)
		output := strings.TrimSpace(cmd)
		info.Output = output
		info.Version = match

		results = append(results, info)
	}
	return results, nil
}
