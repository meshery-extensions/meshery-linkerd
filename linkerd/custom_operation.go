package linkerd

import (
	"github.com/layer5io/meshery-adapter-library/status"
)

func (linkerd *Linkerd) applyCustomOperation(namespace string, manifest string, isDel bool, kubeconfigs []string) (string, error) {
	st := status.Starting

	err := linkerd.applyManifest([]byte(manifest), isDel, namespace, kubeconfigs)
	if err != nil {
		return st, ErrCustomOperation(err)
	}

	return status.Completed, nil  
}
