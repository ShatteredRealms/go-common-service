package config_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
)

var _ = Describe("ServerAddress(s)", func() {
	var (
		serverAddress *config.ServerAddress
	)

	BeforeEach(func() {
		serverAddress = &config.ServerAddress{
			Host: "localhost",
			Port: "8080",
		}
	})

	Describe("Address", func() {
		It("should work", func() {
			Expect(serverAddress.Address()).To(Equal("localhost:8080"))
		})
	})

	Describe("Addresses", func() {
		It("should work", func() {
			serverAddresses := config.ServerAddresses{
				{Host: faker.Username(), Port: faker.Username()},
				{Host: faker.Username(), Port: faker.Username()},
			}
			Expect(serverAddresses.Addresses()).To(HaveLen(2))
			Expect(serverAddresses.Addresses()[0]).To(Equal(serverAddresses[0].Address()))
			Expect(serverAddresses.Addresses()[1]).To(Equal(serverAddresses[1].Address()))
		})
	})
})
