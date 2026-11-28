package oam

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMeshmodelComponentsFiltersFiles(t *testing.T) {
	root := t.TempDir()

	previous := AvailableVersions
	AvailableVersions = map[string]bool{}
	defer func() {
		AvailableVersions = previous
	}()

	files := []string{
		"edge-26.6.3/linkerdmesh_meshmodel.json",
		"edge-26.6.3/serviceprofile.Linkerd.meshery.layer5.io_meshmodel.json",
		"edge-26.6.3/README.md",
		"edge-26.6.3/unrelated.json",
		"edge-26.7.1/temporary.txt",
	}

	for _, name := range files {
		path := filepath.Join(root, name)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	paths, err := loadMeshmodelComponents(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 component files, got %d", len(paths))
	}

	if !AvailableVersions["edge-26.6.3"] {
		t.Fatal("component version was not recorded")
	}

	if AvailableVersions["edge-26.7.1"] {
		t.Fatal("version without components was incorrectly recorded")
	}
}
