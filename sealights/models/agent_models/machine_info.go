package agent_models

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

type MachineInfo struct {
	MachineName        string   `json:"machineName"`
	Arch               string   `json:"arch"`
	Os                 string   `json:"os"`
	LocalDateTime      string   `json:"localDateTime"`
	LocalDateTimeUnixS int64    `json:"localDateTimeUnix_s"`
	IpAddress          []string `json:"ipAddress"`
}

func NewMachineInfo() *MachineInfo {
	host, _ := os.Hostname()
	var ips []string
	addrs, err := net.LookupHost(host)
	if err == nil {
		for _, a := range addrs {
			ips = append(ips, a)
		}
	}
	mi := &MachineInfo{
		MachineName:        host,
		Arch:               runtime.GOARCH,
		Os:                 runtime.GOOS,
		LocalDateTime:      time.Now().Format(time.RFC3339),
		LocalDateTimeUnixS: time.Now().UnixNano() / int64(time.Millisecond),
		IpAddress:          ips,
	}
	return mi
}

func (i *MachineInfo) Validate() error {
	if i.MachineName == "" {
		return fmt.Errorf("MachineInfo: machineName is not valid")
	}
	if i.Arch == "" {
		return fmt.Errorf("MachineInfo: arch is not valid")
	}
	if i.Os == "" {
		return fmt.Errorf("MachineInfo: os is not valid")
	}
	if i.LocalDateTime == "" {
		return fmt.Errorf("MachineInfo: localDateTime is not valid")
	}
	if i.LocalDateTimeUnixS == 0 {
		return fmt.Errorf("MachineInfo: localDateTimeUnix_s is not valid")
	}
	return nil
}

var _ models.Model = &MachineInfo{}
