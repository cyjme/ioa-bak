package monitor

import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"strconv"
)

type sysInfo struct {
	System string `json:"system"`
	Cpu    string `json:"cpu"`
	Memory memory `json:"memory"`
}

type memory struct {
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	cpuPercent, _ := cpu.Percent(0, false)
	v, _ := mem.VirtualMemory()
	hostInfo, _ := host.Info()

	memory := memory{
		Used:  v.Used / 1024 / 1024,
		Total: v.Total / 1024 / 1024,
	}

	sysInfo := sysInfo{
		System: hostInfo.OS,
		Cpu:    strconv.FormatFloat(cpuPercent[0], 'f', 6, 64),
		Memory: memory,
	}

	response, _ := json.Marshal(sysInfo)

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
