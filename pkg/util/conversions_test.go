package util_test

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/util"
)

var _ = Describe("Conversion util", func() {
	Describe("ArrayOfUint64ToUint", func() {
		It("should work", func() {
			ints, err := faker.RandomInt(10, 100, 1)
			Expect(err).To(Succeed())
			in := make([]uint64, ints[0])
			for idx := range in {
				in[idx] = uint64(faker.RandomUnixTime())
			}
			out := *util.ArrayOfUint64ToUint(&in)

			for idx := range out {
				Expect(out[idx]).To(Equal(uint(in[idx])), "each index should match")
			}
		})
	})

	Describe("ParseUUIDs", func() {
		Context("invalid input", func() {
			It("should err when all invalid uuid", func() {
				out, err := util.ParseUUIDs([]string{"asdf"})
				Expect(out).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
			It("should err when any uuid is invalid", func() {
				out, err := util.ParseUUIDs([]string{uuid.NewString(), "asdf"})
				Expect(out).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("valid input", func() {
			It("should work", func() {
				ints, err := faker.RandomInt(5, 50, 1)
				Expect(err).To(BeNil())
				in := make([]string, ints[0])
				for idx := range in {
					in[idx] = uuid.NewString()
				}

				out, err := util.ParseUUIDs(in)
				Expect(err).To(BeNil())
				Expect(len(out)).To(Equal(len(in)))
				for idx := range out {
					Expect(out[idx].String()).To(Equal(in[idx]))
				}
			})
		})
	})
})
