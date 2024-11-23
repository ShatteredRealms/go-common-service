package config

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/spf13/viper"
)

var (
	FailedEnvBinding = errors.New("failed to bind to environment variable(s)")
	FailedReadingConfig = errors.New("failed to read config file")
)

func BindConfigEnvs(ctx context.Context, name string, config *struct{}) error {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/sro/")
	viper.AddConfigPath("./test/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			log.Logger.WithContext(ctx).Errorf("read app config parse error: %v", err)
			return errors.Join(FailedEnvBinding, err)
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Logger.WithContext(ctx).Infof("Using default config: %v", err)
		} else {
			log.Logger.WithContext(ctx).Errorf("unknown error prasing config: %v", err)
			return errors.Join(FailedEnvBinding, err)
		}
	}

	viper.SetEnvPrefix("SRO")

	// Read from environment variables
	errs := bindEnvsToStruct(config)
	if len(errs) > 0 {
		log.Logger.WithContext(ctx).Errorf("failed to bind to %d environment variables", len(errs))
		for _, err := range errs {
			log.Logger.WithContext(ctx).Debugf("failed binding to env: %v", err)
		}
		return FailedEnvBinding
	}

	// Save to struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Logger.WithContext(ctx).Errorf("unmarshal appConfig: %v", err)
		return err
	}

	return nil
}

func bindEnvsToStruct(obj interface{}) []error {
	errs := make([]error, 0)

	viper.AutomaticEnv()

	val := reflect.ValueOf(obj)
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		key := field.Name
		if field.Anonymous {
			key = ""
		}

		err := bindRecursive(key, val.Field(i))
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func bindRecursive(key string, val reflect.Value) error {
	if val.Kind() != reflect.Struct {
		env := "SRO_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		return viper.BindEnv(key, env)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		newKey := field.Name
		if field.Anonymous {
			newKey = ""
		} else if key != "" {
			newKey = "." + newKey
		}

		bindRecursive(key+newKey, val.Field(i))
	}

	return nil
}
