package agent_models

import (
	"os"
	"runtime"

	"github.com/konflux-ci/release-service/__sealights__/models"
)

var awsEnvVars = []string{
	"_HANDLER",
	"_X_AMZN_TRACE_ID",
	"AWS_REGION",
	"AWS_EXECUTION_ENV",
	"AWS_LAMBDA_FUNCTION_NAME",
	"AWS_LAMBDA_FUNCTION_MEMORY_SIZE",
	"AWS_LAMBDA_FUNCTION_VERSION",
	"AWS_LAMBDA_FUNCTION_VERSION",
	"AWS_LAMBDA_INITIALIZATION_TYPE",
	"AWS_LAMBDA_RUNTIME_API",
	"LAMBDA_TASK_ROOT",
	"LAMBDA_RUNTIME_DIR",
	"TZ",
	"LANG",
	"LD_LIBRARY_PATH",
}

var gcpEnvVars = []string{
	"ENTRY_POINT",
	"GCP_PROJECT",
	"GCLOUD_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"FUNCTION_TRIGGER_TYPE",
	"FUNCTION_NAME",
	"FUNCTION_MEMORY_MB",
	"FUNCTION_TIMEOUT_SEC",
	"FUNCTION_IDENTITY",
	"FUNCTION_REGION",
	"FUNCTION_TARGET",
	"FUNCTION_SIGNATURE_TYPE",
	"K_SERVICE",
	"K_REVISION",
	"PORT",
}

var azureEnvVars = []string{
	"WEBSITE_FUNCTIONS_ARMCACHE_ENABLED",
	"WEBSITE_MAX_DYNAMIC_APPLICATION_SCALE_OUT",
	"FUNCTIONS_EXTENSION_VERSION",
	"AzureWebJobsSecretStorageType",
	"AzureWebJobsStorage",
	"WEBSITE_CONTENTAZUREFILECONNECTIONSTRING",
	"WEBSITE_CONTENTSHARE",
	"WEBSITE_CONTENTOVERVNET",
	"WEBSITE_ENABLE_BROTLI_ENCODING",
	"WEBSITE_USE_PLACEHOLDER",
	"WEBSITE_PLACEHOLDER_MODE",
	"WEBSITE_DISABLE_ZIP_CACHE",
}

type TechSpecificInfo struct {
	GoVersion       string            `json:"goVersion"`
	Goos            string            `json:"goOS"`
	Goarch          string            `json:"goArch"`
	GoMaxProcess    int               `json:"goMaxProcs"`
	GoNumGoroutine  int               `json:"goNumGoroutine"`
	GoNumCPU        int               `json:"goNumCPU"`
	AWSRuntimeEnv   map[string]string `json:"awsRuntimeEnv,omitempty"`
	GCPRuntimeEnv   map[string]string `json:"gcpRuntimeEnv,omitempty"`
	AzureRuntimeEnv map[string]string `json:"azureRuntimeEnv,omitempty"`
}

func NewTechSpecificInfo() *TechSpecificInfo {
	tsi := &TechSpecificInfo{
		GoVersion:       runtime.Version(),
		Goos:            runtime.GOOS,
		Goarch:          runtime.GOARCH,
		GoMaxProcess:    runtime.GOMAXPROCS(0),
		GoNumGoroutine:  runtime.NumGoroutine(),
		GoNumCPU:        runtime.NumCPU(),
		AWSRuntimeEnv:   getEnvMap(awsEnvVars),
		GCPRuntimeEnv:   getEnvMap(gcpEnvVars),
		AzureRuntimeEnv: getEnvMap(azureEnvVars),
	}
	return tsi
}

func getEnvMap(vars []string) map[string]string {
	m := map[string]string{}
	for _, v := range vars {
		if val, ok := os.LookupEnv(v); ok {
			m[v] = val
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func (i *TechSpecificInfo) Validate() error {
	return nil
}

var _ models.Model = &TechSpecificInfo{}
