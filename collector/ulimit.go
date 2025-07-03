package collector

import (
	"fmt"
	"regexp"
	"strings"
)

func CollectUlimitInfo() (map[string]string, error) {
	cmd, err := RunCmd("sh", "-c", "ulimit -a")

	if err != nil {
		return nil, fmt.Errorf("failed to run ulimit -a: %w", err)
	}

	lines := strings.Split(cmd, "\n")
	result := make(map[string]string)

	parse := func(line string) (string, string, bool) {
		// 匿名函数，作用域限定在这个函数，删除ulimit结果中括号及其内部内容，并用_替换空格
		// 格式可能是：core file size          (blocks, -c) unlimited
		re := regexp.MustCompile(`\([^)]*\)`)
		line = re.ReplaceAllString(line, "")

		// 去除多余空格
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return "", "", false
		}

		// key: 前面的字段，value: 最后一个字段
		key := strings.Join(fields[:len(fields)-1], " ")
		// 空格用_替换
		key = strings.ReplaceAll(key, " ", "_")
		value := fields[len(fields)-1]
		return key, value, true
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key, value, _ := parse(line)
		result[key] = value
	}

	return result, nil
}
