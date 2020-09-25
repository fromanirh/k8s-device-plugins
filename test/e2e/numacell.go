package k8s_device_plugins

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	"github.com/fromanirh/k8s-device-plugins/pkg/numacell"
)

func countNUMACells(resources corev1.ResourceList) int {
	expectedResName := fmt.Sprintf("%s/%s", numacell.NUMACellResourceNamespace, numacell.NUMACellResourceName) // TODO

	count := 0
	for resName := range resources {
		if strings.HasPrefix(string(resName), expectedResName) {
			count++
		}
	}
	return count
}

func getNUMATestPod(cpus int, numacellid int) *corev1.Pod {
	pod := GetTestPod()
	numacellResourceName := numacell.MakeResourceName(numacellid)
	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", cpus)),
			numacellResourceName:  resource.MustParse("1"),
		},
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", cpus)),
			numacellResourceName:  resource.MustParse("1"),
		},
	}
	return pod
}

func checkNUMACellsOnNode(node *corev1.Node) {
	numacellAllocatable := countNUMACells(node.Status.Allocatable)
	Expect(numacellAllocatable).To(BeNumerically(">", 0), "missing numacells from allocatable resources on %q", node.Name)

	numacellCapacity := countNUMACells(node.Status.Capacity)
	Expect(numacellCapacity).To(BeNumerically(">", 0), "missing numacells from allocatable resources on %q", node.Name)

	Expect(numacellAllocatable).To(Equal(numacellCapacity), "allocatable=%d/capacity=%d mismatch on %q", numacellAllocatable, numacellCapacity, node.Name)
}

var _ = Describe("numacell device plugin", func() {
	Context("with node objects", func() {
		It("Should be declare the resources", func() {
			nodes, err := GetByRole("worker")
			Expect(err).ToNot(HaveOccurred())

			for _, node := range nodes {
				checkNUMACellsOnNode(&node)
			}
		})

		It("Should declare infinite resources", func() {
			nodes, err := GetByRole("worker")
			Expect(err).ToNot(HaveOccurred())

			for _, node := range nodes {
				numacellAllocatable := countNUMACells(node.Status.Allocatable)
				Expect(numacellAllocatable).To(BeNumerically(">", 0), "missing numacells from allocatable resources on %q", node.Name)
			}

			testpod := getNUMATestPod(8, 0)
			testpod.Namespace = TestingNamespace.Name

			err = Client.Create(context.TODO(), testpod)
			Expect(err).ToNot(HaveOccurred())

			err = WaitForPodCondition(testpod, corev1.PodReady, corev1.ConditionTrue, 10*time.Minute)
			Expect(err).ToNot(HaveOccurred())

			var key types.NamespacedName

			key.Namespace = testpod.Namespace
			key.Name = testpod.Name
			updatedPod := &corev1.Pod{}
			err = Client.Get(context.TODO(), key, updatedPod)
			Expect(err).ToNot(HaveOccurred())

			key.Namespace = ""
			key.Name = updatedPod.Spec.NodeName
			updatedNode := &corev1.Node{}
			err = Client.Get(context.TODO(), key, updatedNode)
			Expect(err).ToNot(HaveOccurred())

			checkNUMACellsOnNode(updatedNode)
		})
	})
})
