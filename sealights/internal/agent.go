package internal

import (
	"context"
	"fmt"
	"github.com/konflux-ci/release-service/__sealights__/models/buildsession_models"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/config"
	"github.com/konflux-ci/release-service/__sealights__/models/agent_models"
	"github.com/konflux-ci/release-service/__sealights__/models/buildmap_models"
	clock_models "github.com/konflux-ci/release-service/__sealights__/models/clock_modles"
	"github.com/konflux-ci/release-service/__sealights__/models/configuration_models"
	"github.com/konflux-ci/release-service/__sealights__/models/execution_models"
	"github.com/konflux-ci/release-service/__sealights__/models/footprint_models"
	"github.com/konflux-ci/release-service/__sealights__/models/logs_models"
	"github.com/konflux-ci/release-service/__sealights__/models/tests_models"
	"github.com/konflux-ci/release-service/__sealights__/models/tia_models"
	clock "github.com/konflux-ci/release-service/__sealights__/pkg/clock"
	"github.com/konflux-ci/release-service/__sealights__/pkg/common"
	"github.com/konflux-ci/release-service/__sealights__/pkg/logger"
	"github.com/konflux-ci/release-service/__sealights__/pkg/target"
	"github.com/konflux-ci/release-service/__sealights__/pkg/target/http"
)

// agentModeType - the runningMode of the Agent is running in
type agentModeType int32

const (
	agentModeTypeUnknown          agentModeType = 0 // unknown runningMode - in start up
	agentModeTypeUnitTests        agentModeType = 1 // running in unit tests runningMode
	agentModeTypeIntegrationTests agentModeType = 2 // running in integration tests runningMode
	agentModeTypeLightMode        agentModeType = 3 // running agent in light mode
	agentModeTypeTestsRunner
)

// Agent - the sealerights Agent
type Agent struct {
	sync.Mutex
	target                    target.Target
	agentConfig               *agent_models.AgentConfig
	agentMetadata             *agent_models.AgentMetadata
	ctx                       context.Context    // context - the context of the Agent
	cancelFunc                context.CancelFunc // cancel function - cancels the context
	activeTestFuncMap         map[string]int64   // active test function map - map of active test functions for test duration calculations
	eventsBuffer              *eventsBuffer      // test events buffer - buffer for test events processing
	footprintBuffer           *footprintsMap     // footprint buffer - buffer for footprint events processing
	runningMode               int32              // runningMode - the runningMode of the Agent
	isAgentInIdleMode         int32
	setAgentModeSync          sync.Once // setAgentModeSync - the sync.once for setting the runningMode only one time
	onStopFunc                func()
	logger                    *logger.Logger
	tiaMap                    map[string]string
	activeTests               *activeTestMap
	footprintsCollectionClock *footprintsCollectionClockInterval
	testSelectionStatus       string
	fnCollectorStats          func() string
}

// NewAgent - creates a new Agent
func NewAgent() *Agent {
	return &Agent{
		activeTestFuncMap:   map[string]int64{},
		runningMode:         int32(agentModeTypeUnknown),
		tiaMap:              map[string]string{},
		testSelectionStatus: "",
		activeTests:         newActiveTestMap(),
		isAgentInIdleMode:   1,
	}
}

func (a *Agent) SetOnStopFunc(onStopFunc func()) {
	a.onStopFunc = onStopFunc
}

func (a *Agent) SetCollectorStatsFunc(fn func() string) {
	a.fnCollectorStats = fn
}

// Init - initializes the Agent and starting events loop processing
func (a *Agent) Init(ctx context.Context, cfg *config.Config) error {
	a.ctx, a.cancelFunc = context.WithCancel(ctx)
	var err error
	a.logger = logger.NewLogger("Sealights-Agent", cfg.LogLevel)
	a.logger.Debugf("Sealights agent configuration: %s", cfg.String())
	a.logger.Debugf("Sealights environment configuration: %s", cfg.GetAllSealightsEnv())
	a.agentConfig, err = a.getAgentConfig(cfg)
	a.agentMetadata = agent_models.NewAgentMetadata()
	a.footprintsCollectionClock = newFootprintsClockInterval(a.agentConfig.GetFootprintsCollectionIntervalDuration())
	a.footprintsCollectionClock.run(ctx)
	a.agentMetadata.AgentInfo.
		SetAgentType(agent_models.AgentTypeTestListener).
		SetAgentVersion(a.agentConfig.GetAgentVersion())
	if err != nil {
		return err
	}
	if cfg.LightMode {
		atomic.StoreInt32(&a.runningMode, int32(agentModeTypeLightMode))
		if a.agentConfig.GetTestStage() == "" {
			a.agentConfig.SetTestStage("IntegrationTests")
		}
	}
	a.eventsBuffer = newEventBuffer(int(cfg.EventBufferSize))
	a.footprintBuffer = newFootprintsMap()

	targetConfig := http.
		NewConfig().
		SetToken(cfg.Token).
		SetServerUrl(cfg.ServerUrl).
		SetProxyUrl(cfg.ProxyUrl).
		SetCollectorUrl(cfg.CollectorUrl).
		SetLogLevel(cfg.LogLevel).
		SetRetryCount(int(cfg.ConnectionRetryCount)).
		SetRetryMaxWaitTime(cfg.ConnectionRetryMaxWaitTime).
		SetRetryWaitTime(cfg.ConnectionRetryWaitTime).
		SetTimeout(cfg.ConnectionTimeout)
	if cfg.LightMode {
		targetConfig.SetHeaders(map[string]string{
			"X-Sealights-Agent-Mode": "light",
		})
		targetConfig.SetUseCollectorIfExists(true)
	}
	targetClient := http.NewHTTPTarget()
	if err := targetClient.Init(a.ctx, targetConfig); err != nil {
		return err
	}
	a.target = targetClient
	if a.agentConfig.GetTestsRunnerMode() {
		if err := a.initAsTestRunner(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Agent) getAgentConfig(cfg *config.Config) (*agent_models.AgentConfig, error) {
	agentConfig := agent_models.NewAgentConfig().
		SetAppName(cfg.AppName).
		SetBuild(cfg.Build).
		SetModuleName(cfg.ModuleName).
		SetBranch(cfg.Branch).
		SetAgentId(cfg.AgentId).
		SetBuildSessionId(cfg.BuildSessionId).
		SetLabId(cfg.LabId).
		SetServerUrl(cfg.ServerUrl).
		SetProxyUrl(cfg.ProxyUrl).
		SetLogLevel(cfg.LogLevel).
		SetTestStage(cfg.TestStage).
		SetAgentVersion(cfg.AgentVersion).
		SetDisableTests(cfg.NoTests).
		SetEnableRemoteConfig(cfg.EnableRemoteConfig).
		SetTestSelection(cfg.TestSelection).
		SetRemoteConfigIntervalSec(cfg.RemoteConfigIntervalSec).
		SetTestGroupId(cfg.TestGroupId).
		SetTestRecommendationSleepSeconds(cfg.TestRecommendationSleepSeconds).
		SetTestsRunnerMode(cfg.TestsRunnerMode).
		SetGinkgoEnabled(cfg.GinkgoEnabled)

	if err := agentConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid agent config: %s", err.Error())
	}

	return agentConfig, nil
}

func (a *Agent) initAsTestRunner() error {
	a.logger.Infof("agent is initializing as test runner, checking build session id or lab id and test stage")
	currentTestStage := os.Getenv("SEALIGHTS_TEST_STAGE")
	if currentTestStage == "" {
		return fmt.Errorf("cannot initialize as test runner, missing SEALIGHTS_TEST_STAGE env var")
	}
	a.agentConfig.SetTestStage(currentTestStage)
	currentBuildSessionId := os.Getenv("SEALIGHTS_BUILD_SESSION_ID")
	if currentBuildSessionId != "" {
		a.agentConfig.SetBuildSessionId(currentBuildSessionId)
		a.logger.Infof("agent is initialized as test runner with build session id: %s and test stage: %s", currentBuildSessionId, currentTestStage)
		return nil
	}
	currentLabId := os.Getenv("SEALIGHTS_LAB_ID")
	if currentLabId == "" {
		return fmt.Errorf("cannot initialize as test runner, missing SEALIGHTS_LAB_ID env var as no build session id was provided")
	}
	a.logger.Infof("agent is getting build session id by lab id: %s", currentLabId)
	a.agentConfig.SetLabId(currentLabId)
	req := target.NewRequest().
		SetType(target.RequestTypeAgentGetBuildSessionIdByLabId).
		SetData(buildsession_models.NewBuildSessionByLabIdRequest().SetLabId(currentLabId))
	resp, err := a.target.Do(a.ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get build session id by lab id: %s", err.Error())
	}
	buildSessionResp := resp.GetData().(*buildsession_models.BuildSessionResponse)
	if buildSessionResp.BuildSessionId == "" {
		return fmt.Errorf("failed to get build session id by lab id, empty build session id")
	}
	a.agentConfig.SetBuildSessionId(buildSessionResp.BuildSessionId)
	a.logger.Infof("agent is initialized as test runner with build session id: %s and test stage: %s", buildSessionResp.BuildSessionId, currentTestStage)
	return nil
}
func (a *Agent) Start() error {
	if atomic.LoadInt32(&a.runningMode) == int32(agentModeTypeLightMode) {
		a.logger.Debugf("agent is running in light mode")
		go a.runProcessFootprintsEvents(a.ctx)
		return nil
	}

	a.checkForceMode()
	// sending Agent start event
	if err := a.sendAgentStart(); err != nil {
		return err
	}

	// try to get clock offset from the server
	_ = a.sendCLockSync()

	if a.agentConfig.GetEnableRemoteConfig() {
		a.logger.Debugf("agent is checking remote config")
		if err := a.processRemoteConfig(a.ctx); err != nil {
			return err
		}
		a.logger.Debugf("agent  checking remote config finished")
	}
	go a.runAgentHeartbeat(a.ctx)
	go a.runProcessTestEvents(a.ctx)
	go a.runProcessFootprintsEvents(a.ctx)
	go a.runGetRemoteConfig(a.ctx)
	go a.runReportLogs(a.ctx)
	go a.runNotifyIdleMode(a.ctx)
	return nil
}

// this function is checking if env SEALIGHTS_FORCE_MODE is set  and apply the force mode
func (a *Agent) checkForceMode() {
	a.logger.Debugf("agent is checking force mode")
	defer a.logger.Debugf("agent finished checking force mode")
	forceMode := os.Getenv("SEALIGHTS_FORCE_MODE")
	if forceMode != "" {
		forceModeInt, err := strconv.Atoi(forceMode)
		if err != nil {
			a.logger.Errorf("failed to parse SEALIGHTS_FORCE_MODE env var: %s", err.Error())
			return
		}
		if forceModeInt < 0 || forceModeInt > 3 {
			a.logger.Errorf("invalid SEALIGHTS_FORCE_MODE env var: %s", forceMode)
			return
		}
		a.setAgentMode(agentModeType(forceModeInt))
		a.logger.Warnf("agent is running in force mode: %s", forceMode)
	} else {
		a.logger.Debugf("SEALIGHTS_FORCE_MODE is not set")
	}

}
func (a *Agent) loadTiaMap() error {
	a.logger.Debugf("agent is loading tia map")
	req := target.NewRequest().
		SetType(target.RequestTypeAgentGetTestRecommendations).
		SetData(tia_models.NewTestRecommendationsRequest().
			SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
			SetStage(a.agentConfig.GetTestStage()).
			SetTestGroupId(a.agentConfig.GetTestGroupId()))
	var excludeTests []*tia_models.TestRecommendationsExcludedTests
	test5SecInterval := int(a.agentConfig.GetTestRecommendationSleepSeconds() / 5)
	for i := 0; i < test5SecInterval; i++ {
		resp, err := a.target.Do(a.ctx, req)
		if err != nil {
			a.logger.Debugf("agent failed to get test selection from server: %s", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		respData := resp.GetData().(*tia_models.TestRecommendationsResponse)
		excludeTests = respData.ExcludedTests
		a.testSelectionStatus = respData.GetTestSelectionStatus(a.agentConfig.GetTestSelection())
		if a.testSelectionStatus != "recommendationsTimeout" {
			a.logger.Debugf("agent got test selection status: %s", a.testSelectionStatus)
			break
		} else {
			a.logger.Errorf("agent got test selection status: %s, retry %d of 12, next retry in 5 seconds", a.testSelectionStatus, i+1)
			time.Sleep(5 * time.Second)
		}
	}
	for _, test := range excludeTests {
		a.tiaMap[test.Name] = test.Name
	}
	return nil
}

func (a *Agent) CheckSkipTest(testName string) bool {
	if _, ok := a.tiaMap[testName]; ok {
		return true
	}
	return false
}

// RecordTestEvent - records a test event
func (a *Agent) RecordTestEvent(event Event) {
	a.eventsBuffer.addPre(event)
}

// RecordFootprint - records a footprint event
func (a *Agent) RecordFootprint(event Event) {
	a.footprintBuffer.add(event)
}

// AddActiveTest - adds an active test
func (a *Agent) AddActiveTest(testName string) {
	a.activeTests.add(testName)
}

// RemoveActiveTest - removes an active test
func (a *Agent) RemoveActiveTest(testName string) {
	a.activeTests.remove(testName)
}

// GetActiveTests - returns the active tests
func (a *Agent) GetActiveTests() []string {
	return a.activeTests.list()
}

func (a *Agent) GetFootprintsTimestamps() (int64, int64) {
	return a.footprintsCollectionClock.getsTimestamps()
}

// RecordControlEvent - records a control event
func (a *Agent) RecordControlEvent(event Event) {
	switch event.eventType {
	case EventTypeStartMain: // start main events
		a.setAgentModeSync.Do(func() {
			a.setAgentMode(agentModeTypeIntegrationTests)
		}) // calling to setAgentModeSync.Do() will make sure that setAgentMode() is called only once
		a.footprintBuffer.add(event)
		a.logger.Debug("agent received start main event")
	case EventTypeEndMain: // end main event - no action required
		a.logger.Debug("agent received end main event")
	case EventTypeStartTestMain: // start test main events
		a.setAgentModeSync.Do(func() {
			a.setAgentMode(agentModeTypeUnitTests)
		}) // calling to setAgentModeSync.Do() will make sure that setAgentMode() is called only once
		a.logger.Debug("agent received start test main event")
	case EventTypeEndTestMain: // end test main events
		a.logger.Debug("agent received end test main event")
	default:
		a.logger.Warnf("Unknown event type: %d", event.eventType)
	}
}

// Stop - stops the Agent and flush all pending events
func (a *Agent) Stop() {
	if atomic.LoadInt32(&a.runningMode) == int32(agentModeTypeLightMode) {
		a.logger.Debugf("stopping agent in light mode and flushing all pending events")
		a.processFootprintEvents()
		return
	}
	a.logger.Debugf("Stopping Agent")
	a.cancelFunc()
	a.logger.Debugf("Flushing buffered test events")
	a.processTestEvents()
	a.logger.Debugf("Flushing buffered footprint events")
	a.processFootprintEvents()
	a.logger.Debugf("Report execution end")
	mode := atomic.LoadInt32(&a.runningMode)
	if a.agentConfig.GetExecutionId() != "" && mode == int32(agentModeTypeUnitTests) {
		if err := a.sendEndExecutionEvent(); err != nil {
			a.logger.Warnf("Failed to send enc execution event: %s", err.Error())
		}
	}
	a.logger.Debugf("report agent stopped")
	if err := a.sendAgentStop(); err != nil {
		a.logger.Warnf("Failed to send agent atop event: %s", err.Error())
	}
}

// processTestEvents - process all buffered test events
func (a *Agent) processTestEvents() {
	if !a.readyToReport() { // are we ready to report? aka,we have execution id
		return
	}

	events := a.eventsBuffer.get() // get all buffered events
	if len(events) == 0 {          // no events to report - return
		return
	}
	if a.agentConfig.GetDisableTests() {
		a.logger.Debugf("report tests are disabled, ignoring test %d events", len(events))
		return
	}
	a.logger.Debugf("processing %d buffered test events", len(events))
	testEvents := tests_models.NewTestsEvents(). // create test events
							SetBuild(a.agentConfig.GetBuild()).
							SetBranch(a.agentConfig.GetBranch()).
							SetAppName(a.agentConfig.GetAppName()).
							SetAgentId(a.agentConfig.GetAgentId()).
							SetTestSelectionStatus(a.testSelectionStatus).
							SetLabId(a.agentConfig.GetLabId(), a.agentConfig.GetBuildSessionId()).
							SetTestStage(a.agentConfig.GetTestStage())
	for _, event := range events { // add all buffered events to test events
		switch event.eventType {
		case EventTypeStartTest: // process start test and start main event
			_, found := a.activeTestFuncMap[event.function] // check if the test function is already active and reported
			if found {
				continue
			}
			a.Lock()
			a.activeTestFuncMap[event.function] = event.timestamp // add the test function to the active test function map
			a.Unlock()
			testEvents.AddEvent(
				tests_models.NewTestEvent(tests_models.TestEventTypeTestStart).
					SetExecutionId(a.agentConfig.GetExecutionId()).
					SetTestName(event.function).
					SetTimestamp(event.timestamp))
		case EventTypeStartTestMain:

		case EventTypeEndTest: // process end test event
			a.Lock()                                                // lock the active test function map
			startTime, found := a.activeTestFuncMap[event.function] // get the start time of the test function
			if !found {                                             // test function is not active
				a.Unlock()
				continue
			}
			delete(a.activeTestFuncMap, event.function) // remove the test function from the active test function map
			a.Unlock()
			endTestEvent := tests_models.NewTestEvent(tests_models.TestEventTypeTestEnd).
				SetExecutionId(a.agentConfig.GetExecutionId()).
				SetTestName(event.function).
				SetTimestamp(event.timestamp)
			if duration := event.timestamp - startTime; duration > 0 { // calculate the test duration
				endTestEvent.SetDuration(duration)
			} else {
				endTestEvent.SetDuration(1)
			}
			skippedStr := event.attributes["skipped"]
			result := event.attributes["result"]
			if skippedStr == "true" {
				endTestEvent.SetResult(tests_models.TestEventResultSkipped)
			} else {
				switch result {
				case "true": // test passed
					endTestEvent.SetResult(tests_models.TestEventResultPassed)
				case "false": // test failed
					endTestEvent.SetResult(tests_models.TestEventResultFailed)
				default: // test result is unknown
					endTestEvent.SetResult(tests_models.TestEventResultSkipped)
				}
			}
			testEvents.AddEvent(endTestEvent) // add the test end event to the test events
		case EventTypeEndTestMain: // process end test main event

		default:
			a.logger.Debugf("Unknown test event type %d", event.eventType)
		}
	}
	if len(testEvents.Events) > 0 { // if there are test events to report
		a.logger.Debugf("Reporting %d test events", len(testEvents.Events))
		req := target.NewRequest().
			SetType(target.RequestTypeAgentSendTestEvents).
			SetData(testEvents)
		if _, err := a.target.Do(context.Background(), req); err != nil { // send the test events
			a.logger.Warnf("Failed to send test event: %s", err.Error())
			a.eventsBuffer.addBuffer(events...)                                       // add the buffered events back to the buffer
			if err := a.sendAgentTestEventsSubmissionError(err.Error()); err != nil { // send the Agent test events submission error
				a.logger.Warnf("Failed to send agent test events submission error: %s", err.Error())
			}
		} else {
			a.logger.Debugf("Successfully reported %d test events", len(testEvents.Events))
		}
	}
}

// processFootprintEvents - process all buffered footprint events
func (a *Agent) processFootprintEvents() {
	if !a.readyToReport() { // are we ready to report? aka,we have execution id
		return
	}
	events := a.footprintBuffer.get() // get all pending events from buffer
	if len(events) == 0 {             // no events to report, return
		return
	}
	//if a.agentConfig.GetTestsRunnerMode() {
	//	a.logger.Debugf("agent is running in test runner mode, skipping footprint events")
	//	return
	//}
	if a.agentConfig.GetDisableFootprints() {
		a.logger.Debugf("Footprints are disabled, ignoring %d footprint events", len(events))
		return
	}
	a.logger.Debugf("processing %d footprint events", len(events))
	footprint := footprint_models.NewFootprint().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetExecutionId(a.agentConfig.GetExecutionId()).
		SetLabId(a.agentConfig.GetLabId()).
		SetTimedFootprintsCollectionIntervalSeconds(int64(a.agentConfig.GetFootprintsReportingIntervalDuration().Seconds()))
	if atomic.LoadInt32(&a.runningMode) == int32(agentModeTypeLightMode) {
		footprint.SetAllowEmptyExecutionId(true).
			SetAgentConfig(a.agentConfig).
			SetAgentMetadata(a.agentMetadata)
	}
	for _, event := range events { // iterate over all events and add them to footprint model
		footprint.AddFootprint(event.activeTestList, event.function, event.footprintsStartTime, event.footprintsEndTime, event.timestamp < a.agentConfig.GetExecutionOpenTime())
	}
	switch atomic.LoadInt32(&a.runningMode) {
	case int32(agentModeTypeLightMode), int32(agentModeTypeIntegrationTests):
		footprint.SetAgentConfig(a.agentConfig).
			SetAgentMetadata(a.agentMetadata).
			Complete(footprint_models.FootprintAggregationModeAllHits)
	case int32(agentModeTypeUnitTests):
		footprint.SetAgentConfig(a.agentConfig).
			SetAgentMetadata(a.agentMetadata).
			Complete(footprint_models.FootprintAggregationModeCoverage)
	}
	methods, hits, span := footprint.Stats()
	a.logger.Debugf("Reporting footprint with %d methods,%d hits in %d seconds interval", methods, hits, span)

	if len(footprint.Methods) > 0 { // we have events to report
		req := target.NewRequest().
			SetType(target.RequestTypeAgentSendFootprints).
			SetData(footprint).
			SetMetadataKeyValue("testStage", a.agentConfig.GetTestStage()).
			SetMetadataKeyValue("buildSessionId", a.agentConfig.GetBuildSessionId()).
			SetMetadataKeyValue("executionBuildSessionId", a.agentConfig.GetExecutionBuildSessionId())
		a.logger.Debugf("Reporting footprint event with %d events", len(footprint.Methods))
		if _, err := a.target.Do(context.Background(), req); err != nil { // send footprint to server
			a.logger.Warnf("Failed to send footprint: %s", err.Error())
			a.footprintBuffer.add(events...)                                         // add events back to buffer for trying again later
			if err := a.sendAgentFootprintSubmissionError(err.Error()); err != nil { // send Agent footprint submission error to server
				a.logger.Warnf("Failed to send agent footprint events submission error: %s", err.Error())
			}
		} else {
			a.logger.Debugf("Successfully reported footprint event with %d events", len(footprint.Methods))
		}
	}
}

func (a *Agent) runAgentHeartbeat(ctx context.Context) {
	for {
		select {
		case <-time.After(a.agentConfig.GetAgentHeartbeatIntervalDuration()):
			a.sendAgentHeartbeat() // process heartbeat in separate goroutine
			_ = a.sendCLockSync()
			a.printDebugInfo() // update clock sync
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) runProcessTestEvents(ctx context.Context) {
	for {
		select {
		case <-time.After(a.agentConfig.GetTestsIntervalDuration()):
			a.processTestEvents() // process heartbeat in separate goroutine
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) runProcessFootprintsEvents(ctx context.Context) {
	for {
		select {
		case <-time.After(a.agentConfig.GetFootprintsReportingIntervalDuration()):
			a.processFootprintEvents() // process heartbeat in separate goroutine
		case <-ctx.Done():
			return
		}
	}
}

// sendAgentStart - sends Agent start event to server
func (a *Agent) sendAgentStart() error {
	a.logger.Debugf("Sending Agent start event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			agent_models.NewAgentStartEvent().
				SetLabId(a.agentConfig.GetLabId()).
				SetAgentVersion(a.agentConfig.GetAgentVersion()).
				SetTestStage(a.agentConfig.GetTestStage()).
				SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
				SetAgentType(agent_models.AgentTypeTestListener))

	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent start event: %s", err.Error())
		return err
	}
	a.logger.Debugf("Sending Agent start event completed")
	return nil
}

// sendAgentStop - sends Agent Stop event to server
func (a *Agent) sendAgentStop() error {
	a.logger.Debugf("send agent Stop event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			agent_models.NewAgentStopEvent())
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent Stop event: %s", err.Error())
		return err
	}
	a.logger.Debugf("send agent Stop event completed")
	return nil
}

func (a *Agent) sendCLockSync() error {
	a.logger.Debugf("sending clock sync request")
	clockSyncRequest := clock_models.NewClockSyncRequest()
	req := target.NewRequest().
		SetType(target.RequestTypeClockSync).
		SetData(clockSyncRequest)
	var res target.TargetResponse
	var err error
	if res, err = a.target.Do(context.Background(), req); err != nil {
		a.logger.Debugf("Failed to get clock sync: %s", err.Error())
		return err
	}
	clockSyncResponse, ok := res.GetData().(*clock_models.ClockSyncResponse)
	if !ok {
		a.logger.Warnf("Failed to parse clock sync response")
		return fmt.Errorf("failed to parse clock sync response")
	}
	clock.GlobalClockSync.SetCurrentOffset(clockSyncResponse.Offset)
	a.logger.Debugf("global clock update with offset: %d", clockSyncResponse.Offset)
	a.logger.Debugf("Sending clock sync request completed")
	return nil
}

// sendAgentFootprintSubmissionError - sends Agent footprint submission error event to server
func (a *Agent) sendAgentFootprintSubmissionError(msg string) error {
	a.logger.Debugf("send agent footprint submission error event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			footprint_models.NewFootprintErrorEvent().SetData(msg))
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent footprint submission error event: %s", err.Error())
		return err
	}
	a.logger.Debugf("send agent footprint submission error event completed")
	return nil
}

// sendAgentTestEventsSubmissionError - sends Agent test events submission error event to server
func (a *Agent) sendAgentTestEventsSubmissionError(msg string) error {
	a.logger.Debugf("send agent test events submission error event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			tests_models.NewTestErrorEventEvent().SetData(msg))
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent test events submission error event: %s", err.Error())
		return err
	}
	a.logger.Debugf("send agent test events submission error event completed")
	return nil
}

// sendAgentBuildMapSubmissionError - sends Agent build map submission error event to server
func (a *Agent) sendAgentBuildMapSubmissionError(msg string) error {
	a.logger.Debugf("Send agent build map events submission error event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			buildmap_models.NewBuildMapErrorEvent().SetData(msg))
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent build map events submission error event: %s", err.Error())
		return err
	}
	a.logger.Debugf("send agent build map events submission error event completed")
	return nil
}

// sendStartExecutionEvent - sends start execution event to the server
func (a *Agent) sendStartExecutionEvent(executionId string) error {
	a.logger.Debugf("Sending execution started event")
	executionEvent := execution_models.NewExecutionEvent().
		SetExecutionId(executionId).
		SetAgentId(a.agentConfig.GetAgentId()).
		SetTestStage(a.agentConfig.GetTestStage()).
		SetTestGroupId(a.agentConfig.GetTestGroupId()).
		SetBuildName(a.agentConfig.GetBuild()).
		SetBranchName(a.agentConfig.GetBranch()).
		SetAppName(a.agentConfig.GetAppName()).
		SetLabId(a.agentConfig.GetLabId())

	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendExecutionStart).
		SetData(executionEvent)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send test execution started: %s", err.Error())
		return err
	}
	a.logger.Debugf("Sending execution started event completed")
	return nil
}

// sendEndExecutionEvent - sends execution ended event to the server
func (a *Agent) sendEndExecutionEvent() error {
	a.logger.Debugf("Sending execution ended event")
	query := fmt.Sprintf("labId=%s&executionId=%s", a.agentConfig.GetLabId(), a.agentConfig.GetExecutionId())
	if a.agentConfig.GetTestGroupId() != "" {
		query = fmt.Sprintf("%s&testGroupId=%s", query, a.agentConfig.GetTestGroupId())
	}
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendExecutionEnd).
		SetMetadataKeyValue("query", query)

	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send test execution ended event: %s", err.Error())
		return err
	}
	a.logger.Debugf("Sending execution ended event completed")
	return nil
}

func (a *Agent) addMetadataInfoToRequest(r *target.Request, messageType string) {
	r.SetMetadataKeyValue("x-sl-appname", a.agentConfig.GetAppName())
	r.SetMetadataKeyValue("x-sl-branchname", a.agentConfig.GetBranch())
	r.SetMetadataKeyValue("x-sl-buildname", a.agentConfig.GetBuild())
	r.SetMetadataKeyValue("x-sl-bsid", a.agentConfig.GetBuildSessionId())
	r.SetMetadataKeyValue("x-sl-agentid", a.agentConfig.GetAgentId())
	if labid := a.agentConfig.GetLabId(); labid != "" {
		r.SetMetadataKeyValue("x-sl-labid", labid)
	}
	if executionId := a.agentConfig.GetExecutionId(); executionId != "" {
		r.SetMetadataKeyValue("x-sl-executionid", executionId)
	}
	r.SetMetadataKeyValue("x-sl-messagetype", messageType)
}

// sendAgentHeartbeat - send heartbeat Agent event to the server
func (a *Agent) sendAgentHeartbeat() {
	a.logger.Debugf("Sending Agent heartbeat event")
	agentEvent := agent_models.NewAgentEvents().
		SetAgentId(a.agentConfig.GetAgentId()).
		SetBuildSessionId(a.agentConfig.GetBuildSessionId()).
		SetAppName(a.agentConfig.GetAppName()).
		AddEvent(
			agent_models.NewAgentPingEvent().SetLabId(a.agentConfig.GetLabId()))
	req := target.NewRequest().
		SetType(target.RequestTypeAgentSendEvents).
		SetData(agentEvent)
	a.addMetadataInfoToRequest(req, "1003")
	if _, err := a.target.Do(context.Background(), req); err != nil {
		a.logger.Warnf("Failed to send agent heartbeat: %s", err.Error())
	} else {
		a.logger.Debugf("Sending agent heartbeat event completed")
	}
}

// readyToReport returns true if the Agent is ready to report (the qaent mode was set)
func (a *Agent) readyToReport() bool {
	if atomic.LoadInt32(&a.runningMode) > 0 {
		return true
	}
	return false
}

// setRunningMode - sets the running mode of the Agent
func (a *Agent) setAgentMode(mode agentModeType) {
	if atomic.LoadInt32(&a.runningMode) != 0 {
		return
	}
	switch mode {
	case agentModeTypeUnitTests: // we received unit tests events from the client app
		a.logger.Infof("Agent running mode set to unit tests")
		a.startUnitTestsExecution()
	case agentModeTypeIntegrationTests: // we received integration tests events (using main function) from the client app
		a.logger.Infof("Agent running mode set to integration tests")
		go a.startIntegrationTestsExecution()
	}
}

// startUnitTestsExecution - set Agent uin Unit Test mode and sens execution start event
func (a *Agent) startUnitTestsExecution() {
	atomic.StoreInt32(&a.isAgentInIdleMode, 0)
	a.logger.Debugf("Starting unit tests execution")

	if a.agentConfig.GetTestStage() == "" {
		a.agentConfig.SetTestStage("Unit Tests")
	}
	executionId := common.NewUUID()
	for {
		if a.agentConfig.GetDisableAgent() {
			return
		}
		if err := a.sendStartExecutionEvent(executionId); err != nil {
			time.Sleep(time.Second * 5)
			a.agentConfig.SetExecutionOpenTime(0)
			continue
		}
		a.agentConfig.SetExecutionOpenTime(clock.UnixMilli())
		a.agentConfig.SetExecutionId(executionId)
		break
	}
	atomic.StoreInt32(&a.runningMode, int32(agentModeTypeUnitTests))
	a.logger.Debugf("Unit tests execution started with execution id: %s", a.agentConfig.GetExecutionId())
	if err := a.loadTiaMap(); err != nil {
		a.logger.Warnf("Failed to load TIA map: %s", err.Error())
	}
}

// startIntegrationTestsExecution - set Agent in Integration Test mode and sens execution start event
func (a *Agent) startIntegrationTestsExecution() {
	atomic.StoreInt32(&a.isAgentInIdleMode, 0)
	a.logger.Debugf("Starting integration tests execution")
	if a.agentConfig.GetTestStage() == "" {
		a.agentConfig.SetTestStage("IntegrationTests")
	}
	for {
		if a.agentConfig.GetDisableAgent() {
			return
		}
		time.Sleep(a.agentConfig.GetCheckExecutionIntervalDuration())
		getExecution := execution_models.NewExecutionRequest().
			SetLabId(a.agentConfig.GetLabId())
		req := target.NewRequest().
			SetType(target.RequestTypeGetExecutionId).
			SetData(getExecution)
		resp, err := a.target.Do(context.Background(), req)
		if err != nil {
			atomic.StoreInt32(&a.runningMode, int32(agentModeTypeUnknown))
			a.logger.Debugf("Failed to get execution id: %s", err.Error())
			continue
		}
		execution, ok := resp.GetData().(*execution_models.ExecutionResponse)
		if !ok {
			atomic.StoreInt32(&a.runningMode, int32(agentModeTypeUnknown))
			a.logger.Debugf("Failed to get execution object from response")
			continue
		}
		executionId, state := execution.GetExecutionId()
		switch state {
		case execution_models.ExecutionStateUnknown:
			atomic.StoreInt32(&a.runningMode, int32(agentModeTypeUnknown))
			a.agentConfig.
				SetExecutionOpenTime(0).
				SetExecutionId(executionId).
				SetExecutionBuildSessionId(a.agentConfig.GetBuildSessionId()).
				SetTestStage("IntegrationTests").
				SetTestGroupId("")
			a.logger.Debugf("Failed to get execution id, unknown state")
		case execution_models.ExecutionStateRunning:
			a.agentConfig.
				SetExecutionOpenTime(clock.UnixMilli()).
				SetExecutionId(executionId).
				SetExecutionBuildSessionId(execution.Execution.BuildSessionId).
				SetTestStage(execution.Execution.TestStage).
				SetTestGroupId(execution.Execution.TestGroupId)
			atomic.StoreInt32(&a.runningMode, int32(agentModeTypeIntegrationTests))
			a.logger.Debugf("Integration tests execution started with execution id: %s", a.agentConfig.GetExecutionId())
		case execution_models.ExecutionStateFinished:
			currentMode := atomic.LoadInt32(&a.runningMode)
			if currentMode == int32(agentModeTypeIntegrationTests) {
				a.logger.Debugf("Integration tests execution finished flushing footprints")
				a.processFootprintEvents()
				atomic.StoreInt32(&a.runningMode, int32(agentModeTypeUnknown))
				a.agentConfig.
					SetExecutionOpenTime(0).
					SetExecutionId("").
					SetExecutionBuildSessionId(a.agentConfig.GetBuildSessionId()).
					SetTestStage("IntegrationTests").
					SetTestGroupId("")

			}
		}
	}
}

func (a *Agent) runGetRemoteConfig(ctx context.Context) {
	if !a.agentConfig.GetEnableRemoteConfig() {
		return
	}
	a.logger.Debugf("Starting get remote config every %s seconds", a.agentConfig.GetRemoteConfigIntervalDuration())
	for {
		select {
		case <-time.After(a.agentConfig.GetRemoteConfigIntervalDuration()):
			if err := a.processRemoteConfig(ctx); err != nil {
				a.logger.Debugf("Failed to get remote config: %s", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) processRemoteConfig(ctx context.Context) error {
	req := target.NewRequest().
		SetType(target.RequestTypeAgentGetConfiguration).
		SetData(configuration_models.NewConfigurationRequest().
			SetAppName(a.agentConfig.GetAppName()).
			SetAgentId(a.agentConfig.GetAgentId()).
			SetAgentVersion(a.agentConfig.GetAgentVersion()).
			SetTestStage(a.agentConfig.GetTestStage()).
			SetLabId(a.agentConfig.GetLabId()).
			SetAgentType("TestListener").
			SetBranchName(a.agentConfig.GetBranch()).
			SetBuildName(a.agentConfig.GetBuild()))
	resp, err := a.target.Do(ctx, req)
	if err != nil {
		return err
	}
	remoteConfig, ok := resp.GetData().(*configuration_models.ConfigurationResponse)
	if !ok {
		return fmt.Errorf("remote config response is not valid")
	}

	if remoteConfig.DisableAgent {
		a.agentConfig.SetDisableAgent(true)
		a.logger.Warnf("remote config disable agent received")
		a.onStopFunc()
		a.Stop()
		return nil
	}

	if remoteConfig.AgentHeartbeatIntervalSec > 0 &&
		remoteConfig.AgentHeartbeatIntervalSec != int(a.agentConfig.GetAgentHeartbeatIntervalDuration().Seconds()) {
		a.agentConfig.SetAgentHeartbeatIntervalSec(remoteConfig.AgentHeartbeatIntervalSec)
		a.logger.Debugf("remote config update agent heartbeat interval to %d seconds", remoteConfig.AgentHeartbeatIntervalSec)

	}

	if remoteConfig.EnableLogs != a.agentConfig.GetEnableLogs() {
		a.agentConfig.SetEnableLogs(remoteConfig.EnableLogs)
		a.logger.SetEnableBufferLogs(remoteConfig.EnableLogs)
		a.logger.Debugf("remote config update enable logs to %t", remoteConfig.EnableLogs)
	}

	if remoteConfig.LogsIntervalSec > 0 &&
		remoteConfig.LogsIntervalSec != int(a.agentConfig.GetLogsIntervalDuration().Seconds()) {
		a.agentConfig.SetLogsIntervalSec(remoteConfig.LogsIntervalSec)
		a.logger.Debugf("remote config update logs interval to %d seconds", remoteConfig.LogsIntervalSec)

	}

	if remoteConfig.LogLevel != "" &&
		remoteConfig.LogLevel != a.agentConfig.GetLogLevel() {
		a.agentConfig.SetLogLevel(remoteConfig.LogLevel)
		a.logger.Debugf("remote config update log level to %s", remoteConfig.LogLevel)
		a.logger = logger.NewLogger("Sealights-Agent", a.agentConfig.GetLogLevel())
	}
	if remoteConfig.DisableFootprints != a.agentConfig.GetDisableFootprints() {
		a.agentConfig.SetDisableFootprints(remoteConfig.DisableFootprints)
		a.logger.Debugf("remote config update disable footprints to %t", remoteConfig.DisableFootprints)
	}
	if remoteConfig.FootprintsIntervalSec > 0 &&
		remoteConfig.FootprintsIntervalSec != int(a.agentConfig.GetFootprintsReportingIntervalDuration().Seconds()) {
		a.agentConfig.SetFootprintsReportingIntervalSec(remoteConfig.FootprintsIntervalSec)
		a.logger.Debugf("remote config update footprints interval to %d seconds", remoteConfig.FootprintsIntervalSec)
	}

	if remoteConfig.FootprintsCollectIntervalSecs > 0 &&
		remoteConfig.FootprintsCollectIntervalSecs != int(a.agentConfig.GetFootprintsCollectionIntervalDuration().Seconds()) {
		a.agentConfig.SetFootprintsCollectionInterval(remoteConfig.FootprintsCollectIntervalSecs)
		a.footprintsCollectionClock.setTimeInterval(time.Duration(remoteConfig.FootprintsCollectIntervalSecs) * time.Second)
		a.logger.Debugf("remote config update footprints collection interval to %d seconds", remoteConfig.FootprintsCollectIntervalSecs)
	}
	if remoteConfig.DisableTests != a.agentConfig.GetDisableTests() {
		a.agentConfig.SetDisableTests(remoteConfig.DisableTests)
		a.logger.Debugf("remote config update disable tests to %t", remoteConfig.DisableTests)
	}
	if remoteConfig.TestsIntervalSec > 0 &&
		remoteConfig.TestsIntervalSec != int(a.agentConfig.GetTestsIntervalDuration().Seconds()) {
		a.agentConfig.SetTestsIntervalSec(remoteConfig.TestsIntervalSec)
		a.logger.Debugf("remote config update tests interval to %d seconds", remoteConfig.TestsIntervalSec)
	}
	if remoteConfig.CheckExecutionIntervalSec > 0 &&
		remoteConfig.CheckExecutionIntervalSec != int(a.agentConfig.GetCheckExecutionIntervalDuration().Seconds()) {
		a.agentConfig.SetCheckExecutionIntervalSec(remoteConfig.CheckExecutionIntervalSec)
		a.logger.Debugf("remote config update check execution interval to %d seconds", remoteConfig.CheckExecutionIntervalSec)
	}
	return nil
}

func (a *Agent) runReportLogs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(a.agentConfig.GetLogsIntervalDuration()):
			err := a.sendAgentLogs(ctx)
			if err != nil {
				a.logger.Errorf("error sending logs: %s", err.Error())
			}
		}
	}
}

func (a *Agent) sendAgentLogs(ctx context.Context) error {
	logs := a.logger.GetLogs()
	if len(logs) == 0 {
		return nil
	}
	if !a.agentConfig.GetEnableLogs() {
		return nil
	}
	host, _ := os.Hostname()
	req := target.NewRequest().
		SetType(target.RequestTypeAgentLogSubmission).
		SetData(logs_models.NewLogsSubmission().
			SetAppName(a.agentConfig.GetAppName()).
			SetBranchName(a.agentConfig.GetBranch()).
			SetBuildName(a.agentConfig.GetBuild()).
			SetMachineName(host).
			SetLog(logs))
	_, err := a.target.Do(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) printDebugInfo() {
	a.logger.Debugf("collector stats: %s", a.fnCollectorStats())
	a.logger.Debugf("agent internal stats: footprints buffer size: %d", a.footprintBuffer.len())
}

func (a *Agent) runNotifyIdleMode(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 30):
			if atomic.LoadInt32(&a.isAgentInIdleMode) == 1 {
				a.logger.Warn("The SeaLights agent is in idle mode. Neither Test functions or main function code were instrumented. To fix this, you can try instrumenting test and/or main functions code or setting the environment variable SEALIGHTS_FORCE_MODE=1 (Unit Test) or SEALIGHTS_FORCE_MODE=2 (Integration Test))")
			} else {
				return
			}
		}
	}
}
