package installer

import (
	yaml "github.com/jcrossley3/manifestival/pkg/manifestival"
	"k8s.io/client-go/dynamic"
)

func NewInstaller(ns string, dc dynamic.Interface) *Installer {
	return &Installer{ns: ns, dc: dc}
}

type Installer struct {
	ns string
	dc dynamic.Interface

	Manifest *yaml.Manifest
}

func (r *Installer) Install() error {

	return nil
}
