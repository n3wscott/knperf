package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// BuildClientConfig builds the client config specified by the config path and the cluster name
func BuildClientConfigOrDie(kubeConfigPath string, clusterName string) *rest.Config {

	if cfg, err := clientcmd.BuildConfigFromFlags("", ""); err == nil {
		// success!
		return cfg
	}
	// try local...

	overrides := clientcmd.ConfigOverrides{}
	// Override the cluster name if provided.
	if clusterName != "" {
		overrides.Context.Cluster = clusterName
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&overrides).ClientConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
