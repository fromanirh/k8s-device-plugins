package numacell

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fromanirh/numalign/pkg/topologyinfo/cpus"
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	NUMACellPath      = "/dev/null"
	NUMACellName      = "numacell"
	resourceNamespace = "devices.openshift-kni.io" // TODO pick a better one?
)

// NUMACellLister is the object responsible for discovering initial pool of devices and their allocation.
type NUMACellLister struct{}

type message struct{}

// NUMACellDevicePlugin is an implementation of DevicePlugin that is capable of exposing devices to containers.
type NUMACellDevicePlugin struct {
	update chan message
}

func (NUMACellLister) GetResourceNamespace() string {
	return resourceNamespace
}

// Discovery discovers all NUMA cells within the system.
func (NUMACellLister) Discover(pluginListCh chan dpm.PluginNameList) {
	pluginListCh <- dpm.PluginNameList{"numacell"}
}

// NewPlugin initializes new device plugin with NUMACell specific attributes.
func (NUMACellLister) NewPlugin(deviceID string) dpm.PluginInterface {
	glog.V(3).Infof("Creating device plugin %s", deviceID)
	return &NUMACellDevicePlugin{
		update: make(chan message),
	}
}

// ListAndWatch sends gRPC stream of devices.
func (dpi *NUMACellDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	cpuRes, err := cpus.NewCPUs("/sys") // TODO get from cmdline flags
	if err != nil {
		return err
	}

	// Send initial list of devices
	devs := make([]*pluginapi.Device, 0)
	resp := new(pluginapi.ListAndWatchResponse)
	for _, numacell := range cpuRes.NUMANodes {
		// Initialize with one available device
		devs = append(devs, &pluginapi.Device{
			ID:     fmt.Sprintf("%s%02d", NUMACellName, numacell),
			Health: pluginapi.Healthy,
			Topology: &pluginapi.TopologyInfo{
				Nodes: []*pluginapi.NUMANode{
					&pluginapi.NUMANode{
						ID: int64(numacell),
					},
				},
			},
		})
	}
	resp.Devices = devs
	glog.Infof("send devices %v\n", resp)

	if err := s.Send(resp); err != nil {
		glog.Errorf("failed to list NUMA cells: %v\n", err)
		return err
	}

	// TODO handle signals like sriovdp does
	for {
		select {
		case <-dpi.update:
			s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
		}
	}
}

// Allocate allocates a set of devices to be used by container runtime environment.
func (dpi *NUMACellDevicePlugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response pluginapi.AllocateResponse

	dpi.update <- message{}

	glog.Infof("Allocate() called with %+v", r)
	for _, container := range r.ContainerRequests {
		if len(container.DevicesIDs) != 1 {
			return nil, fmt.Errorf("can't allocate more than 1 numacell")
		}
		if !strings.HasPrefix(container.DevicesIDs[0], NUMACellName) {
			return nil, fmt.Errorf("cannot allocate numacell %q", container.DevicesIDs[0])
		}
		cellID := container.DevicesIDs[0]
		envCellID, err := strconv.Atoi(cellID[len(NUMACellName):])
		if err != nil {
			return nil, fmt.Errorf("unrecognized numacell format: %q error: %v", cellID, err)
		}

		dev := new(pluginapi.DeviceSpec)
		dev.HostPath = "/dev/null"      // TODO
		dev.ContainerPath = "/dev/null" // TODO
		dev.Permissions = "rw"

		containerResp := new(pluginapi.ContainerAllocateResponse)
		containerResp.Devices = []*pluginapi.DeviceSpec{dev}
		// this is only meant to improve debuggability
		containerResp.Envs = map[string]string{
			"IO_OPENSHIFT_KNI_NUMA_CELL_ID": fmt.Sprintf("%d", envCellID),
		}

		response.ContainerResponses = append(response.ContainerResponses, containerResp)
	}
	glog.Infof("AllocateResponse send: %+v", response)
	return &response, nil
}

// GetDevicePluginOptions returns options to be communicated with Device
// Manager
func (NUMACellDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return nil, nil
}

// PreStartContainer is called, if indicated by Device Plugin during registeration phase,
// before each container start. Device plugin can run device specific operations
// such as reseting the device before making devices available to the container
func (NUMACellDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}
