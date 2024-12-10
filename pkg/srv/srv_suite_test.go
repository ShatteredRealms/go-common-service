package srv_test

import (
	"bytes"
	"context"
	"encoding/gob"
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

type initializeData struct {
	KeycloakHost string
	AdminToken   string
}

var (
	keycloak *gocloak.GoCloak

	admin = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testadmin"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("adminfirstname"),
		LastName:      gocloak.StringP("adminlastname"),
		Email:         gocloak.StringP("admin@example.com"),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}

	incAdminCtx context.Context

	kcCfg = config.KeycloakConfig{
		Realm:        "default",
		Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
		ClientId:     "sro-character-service",
		ClientSecret: "**********",
	}
)

func TestSrv(t *testing.T) {
	var keycloakCloseFunc func() error
	var data initializeData

	SynchronizedBeforeSuite(func() []byte {
		log.Logger, _ = test.NewNullLogger()
		var err error
		keycloakCloseFunc, kcCfg.BaseURL, err = testsro.SetupKeycloakWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(kcCfg.BaseURL).NotTo(BeNil())

		keycloak = gocloak.NewClient(kcCfg.BaseURL)
		Expect(keycloak).NotTo(BeNil())

		clientToken, err := keycloak.LoginClient(
			context.Background(),
			kcCfg.ClientId,
			kcCfg.ClientSecret,
			kcCfg.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		setupUser := func(user *gocloak.User, roleName string, tokenStr *string) {
			*user.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, kcCfg.Realm, *user)
			Expect(err).NotTo(HaveOccurred())
			role, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, kcCfg.Realm, roleName)
			Expect(err).NotTo(HaveOccurred())
			err = keycloak.AddRealmRoleToUser(
				context.Background(),
				clientToken.AccessToken,
				kcCfg.Realm,
				*user.ID,
				[]gocloak.Role{*role},
			)
			Expect(err).NotTo(HaveOccurred())
			var token *gocloak.JWT
			Eventually(func() error {
				token, err = keycloak.Login(
					context.Background(),
					kcCfg.ClientId,
					kcCfg.ClientSecret,
					kcCfg.Realm,
					*user.Username,
					*(*user.Credentials)[0].Value,
				)
				return err
			}).Within(time.Minute).Should(Succeed())
			Expect(err).NotTo(HaveOccurred())
			(*tokenStr) = token.AccessToken
		}

		setupUser(&admin, "super admin", &data.AdminToken)
		data.KeycloakHost = kcCfg.BaseURL

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(inBytes []byte) {
		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		kcCfg.BaseURL = data.KeycloakHost

		keycloak = gocloak.NewClient(kcCfg.BaseURL)
		Expect(keycloak).NotTo(BeNil())

		Expect(data.AdminToken).NotTo(BeEmpty())
		md := metadata.New(
			map[string]string{
				"authorization": "Bearer " + data.AdminToken,
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
	RunSpecs(t, "Srv Suite")
}
