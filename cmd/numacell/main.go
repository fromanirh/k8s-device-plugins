package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/jaypipes/ghw/pkg/option"
	"github.com/jaypipes/ghw/pkg/topology"

	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func summarize(topoInfo *topology.Info) string {
	var buf strings.Builder
	for _, node := range topoInfo.Nodes {
		fmt.Fprintf(&buf, "NUMA node %d\n", node.ID)
		for _, core := range node.Cores {
			fmt.Fprintf(&buf, "\t%s\n", core.String())
		}
	}
	return buf.String()
}

func main() {
	var sysfsPath string
	flag.StringVar(&sysfsPath, "sysfs", "/sys", "mount path of sysfs")
	flag.Parse()

	glog.Infof("using sysfs at %q", sysfsPath)
	topoInfo, err := topology.New(option.WithPathOverrides(option.PathOverrides{
		"/sys": sysfsPath,
	}))
	if err != nil {
		log.Fatalf("error getting topology info from %q: %v", sysfsPath, err)
	}

	glog.Infof("hardware detected:\n%s", summarize(topoInfo))

	manager := dpm.NewManager(numacell.NewNUMACellLister(topoInfo))
	manager.Run()
}
