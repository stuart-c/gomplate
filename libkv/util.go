package libkv

import (
	"github.com/hairyhenderson/gomplate/env"
	"github.com/hairyhenderson/gomplate/typeconv"
	consulapi "github.com/hashicorp/consul/api"
)

func setupTLS(prefix string) *consulapi.TLSConfig {
	tlsConfig := &consulapi.TLSConfig{
		Address:  env.Getenv(prefix + "_TLS_SERVER_NAME"),
		CAFile:   env.Getenv(prefix + "_CACERT"),
		CAPath:   env.Getenv(prefix + "_CAPATH"),
		CertFile: env.Getenv(prefix + "_CLIENT_CERT"),
		KeyFile:  env.Getenv(prefix + "_CLIENT_KEY"),
	}
	if v := env.Getenv(prefix + "_HTTP_SSL_VERIFY"); v != "" {
		verify := typeconv.MustParseBool(v)
		tlsConfig.InsecureSkipVerify = !verify
	}
	return tlsConfig
}
