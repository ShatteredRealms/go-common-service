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
	gdbCloseFunc   func()
	mdb            *mongo.Database
	mdbCloseFunc   func()
	redisCloseFunc func()

	data initializeData
)

func TestRepository(t *testing.T) {
	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()

		var gormPort string
		gdbCloseFunc, gormPort = testsro.SetupGormWithDocker()
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr = testsro.SetupMongoWithDocker()
		Expect(mdbCloseFunc).NotTo(BeNil())

		redisCloseFunc, data.RedisConfig = testsro.SetupRedisWithDocker()

		data.GormConfig = config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: gormPort,
				Host: "localhost",
			},
			Name:     testsro.DbName,
			Username: testsro.Username,
			Password: testsro.Password,
		}
		gdb = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(gdb).NotTo(BeNil())
		mdb = testsro.ConnectMongoDocker(data.MdbConnStr)
		Expect(mdb).NotTo(BeNil())

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(inBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		gdb = testsro.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(gdb).NotTo(BeNil())
		mdb = testsro.ConnectMongoDocker(data.MdbConnStr)
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
