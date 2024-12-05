package util_test

import (
	"github.com/WilSimpson/gocloak/v13"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/util"
)

var _ = Describe("Misc", func() {
	Describe("RegisterRole", func() {
		It("should append the role to the roles slice", func() {
			roles := []*gocloak.Role{}
			role := &gocloak.Role{}
			role2 := &gocloak.Role{}
			Expect(faker.FakeData(role)).To(Succeed())
			Expect(faker.FakeData(role2)).To(Succeed())
			Expect(util.RegisterRole(role, &roles)).To(Equal(role))
			Expect(roles).To(ContainElement(role))
			Expect(util.RegisterRole(role2, &roles)).To(Equal(role2))
			Expect(roles).To(ContainElement(role2))
			Expect(roles).To(HaveLen(2))
		})
	})
})
