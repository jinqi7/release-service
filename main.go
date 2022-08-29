/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/kcp"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	hasv1alpha1 "github.com/redhat-appstudio/application-service/api/v1alpha1"
	appstudioshared "github.com/redhat-appstudio/managed-gitops/appstudio-shared/apis/appstudio.redhat.com/v1alpha1"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	appstudiov1alpha1 "github.com/redhat-appstudio/release-service/api/v1alpha1"
	"github.com/redhat-appstudio/release-service/controllers"
	kcpUtils "github.com/redhat-appstudio/release-service/kcp"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appstudiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(appstudioshared.AddToScheme(scheme))
	utilruntime.Must(hasv1alpha1.AddToScheme(scheme))
	utilruntime.Must(tektonv1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var apiExportName string
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var mgr ctrl.Manager
	var err error
	flag.StringVar(&apiExportName, "api-export-name", "", "The name of the APIExport.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()
	restConfig := ctrl.GetConfigOrDie()
	setupLog = setupLog.WithValues("api-export-name", apiExportName)

	options := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "f3d4c01a.redhat.com",
		LeaderElectionConfig:   restConfig,
	}
	if kcpUtils.IsKcpApiGroupPresent(restConfig, setupLog) {
		setupLog.Info("Looking up virtual workspace URL")
		cfg, err := kcpUtils.GetRestConfigForApiExport(ctx, restConfig, apiExportName, setupLog)
		if err != nil {
			setupLog.Error(err, "error looking up virtual workspace URL")
		}

		setupLog.Info("Using virtual workspace URL", "url", cfg.Host)

		options.LeaderElectionConfig = restConfig
		mgr, err = kcp.NewClusterAwareManager(cfg, options)
		if err != nil {
			setupLog.Error(err, "unable to start cluster aware manager")
			os.Exit(1)
		}
	} else {
		setupLog.Info("The apis.kcp.dev group is not present - creating standard manager")
		mgr, err = ctrl.NewManager(restConfig, options)
		if err != nil {
			setupLog.Error(err, "unable to start manager")
			os.Exit(1)
		}
	}

	// Set a default value for the DEFAULT_RELEASE_PVC environment variable
	if os.Getenv("DEFAULT_RELEASE_PVC") == "" {
		err := os.Setenv("DEFAULT_RELEASE_PVC", "release-pvc")
		if err != nil {
			setupLog.Error(err, "unable to setup DEFAULT_RELEASE_PVC environment variable")
			os.Exit(1)
		}
	}

	// Set a default value for the DEFAULT_RELEASE_WORKSPACE_NAME environment variable
	if os.Getenv("DEFAULT_RELEASE_WORKSPACE_NAME") == "" {
		err := os.Setenv("DEFAULT_RELEASE_WORKSPACE_NAME", "release-workspace")
		if err != nil {
			setupLog.Error(err, "unable to setup DEFAULT_RELEASE_WORKSPACE_NAME environment variable")
			os.Exit(1)
		}
	}

	err = controllers.SetupControllers(mgr)
	if err != nil {
		setupLog.Error(err, "unable to setup controllers")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		setupLog.Info("setting up webhooks")

		if err = (&appstudiov1alpha1.Release{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Release")
			os.Exit(1)
		}

		if err = (&appstudiov1alpha1.ReleasePlanAdmission{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ReleasePlanAdmission")
			os.Exit(1)
		}

		if err = (&appstudiov1alpha1.ReleasePlan{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ReleasePlan")
			os.Exit(1)
		}
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
