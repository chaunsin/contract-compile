package compile

import (
	"context"
	"testing"

	"github.com/chaunsin/contract-compile/compile/cmd"
	"github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		ctx = context.TODO()
		cfg = Config{
			AuthConfig: registry.AuthConfig{},
			Mode:       "",
			Host:       "",
			CaCertPath: "",
			CertPath:   "",
			KeyPath:    "",
		}
		args = cmd.Args{
			Organization: "eth",
			Images:       "ethereum/solc:0.6.0",
			HostDir:      "/Users/edy/code/contract-compile/testdata",
			TargetDir:    "/data",
			Overwrite:    false,
			Extend:       nil,
		}
	)
	handle, err := New(Docker, &cfg)
	assert.NoError(t, err)
	assert.NoError(t, handle.Execute(ctx, args))
}
