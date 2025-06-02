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

package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshkit/utils"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {
	var adapterVersions []adapter.Version

	versions, _ := utils.GetLatestReleaseTagsSorted("linkerd", "linkerd2")
	for _, v := range versions {
		adapterVersions = append(adapterVersions, adapter.Version(v))
	}
	dev[LinkerdOperation] = &adapter.Operation{
		Type:                 int32(meshes.OpCategory_INSTALL),
		Description:          "Linkerd Service Mesh",
		Versions:             adapterVersions,
		Templates:            []adapter.Template{},
		AdditionalProperties: map[string]string{},
	}

	dev[AnnotateNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Annotate Namespace",
	}
	dev[JaegerAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Jaeger",
		AdditionalProperties: map[string]string{
			ServiceName:      "jaeger",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-jaeger-30.4.5.tgz",
		},
	}
	dev[VizAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Viz",
		AdditionalProperties: map[string]string{
			ServiceName:      "web",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-viz-30.3.5.tgz",
		},
	}
	dev[MultiClusterAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Multi-cluster",
		AdditionalProperties: map[string]string{
			ServiceName:      "linkerd-gateway",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-multicluster-2.10.2.tgz",
		},
	}
	dev[SMIAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: SMI Addon",
		AdditionalProperties: map[string]string{
			// ServiceName:      "linkerd-gateway",
			// ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL: "https://github.com/linkerd/linkerd-smi/releases/download/v0.1.0/linkerd-smi-0.1.0.tgz",
		},
	}
	return dev
}
