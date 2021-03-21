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
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"

	eniconfigscheme "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/v1alpha1"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var version string

func main() {
	os.Exit(_main())
}



func _main() int {
	//Do not add anything before initializing logger
	logConfig := logger.Configuration{
		LogLevel:    logger.GetLogLevel(),
		LogLocation: logger.GetLogLocation(),
	}
	log := logger.New(&logConfig)

	log.Infof("Starting L-IPAMD %s  ...", version)

	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	eniconfigscheme.AddToScheme(scheme)

	kubeConfig := ctrl.GetConfigOrDie()

	syncPeriod := 10*time.Hour

	mgr, err := ctrl.NewManager(kubeConfig, ctrl.Options{
		Scheme:                 scheme,
		SyncPeriod:             &syncPeriod,
		LeaderElection:         false,
		Port:                   9443,
	})

	stopChan := ctrl.SetupSignalHandler()
	go func() {
		mgr.GetCache().Start(stopChan)
	}()
	mgr.GetCache().WaitForCacheSync(stopChan)


	//ipamContext, err := ipamd.New(k8sclient, eniConfigController)
	ipamContext, err := ipamd.New(mgr.GetClient())
	if err != nil {
		log.Errorf("Initialization failure: %v", err)
		return 1
	}

	// Pool manager
	go ipamContext.StartNodeIPPoolManager()

	// Prometheus metrics
	go ipamContext.ServeMetrics()

	// CNI introspection endpoints
	go ipamContext.ServeIntrospection()

	// Start the RPC listener
	err = ipamContext.RunRPCHandler(version)
	if err != nil {
		log.Errorf("Failed to set up gRPC handler: %v", err)
		return 1
	}
	return 0
}
