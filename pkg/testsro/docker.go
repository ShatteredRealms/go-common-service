package testsro

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/cenkalti/backoff"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	Username = "postgres"
	Password = "password"
	DbName   = "test"
)

func SetupKeycloakWithDocker() (closeFn func() error, host string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	if pool.Client.Ping() != nil {
		return nil, "", errors.New("docker not running")
	}

	pool.MaxWait = time.Second * 10
	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	realmExportFile, err := filepath.Abs(fmt.Sprintf("%s/../../pkg/testsro/keycloak-realm-export.json", wd))
	if err != nil {
		return
	}
	keycloakRunDockerOpts := &dockertest.RunOptions{
		Repository: "quay.io/keycloak/keycloak",
		Tag:        "21.0.0",
		Env:        []string{"KEYCLOAK_ADMIN=admin", "KEYCLOAK_ADMIN_PASSWORD=admin"},
		Cmd: []string{
			"start-dev",
			"--import-realm",
			"--health-enabled=true",
			"--features=declarative-user-profile",
		},
		Mounts:       []string{fmt.Sprintf("%s:/opt/keycloak/data/import/realm-export.json", realmExportFile)},
		ExposedPorts: []string{"8080/tcp"},
	}

	keycloakResource, err := pool.RunWithOptions(keycloakRunDockerOpts, fnConfig)
	if err != nil {
		return
	}

	// Uncomment to see docker log
	// go func() {
	// 	pool.Client.Logs(docker.LogsOptions{
	// 		Container:    keycloakResource.Container.ID,
	// 		OutputStream: log.Logger.Out,
	// 		ErrorStream:  log.Logger.Out,
	// 		Follow:       true,
	// 		Stdout:       true,
	// 		Stderr:       true,
	// 		Timestamps:   false,
	// 		RawTerminal:  true,
	// 	})
	// }()

	closeFn = func() error {
		return keycloakResource.Close()
	}

	host = "http://127.0.0.1:" + keycloakResource.GetPort("8080/tcp")
	err = pool.Retry(func() error {
		_, err := http.Get(host + "/health/ready")
		return err
	})

	err = pool.Retry(func() error {
		_, err = http.Get(host + "/realms/default")
		return err
	})
	if err != nil {
		return
	}

	return
}

func SetupKafkaWithDocker() (closeFn func() error, port string, err error) {
	var pool *dockertest.Pool
	pool, err = dockertest.NewPool("")
	if err != nil {
		return
	}

	if pool.Client.Ping() != nil {
		return nil, "", errors.New("docker not running")
	}

	net, err := pool.CreateNetwork(fmt.Sprintf("go-testing-%d", time.Now().UnixNano()), func(config *docker.CreateNetworkOptions) {
		// config.Driver = "bridge"
	})
	if err != nil {
		return
	}

	pool.MaxWait = time.Second * 10
	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	}

	zookeeperResource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Hostname:     "zookeeper",
			Repository:   "confluentinc/cp-zookeeper",
			Tag:          "latest",
			Env:          []string{"ZOOKEEPER_CLIENT_PORT=2181"},
			ExposedPorts: []string{"2181/tcp"},
			Networks:     []*dockertest.Network{net},
		},
		fnConfig,
	)
	if err != nil {
		return
	}

	kafkaResource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Hostname:   "kafka",
			Repository: "confluentinc/cp-kafka",
			Tag:        "latest",
			Env: []string{
				"KAFKA_BROKER_ID=1",
				"KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181",
				"KAFKA_LISTENERS=PLAINTEXT://:29092,PLAINTEXT_HOST://:59092",
				"KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://:29092,PLAINTEXT_HOST://:59092",
				"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
				"KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT",
				"KAFKA_JMX_PORT=9101",
				"KAFKA_JMX_HOSTNAME=localhost",
				"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
				"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1",
				"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1",
				"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0",
				"KAFKA_ADVERTISED_HOST_NAME=localhost",
			},
			ExposedPorts: []string{"59092/tcp", "29092/tcp", "9101/tcp"},
			// PortBindings: map[docker.Port][]docker.PortBinding{
			// 	"59092/tcp": {{HostIP: "localhost", HostPort: "59092/tcp"}},
			// },
			Networks: []*dockertest.Network{net},
		},
		fnConfig)
	if err != nil {
		return
	}

	// Wait for kafka to be ready
	err = pool.Retry(func() error {
		code, err := kafkaResource.Exec([]string{
			"kafka-topics",
			"--bootstrap-server",
			"localhost:29092",
			"--list",
		}, dockertest.ExecOptions{
			// StdOut: os.Stdout,
			// StdErr: os.Stderr,
		})
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("kafka-topics not ready: code %d", code)
		}
		return nil
	})
	if err != nil {
		return
	}

	// Set advertised listeners do the docker forwarded ports
	err = pool.Retry(func() error {
		code, err := kafkaResource.Exec([]string{
			"kafka-configs",
			"--bootstrap-server",
			"localhost:29092",
			"--entity-type",
			"brokers",
			"--entity-name",
			"1",
			"--alter",
			"--add-config",
			fmt.Sprintf(
				"advertised.listeners=[PLAINTEXT://:29092,PLAINTEXT_HOST://:%s]",
				kafkaResource.GetPort("59092/tcp"),
			),
		}, dockertest.ExecOptions{
			// StdOut: os.Stdout,
			// StdErr: os.Stderr,
		})
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("kafka-config able to set: code %d", code)
		}
		return nil
	})
	if err != nil {
		return
	}

	// Create test topic
	err = pool.Retry(func() error {
		code, err := kafkaResource.Exec([]string{
			"kafka-topics",
			"--bootstrap-server",
			"localhost:29092",
			"--create",
			"--if-not-exists",
			"--topic",
			"test",
			"--replication-factor",
			"1",
			"--partitions",
			"1",
		}, dockertest.ExecOptions{
			// StdOut: os.Stdout,
			// StdErr: os.Stderr,
		})
		if err != nil {
			return err
		}
		if code != 0 {
			return fmt.Errorf("kafka-topics not ready: code %d", code)
		}
		return nil
	})
	if err != nil {
		return
	}

	closeFn = func() error {
		var err error
		errors.Join(err, kafkaResource.Close())
		errors.Join(err, zookeeperResource.Close())
		errors.Join(net.Close())
		return err
	}

	port = kafkaResource.GetPort("59092/tcp")
	return
}

func SetupMongoWithDocker() (closeFn func() error, host string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	runDockerOpt := &dockertest.RunOptions{
		Repository: "mongo", // image
		Tag:        "6.0",   // version
		Env:        []string{"MONGO_INITDB_ROOT_USERNAME=root", "MONGO_INITDB_ROOT_PASSWORD=password"},
	}

	pool.MaxWait = time.Second * 10
	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	if err != nil {
		return
	}
	// call clean up function to release resource
	closeFn = func() error {
		return resource.Close()
	}

	host = fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))
	return
}

func ConnectMongoDocker(host string) (mdb *mongo.Database, err error) {
	err = Retry(func() error {
		db, err := mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				host,
			),
		)
		if err != nil {
			return err
		}
		mdb = db.Database("testdb")
		return db.Ping(context.TODO(), nil)
	}, time.Second*30)

	return
}

func SetupGormWithDocker() (closeFn func() error, port string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env:        []string{"POSTGRES_PASSWORD=" + Password, "POSTGRES_DB=" + DbName},
	}

	pool.MaxWait = time.Second * 10
	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	if err != nil {
		return
	}
	// call clean up function to release resource
	closeFn = func() error {
		return resource.Close()
	}

	port = resource.GetPort("5432/tcp")
	return
}

func ConnectGormDocker(connStr string) (gdb *gorm.DB, err error) {
	// retry until db server is ready
	err = Retry(func() (err error) {
		gdb, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: logger.New(
				log.Logger,
				logger.Config{
					SlowThreshold:             0,
					Colorful:                  true,
					IgnoreRecordNotFoundError: true,
					ParameterizedQueries:      true,
					LogLevel:                  logger.Info,
				},
			),
		})
		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}, time.Second*10)

	return
}

func SetupRedisWithDocker() (closeFn func() error, cfg config.DBPoolConfig, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return
	}

	runDockerOpt := &dockertest.RunOptions{
		Repository: "grokzen/redis-cluster", // image
		Tag:        "7.0.10",                // version
		Env: []string{
			"INITIAL_PORT=7000",
			"MASTERS=3",
			"SLAVES_PER_MASTER=1",
		},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	if err != nil {
		return
	}
	// call clean up function to release resource
	closeFn = func() error {
		return resource.Close()
	}

	// container is ready, return *gorm.Db for testing
	cfg = config.DBPoolConfig{
		Master: config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: resource.GetPort("7000/tcp"),
				Host: "localhost",
			},
		},
		Slaves: []config.DBConfig{
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7001/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7002/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7003/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7004/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7005/tcp"),
					Host: "localhost",
				},
			},
		},
	}
	return
}

func Retry(op func() error, timeout time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = timeout
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return fmt.Errorf("reached retry deadline: %w", err)
		}

		return err
	}

	return nil
}
