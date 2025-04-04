package srv_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ShatteredRealms/go-common-service/pkg/srv"
)

var _ = Describe("Health server", func() {
	It("should work", func() {
		srv := srv.NewHealthServiceServer()
		out, err := srv.Health(context.TODO(), &emptypb.Empty{})
		Expect(err).NotTo(HaveOccurred())
		Expect(out.Status).To(Equal("ok"))
	})
})

