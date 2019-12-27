// Package goredis provides implementation of go-redis client.
package goredis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gozix/viper/v2"
	"github.com/sarulabs/di/v2"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Pool is type alias of redis.Client
	Pool = redis.Client
)

// BundleName is default definition name.
const BundleName = "redis"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
				return nil, errors.New("can't get config from container")
			}

			// use this is hack, not UnmarshalKey
			// see https://github.com/spf13/viper/issues/188
			var (
				keys = cfg.Sub(configKey).AllKeys()
				conf = make(Configs, len(keys))
			)

			for _, key := range keys {
				var name = strings.Split(key, ".")[0]
				if _, ok := conf[name]; ok {
					continue
				}

				var suffix = fmt.Sprintf("%s.%s.", configKey, name)

				cfg.SetDefault(suffix+"port", "6379")

				var c = Config{
					Host:         cfg.GetString(suffix + "host"),
					Port:         cfg.GetString(suffix + "port"),
					DB:           cfg.GetInt(suffix + "db"),
					Password:     cfg.GetString(suffix + "password"),
					MaxRetries:   cfg.GetInt(suffix + "max_retries"),
					IdleTimeout:  cfg.GetDuration(suffix + "idle_timeout"),
					ReadTimeout:  cfg.GetDuration(suffix + "read_timeout"),
					WriteTimeout: cfg.GetDuration(suffix + "write_timeout"),
				}

				// validating
				if c.Host == "" {
					return nil, errors.New(suffix + "host should be set")
				}

				if c.MaxRetries < 0 {
					return nil, errors.New(suffix + "max_retries should be greater or equal to 0")
				}

				conf[name] = c
			}

			return NewRegistry(conf), nil
		},
		Close: func(obj interface{}) error {
			return obj.(*Registry).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		viper.BundleName,
	}
}
