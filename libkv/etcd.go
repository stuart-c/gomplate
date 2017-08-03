package libkv

import (
	"net/url"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/hairyhenderson/gomplate/env"
	"github.com/hairyhenderson/gomplate/typeconv"
	consulapi "github.com/hashicorp/consul/api"
)

// NewEtcd - instantiate a new ETCD datasource handler
func NewEtcd(u *url.URL) *LibKV {
	etcd.Register()
	c := etcdURL(u)
	config := etcdConfig(c.Scheme == "https")
	kv, err := libkv.NewStore(store.ETCD, []string{c.String()}, config)
	if err != nil {
		logFatal("Consul setup failed", err)
	}
	return &LibKV{kv}
}

// -- converts a gomplate datasource URL into a usable ETCD URL
func etcdURL(u *url.URL) *url.URL {
	c, _ := url.Parse(env.Getenv("ETCD_ADDR"))
	if c.Scheme == "" {
		c.Scheme = u.Scheme
	}
	switch c.Scheme {
	case "consul+http", "http":
		c.Scheme = "http"
	case "consul+https", "https":
		c.Scheme = "https"
	case "consul":
		if typeconv.MustParseBool(env.Getenv("ETCD_TLS")) {
			c.Scheme = "https"
		} else {
			c.Scheme = "http"
		}
	}

	if c.Host == "" && u.Host == "" {
		c.Host = "localhost:2379"
	} else if c.Host == "" {
		c.Host = u.Host
	}

	return c
}

func etcdConfig(useTLS bool) *store.Config {
	t := typeconv.MustAtoi(env.Getenv("ETCD_TIMEOUT"))
	config := &store.Config{
		ConnectionTimeout: time.Duration(t) * time.Second,
	}
	if useTLS {
		tconf := setupTLS("ETCD")
		var err error
		config.TLS, err = consulapi.SetupTLSConfig(tconf)
		if err != nil {
			logFatal("TLS Config setup failed", err)
		}
	}
	config.Username = env.Getenv("ETCD_USERNAME", "")
	config.Password = env.Getenv("ETCD_PASSWORD", "")
	return config
}
