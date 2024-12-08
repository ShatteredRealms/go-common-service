package auth_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	"github.com/WilSimpson/gocloak/v13"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc/metadata"
)

var (
	keycloak *gocloak.GoCloak

	admin = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testadmin"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(false),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("adminfirstname"),
		LastName:      gocloak.StringP("adminlastname"),
		Email:         gocloak.StringP("admin@example.com"),
		Credentials: &[]gocloak.CredentialRepresentation{{
			Temporary: gocloak.BoolP(false),
			Type:      gocloak.StringP("password"),
			Value:     gocloak.StringP("Password1!"),
		}},
	}

	clientToken *gocloak.JWT
	adminToken  *gocloak.JWT

	incAdminCtx context.Context

	kcCfg = config.KeycloakConfig{
		Realm:        "default",
		Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
		ClientId:     "sro-character-service",
		ClientSecret: "**********",
	}
)

func TestAuth(t *testing.T) {
	var keycloakCloseFunc func() error

	SynchronizedBeforeSuite(func() []byte {
		log.Logger, _ = test.NewNullLogger()
		var host string
		var err error
		keycloakCloseFunc, host, err = testsro.SetupKeycloakWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(host).NotTo(BeNil())

		keycloak = gocloak.NewClient(host)
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			kcCfg.ClientId,
			kcCfg.ClientSecret,
			kcCfg.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		*admin.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, kcCfg.Realm, admin)
		Expect(err).NotTo(HaveOccurred())
		Expect(*admin.ID).NotTo(BeEmpty())

		saRole, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, kcCfg.Realm, "super admin")
		Expect(err).NotTo(HaveOccurred())

		err = keycloak.AddRealmRoleToUser(
			context.Background(),
			clientToken.AccessToken,
			kcCfg.Realm,
			*admin.ID,
			[]gocloak.Role{*saRole},
		)

		out := fmt.Sprintf("%s", host)

		return []byte(out)
	}, func(data []byte) {
		log.Logger, _ = test.NewNullLogger()
		splitData := strings.Split(string(data), "\n")
		Expect(splitData).To(HaveLen(1))

		host := splitData[0]

		keycloak = gocloak.NewClient(string(host))
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			kcCfg.ClientId,
			kcCfg.ClientSecret,
			kcCfg.Realm,
		)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() error {
			adminToken, err = keycloak.Login(
				context.Background(),
				kcCfg.ClientId,
				kcCfg.ClientSecret,
				kcCfg.Realm,
				*admin.Username,
				*(*admin.Credentials)[0].Value,
			)
			return err
		}).Within(time.Minute).Should(Succeed())

		admins, err := keycloak.GetUsers(
			context.Background(),
			clientToken.AccessToken,
			kcCfg.Realm,
			gocloak.GetUsersParams{Username: admin.Username},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(admins).To(HaveLen(1))
		admin = *admins[0]

		md := metadata.New(
			map[string]string{
				"authorization": "Bearer " + adminToken.AccessToken,
			},
		)
		incAdminCtx = metadata.NewIncomingContext(context.Background(), md)
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		if keycloakCloseFunc != nil {
			keycloakCloseFunc()
		}
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}
