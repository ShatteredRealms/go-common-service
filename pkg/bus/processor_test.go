package bus_test

import (
	"context"
	"errors"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
)

var _ = Describe("Bus processor", func() {
	log.Logger, _ = test.NewNullLogger()
	var (
		testRepo TestBusMessageRepository
		testBus  TestingBus
		bp       bus.BusProcessor
	)
	BeforeEach(func() {
		testRepo = TestBusMessageRepository{}
		testBus = TestingBus{
			Queue: make(chan *TestBusMessage, 10),
		}
		bp = &bus.DefaultBusProcessor[TestBusMessage]{
			Reader: &testBus,
			Repo:   &testRepo,
		}
	})

	Describe("Start and starting processing", func() {
		Context("when starting processing", func() {
			It("should work with a valid context", func() {
				bp.StartProcessing(context.Background())
				Eventually(bp.IsProcessing).Should(BeTrue())
			})

			It("should work with a nil context", func() {
				bp.StartProcessing(nil)
				Eventually(bp.IsProcessing).Should(BeTrue())
			})

			It("should only start once", func() {
				bp.StartProcessing(nil)
				Eventually(bp.IsProcessing).Should(BeTrue())
				bp.StartProcessing(nil)
				Expect(bp.IsProcessing()).To(BeTrue())
			})
		})
		Context("when stopping processing", func() {
			It("should work when called once", func() {
				Expect(bp.IsProcessing()).To(BeFalse())
				bp.StopProcessing()
				Expect(bp.IsProcessing()).To(BeFalse())
			})

			It("should work when called multiple times", func() {
				Expect(bp.IsProcessing()).To(BeFalse())
				bp.StopProcessing()
				Expect(bp.IsProcessing()).To(BeFalse())
				bp.StopProcessing()
				Expect(bp.IsProcessing()).To(BeFalse())
			})

			It("should work if started", func() {
				bp.StartProcessing(nil)
				Expect(bp.IsProcessing()).To(BeTrue())
				bp.StopProcessing()
				Eventually(bp.IsProcessing).Should(BeFalse())
			})
		})
	})
	Describe("processing messages", func() {
		Context("bus fetch errors", func() {
			When("continue indefinitely", func() {
				It("should stop processing", func() {
					testBus.ErrOnFetch = errors.New("fetch error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					Eventually(bp.IsProcessing).Should(BeFalse())
				})
			})
			When("happen sparsely", func() {
				It("should continue processing", func() {
					testBus.ErrOnFetch = errors.New("fetch error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					testBus.ErrOnFetch = nil
					Expect(bp.IsProcessing()).To(BeTrue())
				})
			})
			It("should continue trying processing if fetches don't continue", func() {
				testBus.ErrOnFetch = errors.New("fetch error")
				Eventually(bp.IsProcessing).Should(BeFalse())
				bp.StartProcessing(nil)
				Expect(bp.IsProcessing()).To(BeTrue())
				Eventually(bp.IsProcessing).Should(BeFalse())
			})
		})
		Context("saving", func() {
			BeforeEach(func() {
				testBus.Queue <- &TestBusMessage{Id: faker.Username()}
			})
			When("saving fails continuously", func() {
				It("should stop processing", func() {
					testRepo.ErrOnSave = errors.New("save error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					Eventually(bp.IsProcessing).Should(BeFalse())
					Expect(testBus.CurrentMessage).ToNot(BeNil())
				})
			})
			When("saving mostly succeeds", func() {
				It("should continue processing", func() {
					testBus.Queue <- &TestBusMessage{Id: faker.Username()}
					testRepo.ErrOnSave = errors.New("save error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					testRepo.ErrOnSave = nil
					Expect(bp.IsProcessing()).To(BeTrue())
					Eventually(testBus.Queue).Should(BeEmpty())
					Eventually(testBus.CurrentMessage).Should(BeNil())
				})
			})
		})
		Context("deleting", func() {
			BeforeEach(func() {
				testBus.Queue <- &TestBusMessage{Id: faker.Username(), Deleted: true}
			})
			When("deleting fails continuously", func() {
				It("should stop processing", func() {
					testRepo.ErrOnDelete = errors.New("delete error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					Eventually(bp.IsProcessing).Should(BeFalse())
					Expect(testBus.CurrentMessage).ToNot(BeNil())
				})
			})
			When("deleting mostly succeeds", func() {
				It("should continue processing", func() {
					testRepo.ErrOnDelete = errors.New("delete error")
					bp.StartProcessing(nil)
					Expect(bp.IsProcessing()).To(BeTrue())
					testRepo.ErrOnDelete = nil
					Expect(bp.IsProcessing()).To(BeTrue())
					Eventually(testBus.Queue).Should(BeEmpty())
					Eventually(testBus.CurrentMessage).Should(BeNil())
				})
			})
		})
	})
})
