// Package goredis provides implementation of go-redis client.
package goredis

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// DEFAULT is default connection name.
const DEFAULT = "default"

// ConfigKey is root config key.
const configKey = "redis"

type (
	// Config is registry configuration item.
	Config struct {
		Host         string        `json:"host"`
		Port         string        `json:"port"`
		DB           int           `json:"db"`
		Password     string        `json:"password"`
		MaxRetries   int           `json:"max_retries"`
		IdleTimeout  time.Duration `json:"idle_timeout"`
		ReadTimeout  time.Duration `json:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout"`
	}

	// Configs is registry configurations.
	Configs map[string]Config

	// Registry is database connection registry.
	Registry struct {
		mux     sync.Mutex
		clients map[string]*redis.Client
		conf    Configs
	}
)

var (
	// ErrUnknownConnection is error triggered when connection with provided name not founded.
	ErrUnknownConnection = errors.New("unknown connection")
)

// NewRegistry is registry constructor.
func NewRegistry(conf Configs) *Registry {
	return &Registry{
		clients: make(map[string]*redis.Client, 1),
		conf:    conf,
	}
}

// Close is method for close connections.
func (r *Registry) Close() (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for key, client := range r.clients {
		if errClose := client.Close(); errClose != nil {
			err = errClose
		}

		delete(r.clients, key)
	}

	return err
}

// Connection is default connection getter.
func (r *Registry) Connection() (*redis.Client, error) {
	return r.ConnectionWithName(DEFAULT)
}

// ConnectionWithName is connection getter by name.
func (r *Registry) ConnectionWithName(name string) (_ *redis.Client, err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var client, initialized = r.clients[name]
	if initialized {
		return client, nil
	}

	var cfg, exists = r.conf[name]
	if !exists {
		return nil, ErrUnknownConnection
	}

	var options = &redis.Options{
		Addr: net.JoinHostPort(
			cfg.Host,
			cfg.Port,
		),
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	if client = redis.NewClient(options); client == nil {
		return nil, err
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err = client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	r.clients[name] = client

	return client, nil
}
