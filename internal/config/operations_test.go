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
	"slices"
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

func TestGetOperations(t *testing.T) {
	dev := make(adapter.Operations)
	versions := []string{"stable-2.14.10", "stable-2.14.9"}
	dev = getOperationsWithVersions(dev, versions)

	t.Run("LinkerdOperation", func(t *testing.T) {
		assertLinkerdOperation(t, dev, versions)
	})

	t.Run("AnnotateNamespace", func(t *testing.T) {
		assertAnnotateNamespace(t, dev)
	})

	t.Run("JaegerAddon", func(t *testing.T) {
		assertJaegerAddon(t, dev)
	})

	t.Run("MultiClusterAddon", func(t *testing.T) {
		assertMultiClusterAddon(t, dev)
	})

	t.Run("SMIAddon", func(t *testing.T) {
		assertSMIAddon(t, dev)
	})
}

func assertLinkerdOperation(t *testing.T, dev adapter.Operations, versions []string) {
	t.Helper()
	op, exists := dev[LinkerdOperation]
	if !exists {
		t.Errorf("Operation %v not found in operations map", LinkerdOperation)
	}
	if op.Type != int32(meshes.OpCategory_INSTALL) {
		t.Errorf("Expected operation type %v but got %v", int32(meshes.OpCategory_INSTALL), op.Type)
	}
	if op.Description != "Linkerd Service Mesh" {
		t.Errorf("Expected operation description %v but got %v", "Linkerd Service Mesh", op.Description)
	}
	var adapterVersions []adapter.Version
	for _, v := range versions {
		adapterVersions = append(adapterVersions, adapter.Version(v))
	}
	if !slices.Equal(op.Versions, adapterVersions) {
		t.Errorf("Expected versions %v but got %v", adapterVersions, op.Versions)
	}
	if len(op.Templates) != 0 {
		t.Errorf("Expected empty templates slice but got %v templates", len(op.Templates))
	}
	if len(op.AdditionalProperties) != 0 {
		t.Errorf("Expected empty additional properties map but got %v properties", len(op.AdditionalProperties))
	}

}

func assertAnnotateNamespace(t *testing.T, dev adapter.Operations) {
	t.Helper()
	op, exists := dev[AnnotateNamespace]
	if !exists {
		t.Errorf("Operation %v not found in operations map", AnnotateNamespace)
	}
	if op.Type != int32(meshes.OpCategory_CONFIGURE) {
		t.Errorf("Expected operation type %v but got %v", int32(meshes.OpCategory_CONFIGURE), op.Type)
	}
	if op.Description != "Annotate Namespace" {
		t.Errorf("Expected operation description %v but got %v", "Annotate Namespace", op.Description)
	}

}

func assertJaegerAddon(t *testing.T, dev adapter.Operations) {
	t.Helper()
	op, exists := dev[JaegerAddon]
	if !exists {
		t.Errorf("Operation %v not found in operations map", JaegerAddon)
	}
	if op.Type != int32(meshes.OpCategory_CONFIGURE) {
		t.Errorf("Expected operation type %v but got %v", int32(meshes.OpCategory_CONFIGURE), op.Type)
	}
	if op.Description != "Add-on: Jaeger" {
		t.Errorf("Expected operation description %v but got %v", "Add-on: Jaeger", op.Description)
	}
	if len(op.AdditionalProperties) != 3 {
		t.Errorf("Expected 3 additional properties but got %v properties", len(op.AdditionalProperties))
	}
	if op.AdditionalProperties[ServiceName] != "jaeger" {
		t.Errorf("Expected %v for %v additional property but got %v", "jaeger", ServiceName, op.AdditionalProperties[ServiceName])
	}
	if op.AdditionalProperties[HelmChartURL] != "https://helm.linkerd.io/stable/linkerd-jaeger-30.4.5.tgz" {
		t.Errorf("Expected %v for %v additional property but got %v", "https://helm.linkerd.io/stable/linkerd-jaeger-30.4.5.tgz", HelmChartURL, op.AdditionalProperties[HelmChartURL])
	}

}

func assertMultiClusterAddon(t *testing.T, dev adapter.Operations) {
	t.Helper()
	op, exists := dev[MultiClusterAddon]
	if !exists {
		t.Errorf("Operation %v not found in operations map", MultiClusterAddon)
	}
	if op.Type != int32(meshes.OpCategory_CONFIGURE) {
		t.Errorf("Expected operation type %v but got %v", int32(meshes.OpCategory_CONFIGURE), op.Type)
	}
	if op.Description != "Add-on: Multi-cluster" {
		t.Errorf("Expected operation description %v but got %v", "Add-on: Multi-cluster", op.Description)
	}
	if len(op.AdditionalProperties) != 3 {
		t.Errorf("Expected 3 additional properties but got %v properties", len(op.AdditionalProperties))
	}
	if op.AdditionalProperties[ServiceName] != "linkerd-gateway" {
		t.Errorf("Expected %v for %v additional property but got %v", "linkerd-gateway", ServiceName, op.AdditionalProperties[ServiceName])
	}
	if op.AdditionalProperties[HelmChartURL] != "https://helm.linkerd.io/stable/linkerd-multicluster-2.10.2.tgz" {
		t.Errorf("Expected %v for %v additional property but got %v", "https://helm.linkerd.io/stable/linkerd-multicluster-2.10.2.tgz", HelmChartURL, op.AdditionalProperties[HelmChartURL])
	}

}

func assertSMIAddon(t *testing.T, dev adapter.Operations) {
	t.Helper()
	op, exists := dev[SMIAddon]
	if !exists {
		t.Errorf("Operation %v not found in operations map", SMIAddon)
	}
	if op.Type != int32(meshes.OpCategory_CONFIGURE) {
		t.Errorf("Expected operation type %v but got %v", int32(meshes.OpCategory_CONFIGURE), op.Type)
	}
	if op.Description != "Add-on: SMI Addon" {
		t.Errorf("Expected operation description %v but got %v", "Add-on: SMI Addon", op.Description)
	}
	if len(op.AdditionalProperties) != 1 {
		t.Errorf("Expected 1 additional properties but got %v properties", len(op.AdditionalProperties))
	}
	if op.AdditionalProperties[HelmChartURL] != "https://github.com/linkerd/linkerd-smi/releases/download/v0.1.0/linkerd-smi-0.1.0.tgz" {
		t.Errorf("Expected %v for %v additional property but got %v", "https://github.com/linkerd/linkerd-smi/releases/download/v0.1.0/linkerd-smi-0.1.0.tgz", HelmChartURL, op.AdditionalProperties[HelmChartURL])
	}
}
