package srv_test

import (
	"strings"

	"github.com/WilSimpson/gocloak/v13"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/srv"
)

var _ = Describe("Context", func() {
	var context *srv.Context
	var cfg *config.BaseConfig
	var logLevel logrus.Level
	BeforeEach(func() {
		Expect(faker.FakeData(&logLevel)).To(Succeed())
		cfg = &config.BaseConfig{
			Keycloak: kcCfg,
			LogLevel: logLevel,
		}
		context = srv.NewContext(cfg, faker.Username())
		Expect(context).NotTo(BeNil())
	})

	Describe("NewContext", func() {
		It("should setup keycloak", func() {
			Expect(context.KeycloakClient).NotTo(BeNil())
		})
		It("should setup tracer", func() {
			Expect(context.Tracer).NotTo(BeNil())
		})
		It("should set the log level", func() {
			Expect(log.Logger.Level).To(Equal(logLevel))
		})
	})

	Describe("GetJwt", func() {
		When("valid keycloak credentials are provided", func() {
			It("should return a valid JWT", func(ctx SpecContext) {
				jwt, err := context.GetJWT(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(jwt).NotTo(BeNil())
				Expect(jwt.AccessToken).NotTo(BeEmpty())
				Expect(strings.Split(jwt.AccessToken, ".")).To(HaveLen(3))
			})
			It("should should cache the JWT", func(ctx SpecContext) {
				jwt, err := context.GetJWT(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(jwt).NotTo(BeNil())
				Expect(jwt.AccessToken).NotTo(BeEmpty())
				Expect(strings.Split(jwt.AccessToken, ".")).To(HaveLen(3))
				jwt2, err := context.GetJWT(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(jwt).To(Equal(jwt2))
			})
		})
		When("invalid keycloak credentials are provided", func() {
			It("should return an error", func(ctx SpecContext) {
				cfg.Keycloak.ClientSecret = faker.Password()
				context = srv.NewContext(cfg, faker.Username())
				jwt, err := context.GetJWT(ctx)
				Expect(err).To(HaveOccurred())
				Expect(jwt).To(BeNil())
			})
		})
	})

	Describe("CreateRoles", func() {
		var roles []*gocloak.Role
		BeforeEach(func() {
			roles = []*gocloak.Role{
				{
					Name:        gocloak.StringP(faker.Username()),
					ClientRole:  gocloak.BoolP(true),
					ContainerID: new(string),
					Description: gocloak.StringP(faker.Sentence()),
				},
			}
		})
		When("valid keycloak credentials are provided", func() {
			It("should create client roles if the do not exist", func(ctx SpecContext) {
				Expect(context.CreateRoles(ctx, &roles)).To(Succeed())
			})
			It("should ignore errors if they already exist", func(ctx SpecContext) {
				Expect(context.CreateRoles(ctx, &roles)).To(Succeed())
				Expect(context.CreateRoles(ctx, &roles)).To(Succeed())
			})
		})
		When("invalid keycloak credentials are provided", func() {
			It("should return an error", func(ctx SpecContext) {
				cfg.Keycloak.ClientSecret = faker.Password()
				context = srv.NewContext(cfg, faker.Username())
				Expect(context.CreateRoles(ctx, nil)).To(HaveOccurred())
			})
		})
	})
})
