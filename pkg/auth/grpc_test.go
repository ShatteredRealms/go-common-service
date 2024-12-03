package auth_test

import (
	"context"
	"io"

	"github.com/ShatteredRealms/go-common-service/pkg/auth"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/go-faker/faker/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	iauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Auth gRPC", func() {
	log.Logger.SetOutput(io.Discard)

	Context("public service methods", func() {
		var callMeta interceptors.CallMeta
		BeforeEach(func() {
			callMeta = interceptors.CallMeta{
				Service: faker.Username(),
				Method:  faker.Username(),
			}
		})
		Describe("RegisterPublicServiceMethods", func() {
			It("should work", func() {
				Expect(func() { auth.RegisterPublicServiceMethods() }).NotTo(Panic())
				Expect(func() { auth.RegisterPublicServiceMethods("") }).NotTo(Panic())
				Expect(func() { auth.RegisterPublicServiceMethods("a") }).NotTo(Panic())
				Expect(func() { auth.RegisterPublicServiceMethods("a", "a") }).NotTo(Panic())
			})
		})

		Describe("NotPublicServiceMatcher", func() {
			It("should allow known services", func() {
				Expect(auth.NotPublicServiceMatcher(nil, interceptors.CallMeta{Service: "sro.HealthService", Method: "Health"})).To(BeFalse())
			})
			It("should deny unknown services", func() {
				Expect(auth.NotPublicServiceMatcher(nil, callMeta)).To(BeTrue())
			})
		})

		Describe("registering public methods", func() {
			It("should allow after registering", func() {
				Expect(auth.NotPublicServiceMatcher(nil, callMeta)).To(BeTrue())
				Expect(func() { auth.RegisterPublicServiceMethods(callMeta.FullMethod()) }).NotTo(Panic())
				Expect(auth.NotPublicServiceMatcher(nil, callMeta)).To(BeFalse())
			})
		})
	})

	Describe("AuthFunc", func() {
		Context("invalid parameters", func() {
			It("should return an authorization error with invalid keycloak client", func() {
				fn := auth.AuthFunc(nil, "default")
				inCtx := context.Background()
				outCtx, err := fn(inCtx)
				Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
				Expect(outCtx).To(BeNil())
			})

			It("should return an authorization error with invalid realm", func() {
				fn := auth.AuthFunc(keycloak, "")
				inCtx := context.Background()
				outCtx, err := fn(inCtx)
				Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
				Expect(outCtx).To(BeNil())
			})
		})

		Context("valid keycloak client", func() {
			var fn iauth.AuthFunc
			BeforeEach(func() {
				Expect(keycloak).NotTo(BeNil())
				Expect(func() { fn = auth.AuthFunc(keycloak, "default") }).NotTo(Panic())
			})
			When("invalid context", func() {
				It("should fail with empty context", func() {
					inCtx := context.Background()
					outCtx, err := fn(inCtx)
					Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
					Expect(outCtx).To(BeNil())
				})
				It("should fail with invalid authorization scheme", func() {
					md := metadata.New(
						map[string]string{
							"authorization": faker.Word(),
						},
					)
					inCtx := metadata.NewIncomingContext(context.Background(), md)
					outCtx, err := fn(inCtx)
					Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
					Expect(outCtx).To(BeNil())
				})
				It("should fail with empty authorization", func() {
					md := metadata.New(
						map[string]string{
							"authorization": "",
						},
					)
					inCtx := metadata.NewIncomingContext(context.Background(), md)
					outCtx, err := fn(inCtx)
					Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
					Expect(outCtx).To(BeNil())
				})
			})
			When("given invalid authorization", func() {
				It("should fail given an invalid token", func() {
					md := metadata.New(
						map[string]string{
							"authorization": "Bearer " + faker.Username(),
						},
					)
					inCtx := metadata.NewIncomingContext(context.Background(), md)
					outCtx, err := fn(inCtx)
					Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
					Expect(outCtx).To(BeNil())
				})
				It("should fail given a random token", func() {
					md := metadata.New(
						map[string]string{
							"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
						},
					)
					inCtx := metadata.NewIncomingContext(context.Background(), md)
					outCtx, err := fn(inCtx)
					Expect(err).To(Equal(auth.ErrUnauthorized.Err()))
					Expect(outCtx).To(BeNil())
				})
			})

			It("should pass with valid token", func() {
				md := metadata.New(
					map[string]string{
						"authorization": "Bearer " + adminToken.AccessToken,
					},
				)
				inCtx := metadata.NewIncomingContext(context.Background(), md)
				outCtx, err := fn(inCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(outCtx).NotTo(BeNil())

				claims, ok := auth.RetrieveClaims(outCtx)
				Expect(ok).To(BeTrue())
				Expect(claims).NotTo(BeNil())
			})
		})
	})
})
