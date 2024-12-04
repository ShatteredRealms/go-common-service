package repository_test

import (
	"context"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Db repository", func() {
	var (
		pool config.DBPoolConfig
	)
	BeforeEach(func() {
		pool = config.DBPoolConfig{
			Master: data.GormConfig,
			Slaves: []config.DBConfig{data.GormConfig},
		}
	})
	Describe("ConnectDb", func() {
		When("given valid input", func() {
			It("should work", func() {
				out, err := repository.ConnectDB(context.TODO(), pool, data.RedisConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error", func() {
				pool.Master.Host = "a"
				out, err := repository.ConnectDB(context.TODO(), pool, data.RedisConfig)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

