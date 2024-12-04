package config_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
)

var _ = Describe("Db config", func() {
	var (
		config *config.DBConfig
		dsn    string
	)

	BeforeEach(func() {
		Expect(faker.FakeData(&config)).To(Succeed())
		dsn = ""
	})

	Describe("MySQLDSN", func() {
		It("should work", func() {
			dsn = config.MySQLDSN()
			Expect(dsn).To(ContainSubstring(config.Name))
		})
	})

	Describe("PostgresDSN", func() {
		It("should work", func() {
			dsn = config.PostgresDSN()
			Expect(dsn).To(ContainSubstring(config.Name))
		})
	})

	Describe("PostgresDSNWithoutName", func() {
		It("should work", func() {
			dsn = config.PostgresDSNWithoutName()
			Expect(dsn).NotTo(ContainSubstring(config.Name))
		})
	})

	Describe("MongoDSN", func() {
		It("should work", func() {
			dsn = config.MongoDSN()
		})
	})

	AfterEach(func() {
		Expect(dsn).To(ContainSubstring(config.Username))
		Expect(dsn).To(ContainSubstring(config.Password))
		Expect(dsn).To(ContainSubstring(config.ServerAddress.Host))
		Expect(dsn).To(ContainSubstring(config.ServerAddress.Port))
	})
})

var _ = Describe("Db pool config", func() {
	var (
		config *config.DBPoolConfig
	)

	It("should be able to list addresses and ports", func() {
		Expect(faker.FakeData(&config)).To(Succeed())
		addrs := config.Addresses()
		Expect(addrs).To(HaveLen(len(config.Slaves) + 1))
	})
})
