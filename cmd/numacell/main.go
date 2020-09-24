package main

import (
	"flag"
	"log"
	"strings"
	"text/tabwriter"

	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
	"github.com/fromanirh/numalign/pkg/topologyinfo/cpus"
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func summarize(cpuInfos *cpus.CPUs) string {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	cpus.MakeSummary(cpuInfos, w)
	w.Flush()
	return buf.String()
}

func main() {
	var sysfsPath string
	flag.StringVar(&sysfsPath, "sysfs", "/sys", "mount path of sysfs")
	flag.Parse()

	cpuInfos, err := cpus.NewCPUs(sysfsPath)
	if err != nil {
		log.Fatalf("error getting topology info from %q: %v", sysfsPath, err)
	}

	glog.Infof("detected:\n%s", summarize(cpuInfos))

	manager := dpm.NewManager(numacell.NewNUMACellLister(cpuInfos))
	manager.Run()
}
