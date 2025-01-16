// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package goredis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gozix/di"
	"github.com/gozix/glue/v3"
	gzViper "github.com/gozix/viper/v3"

	"github.com/spf13/viper"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct{}

// BundleName is default definition name.
const BundleName = "redis"

// Bundle implements the glue.Bundle interface.
var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder di.Builder) error {
	return builder.Provide(b.provideRegistry)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		gzViper.BundleName,
	}
}

func (b *Bundle) provideRegistry(cfg *viper.Viper) (*Registry, func() error, error) {
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
			Host:             cfg.GetString(suffix + "host"),
			Port:             cfg.GetString(suffix + "port"),
			DB:               cfg.GetInt(suffix + "db"),
			Username:         cfg.GetString(suffix + "username"),
			Password:         cfg.GetString(suffix + "password"),
			MaxRetries:       cfg.GetInt(suffix + "max_retries"),
			IdleTimeout:      cfg.GetDuration(suffix + "idle_timeout"),
			ReadTimeout:      cfg.GetDuration(suffix + "read_timeout"),
			WriteTimeout:     cfg.GetDuration(suffix + "write_timeout"),
			DisableIndentity: cfg.GetBool(suffix + "disable_indentity"),
		}

		// validating
		if c.Host == "" {
			return nil, nil, errors.New(suffix + "host should be set")
		}

		if c.DB < 0 {
			return nil, nil, errors.New(suffix + "db should be greater or equal to 0")
		}

		if c.MaxRetries < 0 {
			return nil, nil, errors.New(suffix + "max_retries should be greater or equal to 0")
		}

		conf[name] = c
	}

	var (
		registry = NewRegistry(conf)
		closer   = func() error {
			return registry.Close()
		}
	)

	return registry, closer, nil
}
