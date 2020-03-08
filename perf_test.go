package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

type D struct {
	p int32
	c float64
	m float32
}

var ferrs, perrs, serrs int

func BenchmarkProcsForking(b *testing.B) {
	ferrs = 0
	for n := 0; n < b.N; n++ {
		getProcsForking()
	}
	b.Logf("forking errors: %d", ferrs)
}

func BenchmarkProcsGoPSUtil(b *testing.B) {
	perrs = 0
	for n := 0; n < b.N; n++ {
		getProcsGoPsUtil()
	}
	b.Logf("GoPsUtil errors: %d", perrs)
}

func BenchmarkProcsGoSysInfo(b *testing.B) {
	serrs = 0
	for n := 0; n < b.N; n++ {
		getProcsGoSysInfo()
	}
	b.Logf("GoSysInfo errors: %d", serrs)
}

func BenchmarkTempGoPSUtil(b *testing.B) {
	perrs = 0
	for n := 0; n < b.N; n++ {
		getTempsGoPSUtil()
	}
	b.Logf("GoSysInfo errors: %d", serrs)
}

func BenchmarkTempNEW_ZELCH(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// your's here
	}
}

func getTempsGoPSUtil() {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		perrs += 1
	}
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			_ = sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
		}
	}
}

func getProcsForking() {
	output, err := exec.Command("ps", "-axo", "pid:10,comm:50,pcpu:5,pmem:5,args").Output()
	if err != nil {
		ferrs += 1
	}

	// converts to []string, removing trailing newline and header
	linesOfProcStrings := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")[1:]

	for _, line := range linesOfProcStrings {
		_, err = strconv.Atoi(strings.TrimSpace(line[0:10]))
		if err != nil {
			ferrs += 1
			continue
		}
		_, err = strconv.ParseFloat(strings.TrimSpace(line[63:68]), 64)
		if err != nil {
			ferrs += 1
			continue
		}
		_, err = strconv.ParseFloat(strings.TrimSpace(line[69:74]), 32)
		if err != nil {
			ferrs += 1
			continue
		}
	}
}

func getProcsGoPsUtil() {
	procs, err := process.Processes()
	if err != nil {
		perrs += 1
	}
	for _, p := range procs {
		_, err = p.CPUPercent()
		if err != nil {
			perrs += 1
			continue
		}
		_, err = p.MemoryPercent()
		if err != nil {
			perrs += 1
		}
	}
}

func getProcsGoSysInfo() {
	procs, err := sysinfo.Processes()
	if err != nil {
		serrs += 1
	}
	for _, p := range procs {
		_, err = p.CPUTime()
		if err != nil {
			serrs += 1
			continue
		}
		_, err = p.Memory()
		if err != nil {
			serrs += 1
		}
	}
}
