// Copyright 2022 Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/layer5io/meshery-linkerd/build"
	"github.com/layer5io/meshery-linkerd/linkerd"
	"github.com/layer5io/meshery-linkerd/linkerd/oam"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/events"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-linkerd/internal/config"
	configprovider "github.com/layer5io/meshkit/config/provider"
)

var (
	serviceName = "linkerd-adapter"
	version     = "edge"
	gitsha      = "none"
	instanceID  = uuid.NewString()
)

func init() {
	// Create the config path if it doesn't exists as the entire adapter
	// expects that directory to exists, which may or may not be true
	if err := os.MkdirAll(path.Join(config.RootPath(), "bin"), 0750); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// main is the entrypoint of the adapter
func main() {
	// Initialize Logger instance
	log, err := logger.New(serviceName, logger.Options{
		Format:     logger.SyslogLogFormat,
		DebugLevel: isDebug(),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set $KUBECONFIG environmental variable
	// crucial when adapter's running within the containers
	err = os.Setenv("KUBECONFIG", path.Join(
		config.KubeConfig[configprovider.FilePath],
		fmt.Sprintf("%s.%s", config.KubeConfig[configprovider.FileName], config.KubeConfig[configprovider.FileType])),
	)
	if err != nil {
		// Fail silently
		log.Warn(err)
	}

	// Initialize application specific configs and dependencies
	// App and request config
	cfg, err := config.New(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	err = cfg.GetObject(adapter.ServerKey, service)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	kubeconfigHandler, err := config.NewKubeconfigBuilder(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// // Initialize Tracing instance
	// tracer, err := tracing.New(service.Name, service.TraceURL)
	// if err != nil {
	// 	log.Err("Tracing Init Failed", err.Error())
	// 	os.Exit(1)
	// }
	e := events.NewEventStreamer()
	// Initialize Handler intance
	handler := linkerd.New(cfg, log, kubeconfigHandler, e)
	handler = adapter.AddLogger(log, handler)
	service.EventStreamer = e
	service.Handler = handler
	service.StartedAt = time.Now()
	service.Version = version
	service.GitSHA = gitsha

	go registerCapabilities(service.Port, log)
	go registerDynamicCapabilities(service.Port, log)
	// Server Initialization
	log.Info("Adapter Listening at port: ", service.Port)
	err = grpc.Start(service)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func mesheryServerAddress() string {
	meshReg := os.Getenv("MESHERY_SERVER")

	if meshReg != "" {
		if strings.HasPrefix(meshReg, "http") {
			return meshReg
		}

		return "http://" + meshReg
	}

	return "http://localhost:9081"
}

func serviceAddress() string {
	svcAddr := os.Getenv("SERVICE_ADDR")

	if svcAddr != "" {
		return svcAddr
	}

	return "localhost"
}

func registerCapabilities(port string, log logger.Handler) {
	// Register meshmodel components
	if err := oam.RegisterMeshModelComponents(instanceID, mesheryServerAddress(), serviceAddress(), port); err != nil {
		log.Error(err)
	}
}
func registerDynamicCapabilities(port string, log logger.Handler) {
	registerWorkloads(port, log)
	// Start the ticker
	const reRegisterAfter = 24
	ticker := time.NewTicker(reRegisterAfter * time.Hour)
	for {
		<-ticker.C
		registerWorkloads(port, log)
	}
}

func registerWorkloads(port string, log logger.Handler) {
	log.Info("Generating latest service mesh components...")

	// First we create and store any new components if available
	version := build.LatestVersion
	url := build.DefaultGenerationURL
	gm := build.DefaultGenerationMethod

	// Prechecking to skip comp gen
	if os.Getenv("FORCE_DYNAMIC_REG") != "true" && oam.AvailableVersions[version] {
		log.Info("Components available statically for version ", version, ". Skipping dynamic component registeration")
		return
	}

	//If a URL is passed from env variable, it will be used for component generation with default method being "using manifests"
	// In case a helm chart URL is passed, COMP_GEN_METHOD env variable should be set to Helm otherwise the component generation fails
	if os.Getenv("COMP_GEN_URL") != "" && (os.Getenv("COMP_GEN_METHOD") == "Helm" || os.Getenv("COMP_GEN_METHOD") == "Manifest") {
		url = os.Getenv("COMP_GEN_URL")
		gm = os.Getenv("COMP_GEN_METHOD")
		log.Info("Registering workload components from url ", url, " using ", gm, " method...")
	}
	log.Info("url: ", url, " method: ", gm)
	log.Info("Registering latest workload components for version ", version)

	for name, u := range build.CRDnamesURL {
		if err := adapter.CreateComponents(adapter.StaticCompConfig{
			URL:             u,
			Method:          gm,
			MeshModelPath:   build.MeshModelPath,
			MeshModelConfig: build.MeshModelConfig,
			DirName:         version,
			Config:          build.NewConfig(version),
		}); err != nil {
			log.Error(err)
			return
		}
		log.Info(name, " registered")
	}

	// The below log is checked in the workflows. If you change this log, reflect that change in the workflow where components are generated
	log.Info("Component creation completed for version ", version)

	// Now we will register in case
	log.Info("Registering workloads with Meshery Server for version ", version)
	if err := oam.RegisterMeshModelComponents(instanceID, mesheryServerAddress(), serviceAddress(), port); err != nil {
		log.Error(err)
		return
	}
	log.Info("Latest workload components successfully registered for version ", version)
}
// test change
// dummy test change
// test labeler change
// test change
