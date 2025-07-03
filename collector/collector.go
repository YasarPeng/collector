package collector

import (
	"encoding/json"
	"fmt"
	"os"

	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/integrii/flaggy"
	"gopkg.in/yaml.v3"
)

var (
	BuildVersion = "0.0.1-dev"
)

func GetCollectors() map[string]func() (any, error) {
	return map[string]func() (any, error){
		"cpu": func() (any, error) {
			return CollectCPUInfo()
		},
		"memory": func() (any, error) {
			return CollectMemoryInfo()
		},
		"disk": func() (any, error) {
			return CollectDiskInfo()
		},
		"gpu": func() (any, error) {
			return CollectGPUInfo()
		},
		"network": func() (any, error) {
			return CollectNetworkInfo()
		},
		"os": func() (any, error) {
			return CollectOsInfo()
		},
		"command": func() (any, error) {
			return CollectCommonCmdInfo()
		},
		"port": func() (any, error) {
			return CollectPortsInfo()
		},
		"sysctl": func() (any, error) {
			return CollectSysctlInfo()
		},
		"ulimit": func() (any, error) {
			return CollectUlimitInfo()
		},
	}
}

// 判断是否需要采集该字段
func ShouldCollect(name string, filters map[string]bool) bool {
	return len(filters) == 0 || filters[name]
}

func must[T any](val T, err error) T {
	if err != nil && DebugMode {
		logs.Warning("Collection failed: %v", err)
	}
	return val
}

func Collector() {
	type FlagsConfig struct {
		Output     string
		Debug      bool
		List       bool
		FilterList string
	}
	cfg := &FlagsConfig{
		Output:     "json",
		Debug:      false,
		List:       false,
		FilterList: "",
	}
	flaggy.String(&cfg.Output, "o", "output", "Output format: json or yaml")
	flaggy.Bool(&cfg.Debug, "d", "debug", "Enable debug logging")
	flaggy.Bool(&cfg.List, "l", "list", "List all supported collector fields")
	flaggy.String(&cfg.FilterList, "f", "filter", "Only collect specific info (e.g. all)")
	flaggy.SetDescription("System Resource Collector")
	flaggy.SetVersion(BuildVersion)
	// flaggy.SetName("collector")
	flaggy.Parse()

	if cfg.Debug {
		DebugMode = true
		logs.Info("Debug mode enabled")
	}

	if cfg.List {
		for k := range GetCollectors() {
			fmt.Println("- " + k)
		}
		os.Exit(0)
	}

	// 根据 filter 控制采集行为
	collectFields := map[string]bool{}
	if cfg.FilterList != "" {
		for _, field := range strings.Split(cfg.FilterList, ",") {
			collectFields[strings.TrimSpace(field)] = true
		}
	}

	fullData := map[string]any{}

	collectors := GetCollectors()

	for key, fn := range collectors {
		if ShouldCollect(key, collectFields) {
			fullData[key] = must(fn())
		}
	}

	var outputData any = fullData

	// 输出
	switch strings.ToLower(cfg.Output) {
	case "yaml":
		yamlData, err := yaml.Marshal(outputData)
		if err != nil {
			logs.Error("YAML encoding failed: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(yamlData))
	case "json":
		jsonData, err := json.MarshalIndent(outputData, "", "  ")
		if err != nil {
			logs.Error("JSON encoding failed: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	default:
		logs.Error("Unsupported output format: %s", cfg.Output)
		os.Exit(1)
	}
}
