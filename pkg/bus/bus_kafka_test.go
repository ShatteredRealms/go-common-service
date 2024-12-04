package bus_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-common-service/pkg/bus"
	"github.com/ShatteredRealms/go-common-service/pkg/config"
)

var _ = Describe("Kafka", func() {
	var (
		groupId   string
		busWriter bus.MessageBusWriter[TestBusMessage]
		busReader bus.MessageBusReader[TestBusMessage]
		msg       TestBusMessage
	)
	BeforeEach(func() {
		groupId = fmt.Sprintf("test-%d-%d", GinkgoParallelProcess(), time.Now().UnixNano())
		msg = TestBusMessage{
			Type: bus.BusMessageType(strconv.FormatInt(time.Now().UnixNano(), 10)),
			Id:   faker.Username(),
		}
		busWriter = bus.NewKafkaMessageBusWriter(
			config.ServerAddresses{{
				Host: "localhost",
				Port: kafkaPort,
			}},
			msg,
		)
		busReader = bus.NewKafkaMessageBusReader(
			config.ServerAddresses{{
				Host: "localhost",
				Port: kafkaPort,
			}},
			groupId,
			msg,
		)
	})
	AfterEach(func() {
		Expect(busWriter.Close()).ToNot(HaveOccurred())
		Expect(busReader.Close()).ToNot(HaveOccurred())
	})
	Describe("Kafka bus writer", func() {
		Context("closing", func() {
			When("Writer is open", func() {
				It("should close the connection", func(ctx SpecContext) {
					Eventually(ctx, func() error {
						return busWriter.Publish(ctx, msg)
					}).ShouldNot(HaveOccurred())
					go func() {
						for ctx.Err() == nil && busWriter != nil {
							busWriter.Publish(ctx, TestBusMessage{
								Type: msg.GetType(),
								Id:   faker.Username(),
							})
						}
					}()
					Expect(func() { busWriter.Close() }).ToNot(Panic())
				})
			})
			When("Writer is closed", func() {
				It("nothing should happen", func() {
					Expect(func() {
						Expect(busWriter.Close()).ToNot(HaveOccurred())
					}).ToNot(Panic())
				})
			})
		})
	})
	Describe("Kafka bus reader", func() {
		Context("fetching", func() {
			When("Messages are available", func() {
				BeforeEach(func(ctx SpecContext) {
					Eventually(ctx, func() error {
						return busWriter.Publish(ctx, msg)
					}).ShouldNot(HaveOccurred())
				})
				Context("processing", func() {
					var err error
					var outMsg *TestBusMessage
					BeforeEach(func(ctx SpecContext) {
						outMsg, err = busReader.FetchMessage(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(outMsg).NotTo(BeNil())
						Expect(outMsg).NotTo(BeIdenticalTo(&msg))
						Expect(*outMsg).To(Equal(msg))
					})
					When("ProcessSucceeded is called", func() {
						It("should mark the message as completed", func(ctx SpecContext) {
							Expect(busReader.ProcessSucceeded(ctx)).ToNot(HaveOccurred())
						})

						It("should require a message to be fetched", func(ctx SpecContext) {
							Expect(busReader.ProcessSucceeded(ctx)).NotTo(HaveOccurred())
							Expect(busReader.ProcessSucceeded(ctx)).To(HaveOccurred())
						})

						It("should error if the reader is closed", func(ctx SpecContext) {
							Expect(busReader.Close()).ToNot(HaveOccurred())
							Expect(busReader.ProcessSucceeded(ctx)).To(HaveOccurred())
						})
					})
					When("ProcessFailed is called", func() {
						It("should reject the message to and be processed again", func(ctx SpecContext) {
							Expect(busReader.ProcessFailed()).ToNot(HaveOccurred())
							outMsg, err = busReader.FetchMessage(ctx)
							Expect(err).ToNot(HaveOccurred())
							Expect(outMsg).NotTo(BeNil())
							Expect(outMsg).NotTo(BeIdenticalTo(&msg))
							Expect(*outMsg).To(Equal(msg))
							Expect(busReader.ProcessFailed()).ToNot(HaveOccurred())
						})

						It("should require a message to be fetched", func(ctx SpecContext) {
							Expect(busReader.ProcessSucceeded(ctx)).NotTo(HaveOccurred())
							Expect(busReader.ProcessFailed()).To(HaveOccurred())
						})

						It("should error if the reader was closed", func(ctx SpecContext) {
							Expect(busReader.Close()).ToNot(HaveOccurred())
							Expect(busReader.ProcessFailed()).To(HaveOccurred())
						})
					})
				})
				It("should return the messages in order", func(ctx SpecContext) {
					msg2 := TestBusMessage{
						Type: msg.GetType(),
						Id:   "2" + msg.Id,
					}
					Eventually(ctx, func() error {
						return busWriter.Publish(ctx, msg2)
					}).ShouldNot(HaveOccurred())
					outMsg, err := busReader.FetchMessage(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(outMsg).NotTo(BeNil())
					Expect(outMsg).NotTo(BeIdenticalTo(&msg))
					Expect(*outMsg).To(Equal(msg))
					Expect(*outMsg).NotTo(Equal(msg2))
					busReader.ProcessSucceeded(ctx)
					outMsg, err = busReader.FetchMessage(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(outMsg).NotTo(BeNil())
					Expect(outMsg).NotTo(BeIdenticalTo(&msg2))
					Expect(*outMsg).To(Equal(msg2))
					Expect(*outMsg).NotTo(Equal(msg))
				})
				It("should only allow fetching one message at a time", func(ctx SpecContext) {
					outMsg, err := busReader.FetchMessage(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(outMsg).NotTo(BeNil())
					Expect(outMsg).NotTo(BeIdenticalTo(&msg))
					Expect(*outMsg).To(Equal(msg))
					outMsg, err = busReader.FetchMessage(ctx)
					Expect(err).To(HaveOccurred())
					Expect(outMsg).To(BeNil())
				})
			})
			When("Messages are not available", func() {
				It("should hang", func(gCtx SpecContext) {
					ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
					defer cancel()
					Eventually(gCtx, func() error {
						_, err := busReader.FetchMessage(ctx)
						return err
					}).ShouldNot(Succeed())
				})
			})
		})
		Context("closing", func() {
			When("Reader is open", func() {
				It("should close the connection", func(ctx SpecContext) {
					var err error
					var outMsg *TestBusMessage
					Eventually(ctx, func() error {
						return busWriter.Publish(ctx, msg)
					}).ShouldNot(HaveOccurred())

					outMsg, err = busReader.FetchMessage(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(outMsg).NotTo(BeNil())
					Expect(outMsg).NotTo(BeIdenticalTo(&msg))
					Expect(*outMsg).To(Equal(msg))
					Expect(func() { busReader.Close() }).ToNot(Panic())
				})
			})
			When("Reader is closed", func() {
				It("nothing should happen", func() {
					Expect(func() {
						Expect(busReader.Close()).ToNot(HaveOccurred())
					}).ToNot(Panic())
				})
			})
		})
	})
})
