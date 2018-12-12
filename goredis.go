// Package goredis provides implementation of go-redis client.
package goredis

import (
	"errors"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/gozix/viper"
	"github.com/sarulabs/di"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Pool is type alias of redis.Client
	Pool = redis.Client

	// redisConf is logger configuration struct.
	redisConf struct {
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	}
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
			var cnf *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cnf); err != nil {
				return nil, errors.New("can't get config from container")
			}

			var conf redisConf
			if err = cnf.UnmarshalKey("redis", &conf); err != nil {
				return nil, err
			}

			options := &redis.Options{
				Addr: net.JoinHostPort(
					cnf.GetString("redis.host"),
					cnf.GetString("redis.port"),
				),
				IdleTimeout:  conf.IdleTimeout,
				ReadTimeout:  conf.ReadTimeout,
				WriteTimeout: conf.WriteTimeout,
			}

			var client *redis.Client
			if client = redis.NewClient(options); client == nil {
				return nil, err
			}

			if _, err = client.Ping().Result(); err != nil {
				return nil, err
			}

			return client, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*redis.Client).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		viper.BundleName,
	}
}
