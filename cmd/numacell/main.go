package main

import (
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
	"github.com/fromanirh/numalign/pkg/topologyinfo/cpus"
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	flag.Parse()
	sysfsPath := "/sys" // TODO

	cpuInfos, err := cpus.NewCPUs(sysfsPath)
	if err != nil {
		log.Fatalf("error getting topology info from %q: %v", sysfsPath, err)
	}
	glog.Infof("detected: %s", spew.Sdump(cpuInfos))

	manager := dpm.NewManager(numacell.NewNUMACellLister(cpuInfos))
	manager.Run()
}
