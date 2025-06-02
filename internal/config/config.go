// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package config define the configuration and operation of the adapter
*/
package config

import (
	"path"
	"strings"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/walker"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

const (
	Development = "development"
	Production  = "production"

	AnnotateNamespace = "annotate-namespace"
	ServicePatchFile  = "service-patch-file"
	HelmChartURL      = "helm-chart-url"

	// Addons that the adapter supports
	JaegerAddon       = "jaeger-addon"
	VizAddon          = "viz-addon"
	MultiClusterAddon = "multicluster-addon"
	SMIAddon          = "smi-addon"
	// OAM Metadata constants
	OAMAdapterNameMetadataKey       = "adapter.meshery.io/name"
	OAMComponentCategoryMetadataKey = "ui.meshery.io/category"
)

var (
	// LinkerdOperation is the default name for the install
	// and uninstall commands on the Linkerd
	LinkerdOperation = strings.ToLower(smp.ServiceMesh_LINKERD.Enum().String())

	configRootPath = path.Join(utils.GetHome(), ".meshery")

	Config = configprovider.Options{
		FilePath: configRootPath,
		FileName: "linkerd",
		FileType: "yaml",
	}

	ServerConfig = map[string]string{
		"name":     smp.ServiceMesh_LINKERD.Enum().String(),
		"port":     "10001",
		"type":     "adapter",
		"traceurl": status.None,
	}

	MeshSpec = map[string]string{
		"name":    smp.ServiceMesh_LINKERD.Enum().String(),
		"status":  status.NotInstalled,
		"version": status.None,
	}

	ProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "linkerd",
	}

	// KubeConfig - Controlling the kubeconfig lifecycle with viper
	KubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	Operations = getOperations(common.Operations)
)

// New creates a new config instance
func New(provider string) (h config.Handler, err error) {
	// Config provider
	switch provider {
	case configprovider.ViperKey:
		h, err = configprovider.NewViper(Config)
		if err != nil {
			return nil, err
		}
	case configprovider.InMemKey:
		h, err = configprovider.NewInMem(Config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrEmptyConfig
	}

	// Setup server config
	if err := h.SetObject(adapter.ServerKey, ServerConfig); err != nil {
		return nil, err
	}

	// Setup mesh config
	if err := h.SetObject(adapter.MeshSpecKey, MeshSpec); err != nil {
		return nil, err
	}

	// Setup Operations Config
	if err := h.SetObject(adapter.OperationsKey, Operations); err != nil {
		return nil, err
	}

	return h, nil
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{
		FilePath: configRootPath,
		FileType: "yaml",
		FileName: "kubeconfig",
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}

	return nil, ErrEmptyConfig
}

// RootPath returns the config root path for the adapter
func RootPath() string {
	return configRootPath
}
func threadSafeAppend(fs *map[string]string, name string, url string, m *sync.RWMutex) {
	m.Lock()
	defer m.Unlock()
	(*fs)[name] = url
}

// GetFileNames takes the url of a github repo and the path to a directory. Then returns all the filenames->URL from that directory
func GetFileNames(owner string, repo string, path string) (map[string]string, error) {
	g := walker.NewGithub()
	var filenames = make(map[string]string)
	var m sync.RWMutex
	err := g.Owner(owner).Repo(repo).Root(path).RegisterFileInterceptor(func(gca walker.GithubContentAPI) error {
		threadSafeAppend(&filenames, gca.Name, gca.DownloadURL, &m)
		return nil
	}).Walk()
	if err != nil {
		return nil, ErrGetFileNames(err)
	}
	return filenames, nil
}
