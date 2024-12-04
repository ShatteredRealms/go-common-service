package config_test

import (
	"os"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
)

const (
	envKey         = "SRO_FOO"
	nestedEnvKey   = "SRO_BAR_BAZ"
	embeddedEnvKey = "SRO_BUZ"
)

type TestStruct struct {
	EmbeddedStruct `yaml:",inline" mapstructure:",squash"`

	Foo string
	Bar TestInnerStruct
}

type EmbeddedStruct struct {
	Buz string
}

type TestInnerStruct struct {
	Baz string
}

var _ = Describe("Config bind", func() {
	log.Logger, _ = test.NewNullLogger()
	var (
		testStruct  *TestStruct
		originalVal string
		envVal      string
	)

	BeforeEach(func() {
		Expect(faker.FakeData(&testStruct)).Should(Succeed())
	})

	Context("config file", func() {
		Context("valid config", func() {
			BeforeEach(func(ctx SpecContext) {
				Expect(config.BindConfigEnvs(ctx, "test_config", testStruct)).To(Succeed())
			})
			It("should bind to nested structs", func() {
				Expect(testStruct.Bar.Baz).To(Equal("testbaz"))
			})
			It("should bind to non-nested structs", func() {
				Expect(testStruct.Foo).To(Equal("testfoo"))
			})
			It("should bind to embedded structs", func() {
				Expect(testStruct.Buz).To(Equal("testbuz"))
			})
		})
		Context("invalid config", func() {
			It("should error if not valid yaml", func(ctx SpecContext) {
				Expect(config.BindConfigEnvs(ctx, "invalid_test_config", testStruct)).NotTo(Succeed())
			})
		})
	})

	Context("environment variables", func() {
		It("should bind to nested structs", func(ctx SpecContext) {
			originalVal = testStruct.Bar.Baz
			envVal = originalVal + faker.Username()
			GinkgoT().Setenv(nestedEnvKey, envVal)
			Expect(os.Getenv(nestedEnvKey)).To(Equal(envVal))
			Expect(config.BindConfigEnvs(ctx, faker.Username(), testStruct)).To(Succeed())
			Expect(testStruct.Bar.Baz).To(Equal(envVal))
		})

		It("should bind to non-nested structs", func(ctx SpecContext) {
			originalVal = testStruct.Foo
			envVal = faker.Name()
			GinkgoT().Setenv(envKey, envVal)
			Expect(os.Getenv(envKey)).To(Equal(envVal))
			Expect(config.BindConfigEnvs(ctx, faker.Username(), testStruct)).To(Succeed())
			Expect(testStruct.Foo).To(Equal(envVal))
		})

		It("should bind to embedded structs", func(ctx SpecContext) {
			originalVal = testStruct.Buz
			envVal = faker.Name()
			GinkgoT().Setenv(embeddedEnvKey, envVal)
			Expect(os.Getenv(embeddedEnvKey)).To(Equal(envVal))
			Expect(config.BindConfigEnvs(ctx, faker.Username(), testStruct)).To(Succeed())
			Expect(testStruct.Buz).To(Equal(envVal))
		})
	})
})
