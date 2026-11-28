package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/meshkit/utils"

	"github.com/layer5io/meshery-linkerd/internal/config"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var CRDnamesURL map[string]string
var LatestVersion string
var WorkloadPath string
var MeshModelPath string
var AllVersions []string

const Component = "Linkerd"

var Meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Cloud Native Network",
	Metadata: Meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        strings.ToLower(smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_LINKERD)]),
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema"}, false),
		ExtractCrds: func(manifest string) []string {
			manifests.RemoveHelmTemplatingFromCRD(&manifest)
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	// Initialize Metadata including logo svgs
	byt, err := os.ReadFile("./build/meshmodel_metadata.json")
	if err != nil {
		fmt.Printf("Error reading MeshModel metadata: %v\n", err)
	} else if err := json.Unmarshal(byt, &Meshmodelmetadata); err != nil {
		fmt.Printf("Error decoding MeshModel metadata: %v\n", err)
	}

	wd, _ := os.Getwd()
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, err = utils.GetLatestReleaseTagsSorted("linkerd", "linkerd2")
	if err != nil {
		fmt.Printf("Linkerd release discovery failed: %v\n", err)
		return
	}
	if len(AllVersions) == 0 {
		return
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
	names, err := config.GetFileNames("linkerd", "linkerd2", "charts/linkerd-crds/templates/**")
	if err != nil {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for n := range names {
		if !strings.HasSuffix(n, ".yaml") {
			delete(names, n)
		}
	}
	CRDnamesURL = names
}
