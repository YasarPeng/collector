module collector

go 1.23.1

require (
	github.com/astaxie/beego v1.12.3
	github.com/integrii/flaggy v1.5.2
	github.com/moby/sys/mountinfo v0.7.2
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/shirou/gopsutil/v3 v3.24.5
	github.com/spf13/pflag v1.0.6
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace precheck/collector => ../collector
