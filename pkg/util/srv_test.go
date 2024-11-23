package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
)

var _ = Describe("Srv util", func() {
	var (
		hook *test.Hook
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		hook.Reset()
	})

	Describe("GrpcDialOpts", func() {
		It("should create DialOptions", func() {
			Expect(util.GrpcDialOpts()).NotTo(BeEmpty())
		})
	})

	Describe("GrpcClientWithOtel", func() {
		It("should dial address without error", func() {
			client, err := util.GrpcClientWithOtel("127.0.0.1:9999")
			Expect(client).NotTo(BeNil())
			Expect(err).To(Succeed())
		})
	})

	Describe("InitServerDefaults", func() {
		It("should create default server and mux", func() {
			server, mux := util.InitServerDefaults(nil, "")
			Expect(server).NotTo(BeNil())
			Expect(mux).NotTo(BeNil())
		})
	})
})
