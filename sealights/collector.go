package __sealights__

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"testing"

	"github.com/konflux-ci/release-service/__sealights__/config"
	"github.com/konflux-ci/release-service/__sealights__/internal"
	"github.com/konflux-ci/release-service/__sealights__/pkg/clock"
)

var SealightPkg embed.FS

const (
	sealightsDisableEnvVar = "SEALIGHTS_DISABLE"
	sealightsDisableInit   = "SEALIGHTS_DISABLE_INIT"
	sealightsDisableOnInit = "SEALIGHTS_DISABLE_ON_INIT"
)

var (
	sealightsCollector collector
	isReady            uint32
	instrumentMap      map[string]string
	sealightsAgent     *internal.Agent
	colStats           = newCollectorStats()
)
var IsCLI = false

type collector interface {
	RecordTestEvent(event internal.Event)
	CheckSkipTest(funcName string) bool
	RecordFootprint(event internal.Event)
	RecordControlEvent(event internal.Event)
	AddActiveTest(testName string)
	RemoveActiveTest(testName string)
	GetActiveTests() []string
	GetFootprintsTimestamps() (int64, int64)
	Start() error
	Stop()
	SetOnStopFunc(func())
	SetCollectorStatsFunc(fn func() string)
}

func setIsReady(ready bool) {
	if ready {
		atomic.StoreUint32(&isReady, 1)
	} else {
		atomic.StoreUint32(&isReady, 0)
	}
}

func getIsReady() bool {
	return atomic.LoadUint32(&isReady) == 1
}

func onStop() {
	setIsReady(false)
	sealightsCollector = nil
}

func TraceFunc(id string) {
	if getIsReady() {
		traceFunc(id)
	}
}

func traceFunc(id string) {
	funcName, found := instrumentMap[id]
	if !found {
		log.Printf("TraceFunc: id %s not found\n", id)
		return
	}
	sealightsCollector.RecordFootprint(internal.NewRawEvent().
		SetEventType(internal.EventTypeTrace).
		SetFunction(funcName).
		SetActiveTestList(sealightsCollector.GetActiveTests()).
		SetFootprintsTimestamps(sealightsCollector.GetFootprintsTimestamps()).
		SetTimestamp(clock.UnixMilli()))
}

func StartMainFunc(id string) {
	if getIsReady() {
		startMainFunc(id)
	} else {
		colStats.setMainNotReady()
	}

}
func startMainFunc(id string) {
	colStats.setMainReady()
	funcName, found := instrumentMap[id]
	if !found {
		colStats.setMainNotInMap()
		log.Printf("StartMainFunc: id %s not found\n", id)
		return
	}
	sealightsCollector.RecordControlEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeStartMain).
			SetFunction(funcName).
			SetActiveTestList(sealightsCollector.GetActiveTests()).
			SetFootprintsTimestamps(sealightsCollector.GetFootprintsTimestamps()).
			SetTimestamp(clock.UnixMilli()))
}

func EndMainFunc(id string) {
	if getIsReady() {
		endMainFunc(id)
	}
}

func endMainFunc(id string) {
	funcName, found := instrumentMap[id]
	if !found {
		funcName = "main"
	}
	sealightsCollector.RecordControlEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeEndMain).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()))

	sealightsCollector.Stop()
}
func StartTestFunc(id string, t *testing.T) {
	if getIsReady() {
		startTestFunc(id, t)
	}
}

func startTestFunc(id string, t *testing.T) {
	funcName, found := instrumentMap[id]
	if !found {
		log.Printf("StartTestFunc: id %s not found\n", id)
		return
	}
	if sealightsCollector.CheckSkipTest(funcName) {
		ts := clock.UnixMilli()
		sealightsCollector.RecordTestEvent(
			internal.NewRawEvent().
				SetEventType(internal.EventTypeStartTest).
				SetFunction(funcName).
				SetTimestamp(ts))
		sealightsCollector.RecordTestEvent(
			internal.NewRawEvent().
				SetEventType(internal.EventTypeEndTest).
				SetFunction(funcName).
				SetTimestamp(ts+1).
				SetAttribute("result", "true").
				SetAttribute("skipped", "true"))
		t.Skipf("Sealights TIA: Skipping test %s", t.Name())
	} else {
		sealightsCollector.RecordTestEvent(
			internal.NewRawEvent().
				SetEventType(internal.EventTypeStartTest).
				SetFunction(funcName).
				SetTimestamp(clock.UnixMilli()))
		sealightsCollector.AddActiveTest(funcName)
	}
}

func startGinkgoTestFunc(funcName string) bool {
	if sealightsCollector.CheckSkipTest(funcName) {
		ts := clock.UnixMilli()
		sealightsCollector.RecordTestEvent(
			internal.NewRawEvent().
				SetEventType(internal.EventTypeStartTest).
				SetFunction(funcName).
				SetTimestamp(ts))
		sealightsCollector.RecordTestEvent(
			internal.NewRawEvent().
				SetEventType(internal.EventTypeEndTest).
				SetFunction(funcName).
				SetTimestamp(ts+1).
				SetAttribute("result", "true").
				SetAttribute("skipped", "true"))
		return true
	}
	sealightsCollector.RecordTestEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeStartTest).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()))
	sealightsCollector.AddActiveTest(funcName)
	return false
}

func endGinkgoTestFunc(funcName string, result bool, skipped bool) {
	sealightsCollector.RecordTestEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeEndTest).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()).
			SetAttribute("result", fmt.Sprintf("%t", result)).
			SetAttribute("skipped", fmt.Sprintf("%t", skipped)))
	sealightsCollector.RemoveActiveTest(funcName)
}
func EndTestFunc(id string, t *testing.T) {
	if getIsReady() {
		endTestFunc(id, t)
	}
}

func endTestFunc(id string, t *testing.T) {
	funcName, found := instrumentMap[id]
	if !found {
		log.Printf("EndTestFunc: id %s not found\n", id)
		return
	}
	sealightsCollector.RecordTestEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeEndTest).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()).
			SetAttribute("result", fmt.Sprintf("%t", !t.Failed())).
			SetAttribute("skipped", fmt.Sprintf("%t", t.Skipped())))
	sealightsCollector.RemoveActiveTest(funcName)
}
func StartTestMainFunc(id string) {
	if getIsReady() {
		startTestMainFunc(id)
	}
}

func startTestMainFunc(id string) {
	funcName, found := instrumentMap[id]
	if !found {
		log.Printf("StartTestMainFunc: id %s not found\n", id)
		return
	}
	sealightsCollector.RecordControlEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeStartTestMain).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()))
}
func RunMainTestFunc(id string, m *testing.M) int {
	code := m.Run()
	EndTestMainFunc(id)
	return code
}

func EndTestMainFunc(id string) {
	if getIsReady() {
		endTestMainFunc(id)
	}

}

func endTestMainFunc(id string) {
	funcName, found := instrumentMap[id]
	if !found {
		log.Printf("EndTestMainFunc: id %s not found\n", id)
		return
	}
	sealightsCollector.RecordControlEvent(
		internal.NewRawEvent().
			SetEventType(internal.EventTypeEndTestMain).
			SetFunction(funcName).
			SetTimestamp(clock.UnixMilli()))
	sealightsCollector.Stop()
}
func init() {
	if os.Getenv(sealightsDisableInit) == "true" {
		log.Printf("Sealights agent is in testing mode, skipping init")
		return
	}
	if initCollector() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			setIsReady(false)
			if sealightsCollector != nil {
				sealightsCollector.Stop()
			}
		}()
		setIsReady(true)
	}

}

func initCollector() bool {
	if os.Getenv("SL_MODE") == "cli" {
		return false
	}
	if os.Getenv(sealightsDisableEnvVar) == "true" {
		log.Printf("Sealights agent is disabled")
		return false
	}
	cfg := config.NewConfig()
	if err := cfg.Load(ConfigData); err != nil {
		log.Printf("Failed to load Sealights build scan configuration file: %s, Sealights agent is disabled", err.Error())
		return false
	}
	if cfg.DisableSealightsOnInit {
		val := os.Getenv(sealightsDisableOnInit)
		if val == "false" {
			log.Printf("Disable Sealights On Init was configured during the scan but overridden by the environment variable SL_DISABLE_ON_INIT=false")
		} else {
			log.Printf("Sealights agent is disabled")
			return false
		}
	}
	instrumentMap = cfg.InstrumentMap
	sealightsAgent = internal.NewAgent()
	sealightsAgent.SetOnStopFunc(onStop)
	sealightsAgent.SetCollectorStatsFunc(getCollectorStats)
	err := sealightsAgent.Init(context.Background(), cfg)
	if err != nil {
		log.Printf("Failed to initialized Sealights embedded agent: %s, Sealights agent is disabled", err.Error())
		return false
	}
	if err := sealightsAgent.Start(); err != nil {
		if cfg.DisableOnNoConnection {
			panic(fmt.Sprintf("Failed to start Sealights embedded agent: %s and Disable on No Connection was set to true", err.Error()))
		} else {
			log.Printf("Failed to start Sealights embedded agent: %s, Sealights agent is disabled", err.Error())
		}
		return false
	}
	sealightsCollector = sealightsAgent
	return true
}

func getCollectorStats() string {
	return colStats.String()
}
