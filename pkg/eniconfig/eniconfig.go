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

// Package eniconfig handles eniconfig CRD
package eniconfig

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/apis/v1alpha1"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"
	"github.com/pkg/errors"
)

const (
	defaultEniConfigAnnotationDef = "k8s.amazonaws.com/eniConfig"
	defaultEniConfigLabelDef      = "k8s.amazonaws.com/eniConfig"
	eniConfigDefault              = "default"

	// when "ENI_CONFIG_LABEL_DEF is defined, ENIConfigController will use that label key to
	// search if is setting value for eniConfigLabelDef
	// Example:
	//   Node has set label k8s.amazonaws.com/eniConfigCustom=customConfig
	//   We can get that value in controller by setting environmental variable ENI_CONFIG_LABEL_DEF
	//   ENI_CONFIG_LABEL_DEF=k8s.amazonaws.com/eniConfigOverride
	//   This will set eniConfigLabelDef to eniConfigOverride
	envEniConfigAnnotationDef = "ENI_CONFIG_ANNOTATION_DEF"
	envEniConfigLabelDef      = "ENI_CONFIG_LABEL_DEF"
)

// ENIConfig interface
type ENIConfig interface {
	MyENIConfig(client.Client) (*v1alpha1.ENIConfigSpec, error)
    GetENIConfigName(context.Context, client.Client) (string, error)
}

// ErrNoENIConfig is the missing ENIConfig error
var ErrNoENIConfig = errors.New("eniconfig: eniconfig is not available")

var log = logger.Get()

// ENIConfigController defines global context for ENIConfig controller
type ENIConfigController struct {
	eni                    map[string]*v1alpha1.ENIConfigSpec
	myENI                  string
	eniLock                sync.RWMutex
	myNodeName             string
	eniConfigAnnotationDef string
	eniConfigLabelDef      string
}

// ENIConfigInfo returns locally cached ENIConfigs
type ENIConfigInfo struct {
	ENI                    map[string]v1alpha1.ENIConfigSpec
	MyENI                  string
	EniConfigAnnotationDef string
	EniConfigLabelDef      string
}


// NewENIConfigController creates a new ENIConfig controller
func NewENIConfigController() *ENIConfigController {
	return &ENIConfigController{
		myNodeName:             os.Getenv("MY_NODE_NAME"),
		eni:                    make(map[string]*v1alpha1.ENIConfigSpec),
		myENI:                  eniConfigDefault,
		eniConfigAnnotationDef: getEniConfigAnnotationDef(),
		eniConfigLabelDef:      getEniConfigLabelDef(),
	}
}

/*
// NewHandler creates a new handler for sdk
func NewHandler(controller *ENIConfigController) sdk.Handler {
	return &Handler{controller: controller}
}

// Handler stores the ENIConfigController
type Handler struct {
	controller *ENIConfigController
}

// Handle handles ENIConfig updates from API Server and store them in local cache
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.ENIConfig:
		eniConfigName := o.GetName()
		if event.Deleted {
			log.Debugf("Deleting ENIConfig: %s", eniConfigName)
			h.controller.eniLock.Lock()
			defer h.controller.eniLock.Unlock()
			delete(h.controller.eni, eniConfigName)
			return nil
		}

		curENIConfig := o.DeepCopy()

		log.Debugf("Handle ENIConfig Add/Update: %s, %v, %s", eniConfigName, curENIConfig.Spec.SecurityGroups, curENIConfig.Spec.Subnet)

		h.controller.eniLock.Lock()
		defer h.controller.eniLock.Unlock()
		h.controller.eni[eniConfigName] = &curENIConfig.Spec

	case *corev1.Node:
		log.Debugf("Handle corev1.Node: %s, %v, %v", o.GetName(), o.GetAnnotations(), o.GetLabels())
		// Get annotations if not found get labels if not found fallback use default
		if h.controller.myNodeName == o.GetName() {
			val, ok := o.GetAnnotations()[h.controller.eniConfigAnnotationDef]
			if !ok {
				val, ok = o.GetLabels()[h.controller.eniConfigLabelDef]
				if !ok {
					val = eniConfigDefault
				}
			}
			// If value changes
			if h.controller.myENI != val {
				h.controller.eniLock.Lock()
				defer h.controller.eniLock.Unlock()
				h.controller.myENI = val
				log.Debugf("Setting myENI to: %s", val)
				if val != eniConfigDefault {
					labels := o.GetLabels()
					labels["vpc.amazonaws.com/eniConfig"] = val
					o.SetLabels(labels)
				}
			}
		}
	}
	return nil
}

func printVersion() {
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

 */

// Start kicks off ENIConfig controller
/*
func (eniCfg *ENIConfigController) Start() {
	printVersion()

	sdk.ExposeMetricsPort()

	resource := "crd.k8s.amazonaws.com/v1alpha1"
	kind := "ENIConfig"
	resyncPeriod := time.Hour * 10
	log.Infof("Watching %s, %s, every %v s", resource, kind, resyncPeriod.Seconds())
	sdk.Watch(resource, kind, "", resyncPeriod)
	sdk.Watch("/v1", "Node", corev1.NamespaceAll, resyncPeriod)
	sdk.Handle(NewHandler(eniCfg))
	sdk.Run(context.TODO())
}


func (eniCfg *ENIConfigController) Getter() *ENIConfigInfo {
	output := &ENIConfigInfo{
		ENI: make(map[string]v1alpha1.ENIConfigSpec),
	}
	eniCfg.eniLock.Lock()
	defer eniCfg.eniLock.Unlock()

	output.MyENI = eniCfg.myENI
	output.EniConfigAnnotationDef = getEniConfigAnnotationDef()
	output.EniConfigLabelDef = getEniConfigLabelDef()

	for name, val := range eniCfg.eni {
		output.ENI[name] = *val
	}
	return output
}
*/

// MyENIConfig returns the ENIConfig applicable to the particular node
func MyENIConfig(k8sClient client.Client) (*v1alpha1.ENIConfigSpec, error) {

	ctx := context.Background()
	eniConfigName, err := GetENIConfigName(ctx, k8sClient)
	if err != nil {
       log.Debugf("Error while retrieving Node name")
	}

	log.Debugf("Found ENI Config Name: %s", eniConfigName)
	log.Debugf("Let's get ENIConfig List Info via Manager Client - Cache Synced - New")

	eniConfigsList := v1alpha1.ENIConfigList{}
	err = k8sClient.List(ctx, &eniConfigsList)
	if err != nil {
		fmt.Errorf("Error while EniConfig List Get: %s", err)
	}
	log.Debugf("ENIConfigs Size: %s ", len(eniConfigsList.Items))
	for _, eni := range eniConfigsList.Items {
		if eniConfigName == eni.Name {
			log.Debugf("Matching ENIConfig found: %s - %s - %s ", eni.Name, eni.Spec.Subnet, eni.Spec.SecurityGroups)
			return &v1alpha1.ENIConfigSpec{
				SecurityGroups: eni.Spec.SecurityGroups,
				Subnet:         eni.Spec.Subnet,
			}, nil
		}
		log.Debugf("ENIConfigs Info: %s - %s - %s ", eni.Name, eni.Spec.Subnet, eni.Spec.SecurityGroups)
	}

	return nil, ErrNoENIConfig
}

// getEniConfigAnnotationDef returns eniConfigAnnotation
func getEniConfigAnnotationDef() string {
	inputStr, found := os.LookupEnv(envEniConfigAnnotationDef)

	if !found {
		return defaultEniConfigAnnotationDef
	}
	if len(inputStr) > 0 {
		log.Debugf("Using ENI_CONFIG_ANNOTATION_DEF %v", inputStr)
		return inputStr
	}
	return defaultEniConfigAnnotationDef
}

// getEniConfigLabelDef returns eniConfigLabel name
func getEniConfigLabelDef() string {
	inputStr, found := os.LookupEnv(envEniConfigLabelDef)

	if !found {
		return defaultEniConfigLabelDef
	}
	if len(inputStr) > 0 {
		log.Debugf("Using ENI_CONFIG_LABEL_DEF %v", inputStr)
		return inputStr
	}
	return defaultEniConfigLabelDef
}

func GetENIConfigName(ctx context.Context, k8sClient client.Client) (string, error) {
    var eniConfigName string
	log.Debugf("Let's get Node List Info via Manager Client")
	nodeList := corev1.NodeList{}
	err := k8sClient.List(ctx, &nodeList)
	if err != nil {
		fmt.Errorf("Error while Node List Get: %s", err)
	}
	log.Debugf("Node Count: ", len(nodeList.Items))
	for _, node := range nodeList.Items {
		if node.Name == os.Getenv("MY_NODE_NAME") {
			log.Debugf("Node Info: %s", node.Name)
			val, ok := node.GetAnnotations()[getEniConfigAnnotationDef()]
			if !ok {
				val, ok = node.GetLabels()[getEniConfigLabelDef()]
				if !ok {
					val = eniConfigDefault
				}
			}

			eniConfigName = val
			if val != eniConfigDefault {
				labels := node.GetLabels()
				labels["vpc.amazonaws.com/eniConfig"] = eniConfigName
				node.SetLabels(labels)
				//k8sClient.Update(ctx, )
			}
		}
	}
	log.Debugf("ENI Config Name: %s", eniConfigName)
	return eniConfigName, nil
}
