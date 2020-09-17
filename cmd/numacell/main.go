package main

import (
	"flag"

	"github.com/davecgh/go-spew/spew"
	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	flag.Parse()

	devs, _ := numacell.GetNUMACellDevices("/sys") // TODO
	glog.Infof("detected: %s", spew.Sdump(devs))

	manager := dpm.NewManager(numacell.NUMACellLister{})
	manager.Run()
}
