package main

import (
	"flag"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"

	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
)

func main() {
	flag.Parse()

	manager := dpm.NewManager(numacell.NUMACellLister{})
	manager.Run()
}
