// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// The aws-node ipam daemon binary
package main

import (
	"context"
	cniNode "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/v1alpha1"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/config"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/k8sapi"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/eventrecorder"
	"github.com/aws/amazon-vpc-cni-k8s/utils"
	metrics "github.com/aws/amazon-vpc-cni-k8s/utils/prometheusmetrics"
	"os"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/controllers"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/version"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/spf13/pflag"
)

const (
	appName = "aws-node"
	// metricsPort is the port for prometheus metrics
	metricsPort = 61678

	// Environment variable to disable the metrics endpoint on 61678
	envDisableMetrics = "DISABLE_METRICS"
)

func main() {
	os.Exit(_main())
}

var (
	scheme = runtime.NewScheme()
	//setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(cniNode.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func _main() int {
	// Do not add anything before initializing logger
	log := logger.Get()

	log.Infof("Starting L-IPAMD %s  ...", version.Version)
	version.RegisterMetric()

	// Create Kubernetes client for API server requests
	k8sClient, err := k8sapi.CreateKubeClient(appName)
	if err != nil {
		log.Errorf("Failed to create kube client: %s", err)
		return 1
	}

	// Create EventRecorder for use by IPAMD
	if err := eventrecorder.Init(k8sClient); err != nil {
		log.Errorf("Failed to create event recorder: %s", err)
		return 1
	}

	ctrlConfig, err := loadControllerConfig()
	if err != nil {
		//initLogger.Error("unable to load policy endpoint controller config")
		log.Errorf("loadControllerConfig failure: %v", err)
		os.Exit(1)
	}

	restCFG, err := config.BuildRestConfig(ctrlConfig.RuntimeConfig)
	if err != nil {
		//setupLog.Error(err, "unable to build REST config")
		log.Errorf("BuildRestConfig failure: %v", err)
		os.Exit(1)
	}

	runtimeOpts := config.BuildRuntimeOptions(ctrlConfig.RuntimeConfig, scheme)
	mgr, err := ctrl.NewManager(restCFG, runtimeOpts)
	if err != nil {
		log.Errorf("BuildRuntimeOptions failure: %v", err)
		//setupLog.Error(err, "unable to create controller manager")
		os.Exit(1)
	}

	ctx := context.Background() //ctrl.SetupSignalHandler()
	cniNodeReconciler, err := controllers.NewCNINodeReconciler(mgr.GetClient(), log, false)
	//	ctrl.Log.WithName("controllers").WithName("cniNode"), false)
	if err != nil {
		log.Errorf("unable to setup controller: %v", err)
		//setupLog.Error(err, "unable to setup controller", "controller", "PolicyEndpoints init failed")
		os.Exit(1)
	}

	if err = cniNodeReconciler.SetupWithManager(ctx, mgr); err != nil {
		log.Errorf("unable to create controller: %v", err)
		//setupLog.Error(err, "unable to create controller", "controller", "PolicyEndpoints")
		os.Exit(1)
	}

	if !utils.GetBoolAsStringEnvVar(envDisableMetrics, false) {
		// Prometheus metrics
		go metrics.ServeMetrics(metricsPort)
	}

	// TODO: Do we need it? - CNI introspection endpoints
	/*
		if !utils.GetBoolAsStringEnvVar(envDisableIntrospection, false) {
			go ipamContext.ServeIntrospection()
		}
	*/

	log.Info("starting manager")
	if err = mgr.Start(ctx); err != nil {
		log.Errorf("problem running manager", err)
		os.Exit(1)
	}
	return 0
}

// TODO - getLoggerWithLogLevel returns logger with specific log level.
/*
func getLoggerWithLogLevel(logLevel string, logFilePath string) (logr.Logger, error) {
	logConfig := logger.Configuration{
		LogLevel:    logLevel,
		LogLocation: logFilePath,
	}
	ctrlLogger := logger.New(&logConfig)
	//return ctrlLogger, nil
	return zapr.NewLogger(ctrlLogger), nil
}
*/

// loadControllerConfig loads the controller configuration
func loadControllerConfig() (config.ControllerConfig, error) {
	controllerConfig := config.ControllerConfig{}
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	controllerConfig.BindFlags(fs)

	if err := fs.Parse(os.Args); err != nil {
		return controllerConfig, err
	}

	return controllerConfig, nil
}
