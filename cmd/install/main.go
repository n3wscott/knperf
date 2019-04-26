package main

import (
	"flag"
	yaml "github.com/jcrossley3/manifestival/pkg/manifestival"
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

var ns = "default"

func main() {
	flag.Parse()

	cfg := config.BuildClientConfigOrDie(kubeconfig, cluster)
	dc := dynamic.NewForConfigOrDie(cfg)
	i := installer.NewInstaller(ns, dc)

	mf, err := yaml.NewYamlManifest("/var/run/ko/install", true, dc)

	if err != nil {
		log.Fatalf("could not install: %v", err)
	}
	i.Manifest = &mf

	if err := i.Install(); err != nil {
		log.Fatalf("could not install: %v", err)
	}

	os.Exit(0)
}
