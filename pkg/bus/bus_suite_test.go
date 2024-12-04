package bus_test

import (
	"testing"

	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	kafkaCloseFunc func() error
	kafkaPort      string
)

func TestBus(t *testing.T) {
	// log.Logger, _ = test.NewNullLogger()
	var err error

	SynchronizedBeforeSuite(func() []byte {
		kafkaCloseFunc, kafkaPort, err = testsro.SetupKafkaWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(kafkaPort).NotTo(BeEmpty())

		return []byte(kafkaPort)
	}, func(data []byte) {
		kafkaPort = string(data)
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		if kafkaCloseFunc != nil {
			kafkaCloseFunc()
		}
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Bus Suite")
}
