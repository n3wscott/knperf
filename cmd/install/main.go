package main

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/knperf/pkg/config"
	"github.com/n3wscott/knperf/pkg/installer"
	"k8s.io/client-go/dynamic"
	"log"
	"os"
	"os/user"
	"path"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	cluster    string
	kubeconfig string
)

type envConfig struct {
	Action    string `envconfig:"ACTION" required:"true"`
	Namespace string `envconfig:"NAMESPACE" required:"true"`
}

// TODO: move all this to ENV VARS

func init() {
	flag.StringVar(&cluster, "cluster", "",
		"Provide the cluster to test against. Defaults to the current cluster in kubeconfig.")

	var defaultKubeconfig string
	if usr, err := user.Current(); err == nil {
		defaultKubeconfig = path.Join(usr.HomeDir, ".kube/config")
	}

	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig,
		"Provide the path to the `kubeconfig` file.")

}

func main() {
	flag.Parse()

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	cfg := config.BuildClientConfigOrDie(kubeconfig, cluster)
	i := installer.NewInstaller(env.Namespace, dynamic.NewForConfigOrDie(cfg))

	if err := i.Do(env.Action); err != nil {
		log.Fatalf("could not install: %v", err)
	}

	os.Exit(0)
}
