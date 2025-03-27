package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/konflux-ci/release-service/__sealights__/models/buildsession_models"
	clock_models "github.com/konflux-ci/release-service/__sealights__/models/clock_modles"
	"github.com/konflux-ci/release-service/__sealights__/models/configuration_models"
	"github.com/konflux-ci/release-service/__sealights__/models/execution_models"
	"github.com/konflux-ci/release-service/__sealights__/models/tia_models"
	"github.com/konflux-ci/release-service/__sealights__/pkg/logger"
	"github.com/konflux-ci/release-service/__sealights__/pkg/target"
	"github.com/konflux-ci/release-service/__sealights__/pkg/target/http/resty"
)

type Target struct {
	ctx       context.Context
	canceler  context.CancelFunc
	cfg       *Config
	client    *resty.Client
	logger    *logger.Logger
	traceFile *os.File
}

func NewHTTPTarget() *Target {
	return &Target{}
}

func (t *Target) Init(ctx context.Context, cfg *Config) error {
	t.ctx, t.canceler = context.WithCancel(ctx)
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	t.cfg = cfg
	t.logger = logger.NewLogger("http-target", cfg.LogLevel)
	t.client = resty.New()

	t.client.SetAuthToken(cfg.Token)
	if cfg.ProxyUrl != "" {
		t.client.SetProxy(cfg.ProxyUrl)

	}
	if strings.HasPrefix(cfg.ServerUrl, "https") {
		t.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	t.client.SetRetryCount(cfg.RetryCount)
	t.client.SetRetryWaitTime(cfg.RetryWaitTime)
	t.client.SetRetryMaxWaitTime(cfg.RetryMaxWaitTime)
	t.client.SetTimeout(cfg.Timeout)

	if t.cfg.isDebugLevel() {
		t.client.SetLogger(t.logger)
		if err := t.initDebugging(); err != nil {
			return nil
		}
	}
	return nil
}

func (t *Target) initDebugging() error {
	t.client.SetDebug(true)
	var err error
	t.traceFile, err = os.OpenFile(fmt.Sprintf(".sealights-debug-log.json"), os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("error opening agent log trace file: %w", err)
	}
	t.client.OnRequestLog(func(requestLog *resty.RequestLog) error {
		if requestLog == nil {
			return nil
		}
		_, _ = t.traceFile.WriteString(NewDebugLogForRequest(requestLog).String())
		return nil
	})

	t.client.OnResponseLog(func(responseLog *resty.ResponseLog) error {
		if responseLog == nil {
			return nil
		}
		_, _ = t.traceFile.WriteString(NewDebugLogForResponse(responseLog).String())
		return nil
	})

	t.client.OnError(func(request *resty.Request, err error) {
		if err == nil {
			return
		}
		_, _ = t.traceFile.WriteString(NewDebugLogForError(request, err).String())
	})

	return nil
}

func (t *Target) Do(ctx context.Context, req target.TargetRequest) (resp target.TargetResponse, err error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	switch req.GetType() {
	case target.RequestTypeUndefined:
		t.logger.Errorf("request type is undefined")
		return nil, fmt.Errorf("undefined request type")
	case target.RequestTypeAgentSendEvents:
		return t.doPostAgentEvent(ctx, req)
	case target.RequestTypeAgentSendFootprints:
		return t.doPostFootprints(ctx, req)
	case target.RequestTypeGetExecutionId:
		return t.doGetExecutionId(ctx, req)
	case target.RequestTypeAgentSendBuildMap:
		return t.doPostSendBuildMap(ctx, req)
	case target.RequestTypeAgentSendBuildMapEnd:
		return t.doPostSendBuildMapEnd(ctx, req)
	case target.RequestTypeAgentSendExecutionStart:
		return t.doPostSendExecutionStart(ctx, req)
	case target.RequestTypeAgentSendExecutionEnd:
		return t.doDeleteSendExecutionEnd(ctx, req)
	case target.RequestTypeAgentGetBuildSessionId:
		return t.doPostGetBuildSessionId(ctx, req)
	case target.RequestTypeAgentSendTestEvents:
		return t.doPostSendTestEvent(ctx, req)
	case target.RequestTypeAgentGetConfiguration:
		return t.doGetAgentConfiguration(ctx, req)
	case target.RequestTypeAgentLogSubmission:
		return t.doPostAgentLogSubmission(ctx, req)
	case target.RequestTypeClockSync:
		return t.doGetClockSync(ctx, req)
	case target.RequestTypeAgentGetTestRecommendations:
		return t.doGetTestRecommendations(ctx, req)
	case target.RequestTypeAgentGetBuildSessionIdByLabId:
		return t.doGetBuildSessionIDByLabID(ctx, req)
	default:
		return nil, fmt.Errorf("request type not supported")
	}
}

func (t *Target) getHttpResponseError(resp *resty.Response) error {
	txid := resp.RawResponse.Header.Get("X-Sl-Txid")
	message := ""
	if resp.Body() != nil {
		message = string(resp.Body())
	}
	return fmt.Errorf("code: %d, message: %s, txid: %s", resp.StatusCode(), message, txid)
}

func (t *Target) doPostAgentEvent(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {

	httpRequest := t.createRequest(ctx).
		SetContext(ctx).
		SetBody(req.GetData())
	for k, v := range req.GetMetadata() {
		httpRequest.SetHeader(k, v)
	}
	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v3/agents/agent-events"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doGetExecutionId(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	getExecutionIdRequest, ok := req.GetData().(*execution_models.ExecutionRequest)
	if !ok {
		return nil, fmt.Errorf("doGetExecutionId: invalid request type")
	}
	getExecutionIdResponse := execution_models.NewExecutionResponse()
	httpRequest := t.createRequest(ctx).
		SetResult(getExecutionIdResponse)
	path := fmt.Sprintf("%s%s", t.cfg.getBaseUrl("/v4/testExecution/"), getExecutionIdRequest.LabId)
	httpResp, err := httpRequest.Get(path)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	if err := getExecutionIdResponse.Validate(); err != nil {
		return nil, fmt.Errorf("doGetExecutionId: invalid response: %s", err.Error())
	}
	return target.NewResponse().
		SetType(target.ResponseTypeGetExecution).
		SetData(getExecutionIdResponse), nil
}

func (t *Target) doPostFootprints(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	testStage := req.GetMetadata().Get("testStage")
	if testStage == "" {
		return nil, fmt.Errorf("testStage is empty")
	}
	buildSessionId := req.GetMetadata().Get("buildSessionId")
	if buildSessionId == "" {
		return nil, fmt.Errorf("buildSessionId is empty")
	}
	executionBuildSessionId := req.GetMetadata().Get("executionBuildSessionId")
	if executionBuildSessionId == "" {
		executionBuildSessionId = buildSessionId
	}
	path := fmt.Sprintf("%s/v6/agents/%s/footprints/%s/%s", t.cfg.getBaseUrl(""), executionBuildSessionId, testStage, buildSessionId)
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())
	httpResp, err := httpRequest.Post(path)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doPostSendBuildMap(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())
	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v5/agents/buildmapping"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doPostSendBuildMapEnd(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())
	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v3/agents/buildend"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doPostSendTestEvent(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())
	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v2/agents/events"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doPostSendExecutionStart(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())
	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v3/testExecution"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doDeleteSendExecutionEnd(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetQueryString(req.GetMetadata().
			Get("query"))
	httpResp, err := httpRequest.Delete(t.cfg.getBaseUrl("/v3/testExecution"))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	return nil, nil
}

func (t *Target) doPostGetBuildSessionId(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	resp := buildsession_models.NewGetBuildSessionIdResponse()
	httpRequest := t.createRequest(ctx).
		SetResult(&resp.BuildSessionId).
		SetBody(req.GetData())

	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v2/agents/buildsession"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}

	targetResp := target.NewResponse().SetType(target.ResponseTypesGetBuildSessionId).SetData(resp)
	if err := targetResp.Validate(); err != nil {
		return nil, err
	}
	return targetResp, nil
}

func (t *Target) doGetAgentConfiguration(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	resp := configuration_models.NewConfigurationResponse()
	request, ok := req.GetData().(*configuration_models.ConfigurationRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request model")
	}
	httpRequest := t.createRequest(ctx).
		SetResult(resp).
		SetQueryParams(request.ToQueryMap())
	httpResp, err := httpRequest.Get(t.cfg.getBaseUrl("/v3/config"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}

	targetResp := target.NewResponse().SetType(target.ResponseTypesGetAgentConfiguration).SetData(resp)
	if err := targetResp.Validate(); err != nil {
		return nil, err
	}
	return targetResp, nil
}

func (t *Target) doPostAgentLogSubmission(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	httpRequest := t.createRequest(ctx).
		SetBody(req.GetData())

	httpResp, err := httpRequest.Post(t.cfg.getBaseUrl("/v2/logsubmission"))
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}

	return nil, nil
}

func (t *Target) doGetClockSync(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	getClockSyncRequest, ok := req.GetData().(*clock_models.ClockSyncRequest)
	if !ok {
		return nil, fmt.Errorf("doGetClockSync: invalid request type")
	}
	url := strings.TrimSuffix(t.cfg.getBaseUrl(""), "/api")
	path := fmt.Sprintf("%s/clock/sync?time=%d", url, getClockSyncRequest.MyTime)
	getClockSyncResponse := clock_models.NewClockSyncResponse()
	httpRequest := t.createRequest(ctx).SetResult(getClockSyncResponse)

	httpResp, err := httpRequest.Get(path)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	if err := getClockSyncResponse.Validate(); err != nil {
		return nil, fmt.Errorf("doGetClockSync: invalid response:%s", err.Error())
	}
	return target.NewResponse().
		SetType(target.ResponseTypesClockSync).
		SetData(getClockSyncResponse), nil
}

func (t *Target) createRequest(ctx context.Context) *resty.Request {
	return t.client.R().
		SetAuthToken(t.cfg.Token).
		SetContext(ctx).SetHeaders(t.cfg.Headers)
}

func (t *Target) Close() error {
	t.canceler()
	return nil
}

func (t *Target) doGetTestRecommendations(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	getTestRecommendationsRequest, ok := req.GetData().(*tia_models.TestRecommendationsRequest)
	if !ok {
		return nil, fmt.Errorf("doGetTestRecommendations: invalid request type")
	}
	path := fmt.Sprintf("%s/v3/test-exclusions/%s/%s", t.cfg.getBaseUrl(""), getTestRecommendationsRequest.BuildSessionId, getTestRecommendationsRequest.Stage)
	if getTestRecommendationsRequest.TestGroupId != "" {
		path = fmt.Sprintf("%s/v4/test-exclusions/%s/%s/%s", t.cfg.getBaseUrl(""), getTestRecommendationsRequest.BuildSessionId, getTestRecommendationsRequest.Stage, getTestRecommendationsRequest.TestGroupId)
	}
	getTestRecommendationsResponse := tia_models.NewTestRecommendationsResponse()
	httpRequest := t.createRequest(ctx).SetResult(getTestRecommendationsResponse)
	httpResp, err := httpRequest.Get(path)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	if err := getTestRecommendationsResponse.Validate(); err != nil {
		return nil, fmt.Errorf("doGetTestRecommendations: invalid response:%s", err.Error())
	}
	return target.NewResponse().
		SetType(target.ResponseTypesTestRecommendations).
		SetData(getTestRecommendationsResponse), nil
}

func (t *Target) doGetBuildSessionIDByLabID(ctx context.Context, req target.TargetRequest) (target.TargetResponse, error) {
	getBuildSessionIdRequest, ok := req.GetData().(*buildsession_models.BuildSessionByLabIdRequest)
	if !ok {
		return nil, fmt.Errorf("doGetBuildSessionIDByLabID: invalid request type")
	}
	path := fmt.Sprintf("%s/v1/lab-ids/%s/build-sessions/active", t.cfg.getBaseUrl(""), getBuildSessionIdRequest.LabId)
	getBuildSessionIdResponse := buildsession_models.NewGetBuildSessionIdResponse()
	httpRequest := t.createRequest(ctx).SetResult(getBuildSessionIdResponse)
	httpResp, err := httpRequest.Get(path)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode() >= 400 {
		return nil, t.getHttpResponseError(httpResp)
	}
	if err := getBuildSessionIdResponse.Validate(); err != nil {
		return nil, fmt.Errorf("doGetBuildSessionIDByLabID: invalid response:%s", err.Error())
	}
	return target.NewResponse().
		SetType(target.ResponseTypesGetBuildSessionId).
		SetData(getBuildSessionIdResponse), nil
}
