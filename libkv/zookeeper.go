package libkv

import (
	"net/url"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/hairyhenderson/gomplate/env"
	"github.com/hairyhenderson/gomplate/typeconv"
)

// NewZookeeper - instantiate a new Zookeeper datasource handler
func NewZookeeper(u *url.URL) *LibKV {
	zookeeper.Register()
	c := zookeeperURL(u)
	config := zookeeperConfig()
	kv, err := libkv.NewStore(store.ZK, []string{c.String()}, config)
	if err != nil {
		logFatal("Zookeeper setup failed", err)
	}
	return &LibKV{kv}
}

// -- converts a gomplate datasource URL into a usable Zookeeper URL
func zookeeperURL(u *url.URL) *url.URL {
	c, _ := url.Parse(env.Getenv("ZK_ADDR"))
	if c.Scheme == "" {
		c.Scheme = u.Scheme
	}

	if c.Host == "" && u.Host == "" {
		c.Host = "localhost:2181"
	} else if c.Host == "" {
		c.Host = u.Host
	}

	return c
}

func zookeeperConfig() *store.Config {
	t := typeconv.MustAtoi(env.Getenv("ZK_TIMEOUT"))
	config := &store.Config{
		ConnectionTimeout: time.Duration(t) * time.Second,
	}
	return config
}
