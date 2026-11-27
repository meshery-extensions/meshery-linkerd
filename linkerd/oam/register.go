package oam

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

var (
	basePath, _ = os.Getwd()

	// MeshModelComponents Path will be used by both static and component generation
	MeshmodelComponents = filepath.Join(basePath, "templates", "meshmodel", "components")
)

// AvailableVersions denote the component versions available statically
var AvailableVersions = map[string]bool{}
var availableVersionGlobalMutex sync.Mutex

type meshmodelDefinitionPathSet struct {
	meshmodelDefinitionPath string
}

func RegisterMeshModelComponents(uuid, runtime, host, port string) error {
	pathSets, err := loadMeshmodelComponents(MeshmodelComponents)
	if err != nil {
		return err
	}

	paths := make([]string, 0, len(pathSets))
	for _, pathSet := range pathSets {
		paths = append(paths, pathSet.meshmodelDefinitionPath)
	}

	return adapter.RegisterModernMeshModelComponentFiles(
		context.Background(),
		uuid,
		runtime,
		host,
		port,
		paths,
	)
}

func loadMeshmodelComponents(basepath string) ([]meshmodelDefinitionPathSet, error) {
	res := []meshmodelDefinitionPathSet{}
	if err := filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() ||
			!strings.HasSuffix(info.Name(), "_meshmodel.json") {
			return nil
		}

		res = append(res, meshmodelDefinitionPathSet{
			meshmodelDefinitionPath: path,
		})
		availableVersionGlobalMutex.Lock()
		AvailableVersions[filepath.Base(filepath.Dir(path))] = true // Getting available versions already existing on file system
		availableVersionGlobalMutex.Unlock()
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}
