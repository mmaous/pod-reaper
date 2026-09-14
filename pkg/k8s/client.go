package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(kubeconfigPath string) (kubernetes.Interface, error) {
	var restCfg *rest.Config
	var err error
	if kubeconfigPath != "" {
		restCfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		restCfg, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("build kube client config: %w", err)
	}
	return kubernetes.NewForConfig(restCfg)
}
