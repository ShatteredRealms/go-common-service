package repository

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/ShatteredRealms/go-common-service/pkg/config"
	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
)

var (
	currentConn    *kafka.Conn
	controllerConn *kafka.Conn
)

func ConnectKafka(address config.ServerAddress) (*kafka.Conn, error) {
	if currentConn != nil {
		_ = currentConn.Close()
	}
	if controllerConn != nil {
		_ = controllerConn.Close()
	}
	var err error

	err = retry(func() error {
		currentConn, err = kafka.Dial("tcp", address.Address())
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("kafka connect: %v", err)
	}

	controller, err := currentConn.Controller()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return currentConn, nil
		}
		return nil, fmt.Errorf("controller: %v", err)
	}

	err = retry(func() error {
		controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("kafka controller connection: %v", err)
	}

	return controllerConn, nil
}

func retry(op func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second
	bo.MaxElapsedTime = time.Second * 15
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return fmt.Errorf("reached retry deadline: %v", err)
		}

		return err
	}

	return nil
}

