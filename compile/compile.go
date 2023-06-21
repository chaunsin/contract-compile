package compile

import (
	"context"
	"fmt"

	"github.com/chaunsin/contract-compile/compile/cmd"
)

type ContractCompile interface {
	Execute(ctx context.Context, args cmd.Args) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}

type Engine string

const (
	Docker Engine = "docker"
	K8S    Engine = "k8s"
)

func New(kind Engine, cfg *Config) (ContractCompile, error) {
	switch kind {
	case K8S:
		return NewK8s(cfg)
	case Docker:
		return NewDocker(cfg)
	default:
		return nil, fmt.Errorf("%s unknown engine type", kind)
	}
}
