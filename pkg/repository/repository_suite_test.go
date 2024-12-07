package repository_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/testsro"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type initializeData struct {
	GormConfig  config.DBConfig
	MdbConnStr  string
	RedisConfig config.DBPoolConfig
}

var (
	hook *test.Hook

	gdb            *gorm.DB
	gdbCloseFunc   func() error
	mdb            *mongo.Database
	mdbCloseFunc   func() error
	redisCloseFunc func() error

	data initializeData
)

func TestRepository(t *testing.T) {
	var err error
	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()

		var gormPort string
		gdbCloseFunc, gormPort, err = testsro.SetupGormWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr, err = testsro.SetupMongoWithDocker()
		Expect(err).NotTo(HaveOccurred())
		Expect(mdbCloseFunc).NotTo(BeNil())

		redisCloseFunc, data.RedisConfig, err = testsro.SetupRedisWithDocker()
		Expect(err).NotTo(HaveOccurred())

		data.GormConfig = config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: gormPort,
				Host: "localhost",
			},
			Name:     testsro.DbName,
			Username: testsro.Username,
			Password: testsro.Password,
		}
		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(inBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		gdb, err = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(err).NotTo(HaveOccurred())
		Expect(gdb).NotTo(BeNil())
		mdb, err = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(mdb).NotTo(BeNil())
	})

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gdbCloseFunc()
		mdbCloseFunc()
		redisCloseFunc()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
