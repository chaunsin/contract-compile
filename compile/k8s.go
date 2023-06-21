package compile

import (
	"context"

	"github.com/chaunsin/contract-compile/compile/cmd"
)

type k8sClient struct {
	cfg *Config
}

func NewK8s(c *Config) (ContractCompile, error) {
	cli := &k8sClient{
		cfg: c,
	}

	return cli, nil
}

func (c *k8sClient) Ping(ctx context.Context) error {
	return nil
}

func (c *k8sClient) Close(ctx context.Context) error {
	return nil
}

func (c *k8sClient) Execute(ctx context.Context, args cmd.Args) error {

	return nil
}
