package k8s_device_plugins_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	e2etest "github.com/fromanirh/k8s-device-plugins/test/e2e"
)

var _ = BeforeSuite(func() {
	err := e2etest.Setup()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := e2etest.Teardown()
	Expect(err).ToNot(HaveOccurred())
})

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}
