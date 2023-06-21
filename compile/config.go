package compile

import "github.com/docker/docker/api/types/registry"

type Config struct {
	// 用于访问私有镜像托管仓库访问配置,比如访问私有得docker-hub
	registry.AuthConfig

	Mode       string
	Host       string
	CaCertPath string
	CertPath   string
	KeyPath    string
}

func (c *Config) Valid() error {
	return nil
}
